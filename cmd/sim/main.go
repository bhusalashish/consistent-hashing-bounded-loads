package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/bhusalashish/consistent-hashing-bounded-loads.git/pkg/metrics"
	router "github.com/bhusalashish/consistent-hashing-bounded-loads.git/pkg/router"
	rc "github.com/bhusalashish/consistent-hashing-bounded-loads.git/pkg/routercore"
)

func main() {
	// ----- Flags -----
	mode := flag.String("mode", "dist", "simulation mode: dist | churn")
	algo := flag.String("algo", "jump", "routing algorithm: jump | maglev | chbl")

	nodesN := flag.Int("nodes", 8, "number of nodes (before churn)")
	keysN := flag.Int("keys", 100000, "number of keys to simulate")
	zipfS := flag.Float64("zipf-s", 0.0, "Zipf skew parameter s (0 = uniform)")

	tableSize := flag.Int("table-size", 65537, "Maglev table size (M)")
	loadFactor := flag.Float64("load-factor", 1.25, "CH-BL load factor c (>=1.0)")
	vnodes := flag.Int("vnodes", 100, "CH-BL virtual nodes per physical node")
	walkThreshold := flag.Int("walk-threshold", 8, "CH-BL walk threshold before two-choice fallback")

	seed := flag.Int64("seed", time.Now().UnixNano(), "random seed")
	outPath := flag.String("out", "", "output CSV file path (default stdout)")

	churnOp := flag.String("churn-op", "", "churn operation in churn mode: add | remove")

	flag.Parse()

	if *nodesN <= 0 {
		log.Fatalf("nodes must be > 0")
	}
	if *keysN <= 0 {
		log.Fatalf("keys must be > 0")
	}
	if *loadFactor < 1.0 {
		log.Fatalf("load-factor must be >= 1.0 for CH-BL")
	}
	if *mode != "dist" && *mode != "churn" {
		log.Fatalf("mode must be 'dist' or 'churn'")
	}
	if *mode == "churn" && (*churnOp != "add" && *churnOp != "remove") {
		log.Fatalf("in churn mode, -churn-op must be 'add' or 'remove'")
	}

	// ----- Nodes (before churn) -----
	nodesBefore := make([]string, *nodesN)
	for i := 0; i < *nodesN; i++ {
		nodesBefore[i] = fmt.Sprintf("node-%d", i)
	}

	// ----- Algo enum -----
	var algoEnum rc.Algo
	switch *algo {
	case "jump":
		algoEnum = rc.AlgoJump
	case "maglev":
		algoEnum = rc.AlgoMaglev
	case "chbl":
		algoEnum = rc.AlgoCHBL
	default:
		log.Fatalf("unknown algo %q (expected jump|maglev|chbl)", *algo)
	}

	// ----- Router options -----
	opts := rc.Options{
		TableSize:     *tableSize,
		LoadFactor:    *loadFactor,
		Vnodes:        *vnodes,
		WalkThreshold: *walkThreshold,
		HashSeed:      uint64(*seed),
		ExpectedKeys:  *keysN, // CH-BL uses this; others ignore it
	}

	// ----- Pre-generate keys (so both phases use identical keys) -----
	keys := generateKeys(*keysN, *zipfS, *seed)

	// ----- Run appropriate mode -----
	switch *mode {
	case "dist":
		if err := runDistribution(*algo, algoEnum, nodesBefore, keys, opts, *zipfS, *seed, *outPath); err != nil {
			log.Fatalf("distribution run failed: %v", err)
		}
	case "churn":
		if err := runChurn(*algo, algoEnum, nodesBefore, keys, opts, *zipfS, *seed, *churnOp, *outPath); err != nil {
			log.Fatalf("churn run failed: %v", err)
		}
	}
}

func generateKeys(keysN int, zipfS float64, seed int64) [][]byte {
	keys := make([][]byte, keysN)
	rng := rand.New(rand.NewSource(seed))

	if zipfS <= 0.0 {
		for i := 0; i < keysN; i++ {
			keys[i] = []byte(fmt.Sprintf("key-%d", i))
		}
		return keys
	}

	zipf := rand.NewZipf(rng, zipfS, 1.0, uint64(keysN-1))
	for i := 0; i < keysN; i++ {
		kIdx := zipf.Uint64()
		keys[i] = []byte(fmt.Sprintf("key-%d", kIdx))
	}
	return keys
}

// ------------------ Distribution mode ------------------

func runDistribution(
	algoName string,
	algoEnum rc.Algo,
	nodes []string,
	keys [][]byte,
	opts rc.Options,
	zipfS float64,
	seed int64,
	outPath string,
) error {
	mapper, err := router.New(algoEnum, opts, nodes)
	if err != nil {
		return fmt.Errorf("construct mapper: %w", err)
	}

	// Count per node
	counts := make(map[string]int, len(nodes))
	for _, k := range keys {
		n := mapper.Pick(k)
		counts[n]++
	}

	perNode := make([]int, 0, len(nodes))
	for _, id := range nodes {
		perNode = append(perNode, counts[id]) // consistent order
	}
	stats := metrics.ComputeIntStats(perNode)

	out, w, err := createCSVWriter(outPath)
	if err != nil {
		return err
	}
	defer out.Close()
	defer w.Flush()

	// Header
	if err := w.Write([]string{"node_id", "count"}); err != nil {
		return fmt.Errorf("write header: %w", err)
	}
	// Rows
	for i, id := range nodes {
		if err := w.Write([]string{
			id,
			fmt.Sprintf("%d", perNode[i]),
		}); err != nil {
			return fmt.Errorf("write row: %w", err)
		}
	}

	// Summary meta rows
	summaryRows := [][]string{
		{"#mode", "dist"},
		{"#algo", algoName},
		{"#nodes", fmt.Sprintf("%d", len(nodes))},
		{"#keys", fmt.Sprintf("%d", len(keys))},
		{"#zipf_s", fmt.Sprintf("%.3f", zipfS)},
		{"#table_size", fmt.Sprintf("%d", opts.TableSize)},
		{"#load_factor", fmt.Sprintf("%.3f", opts.LoadFactor)},
		{"#vnodes", fmt.Sprintf("%d", opts.Vnodes)},
		{"#walk_threshold", fmt.Sprintf("%d", opts.WalkThreshold)},
		{"#seed", fmt.Sprintf("%d", seed)},
		{"#mean", fmt.Sprintf("%.3f", stats.Mean)},
		{"#max", fmt.Sprintf("%d", stats.Max)},
		{"#std", fmt.Sprintf("%.3f", stats.Std)},
		{"#cv", fmt.Sprintf("%.5f", stats.CV)},
	}

	for _, row := range summaryRows {
		if err := w.Write(row); err != nil {
			return fmt.Errorf("write summary row: %w", err)
		}
	}

	log.Printf("mode=dist algo=%s nodes=%d keys=%d zipf_s=%.2f mean=%.2f max=%d cv=%.4f",
		algoName, len(nodes), len(keys), zipfS, stats.Mean, stats.Max, stats.CV)

	return nil
}

// ------------------ Churn mode ------------------

func runChurn(
	algoName string,
	algoEnum rc.Algo,
	nodesBefore []string,
	keys [][]byte,
	opts rc.Options,
	zipfS float64,
	seed int64,
	churnOp string,
	outPath string,
) error {
	// Build nodesAfter
	var nodesAfter []string
	switch churnOp {
	case "add":
		nodesAfter = append([]string{}, nodesBefore...)
		newID := fmt.Sprintf("node-%d", len(nodesBefore))
		nodesAfter = append(nodesAfter, newID)
	case "remove":
		if len(nodesBefore) <= 1 {
			return fmt.Errorf("cannot remove from single-node cluster")
		}
		nodesAfter = append([]string{}, nodesBefore[:len(nodesBefore)-1]...)
	default:
		return fmt.Errorf("unknown churn-op %q", churnOp)
	}

	// Mapper before churn
	mapperBefore, err := router.New(algoEnum, opts, nodesBefore)
	if err != nil {
		return fmt.Errorf("construct mapper(before): %w", err)
	}
	// Mapper after churn
	mapperAfter, err := router.New(algoEnum, opts, nodesAfter)
	if err != nil {
		return fmt.Errorf("construct mapper(after): %w", err)
	}

	countsBefore := make(map[string]int, len(nodesBefore))
	countsAfter := make(map[string]int, len(nodesAfter))

	moved := 0
	total := len(keys)

	for _, k := range keys {
		nb := mapperBefore.Pick(k)
		na := mapperAfter.Pick(k)

		countsBefore[nb]++
		countsAfter[na]++

		if nb != na {
			moved++
		}
	}

	// Build unified node list: all nodes before, then any new ones
	nodeSeen := make(map[string]struct{})
	var nodeList []string
	for _, n := range nodesBefore {
		if _, ok := nodeSeen[n]; !ok {
			nodeSeen[n] = struct{}{}
			nodeList = append(nodeList, n)
		}
	}
	for _, n := range nodesAfter {
		if _, ok := nodeSeen[n]; !ok {
			nodeSeen[n] = struct{}{}
			nodeList = append(nodeList, n)
		}
	}

	perBefore := make([]int, len(nodeList))
	perAfter := make([]int, len(nodeList))
	for i, n := range nodeList {
		perBefore[i] = countsBefore[n]
		perAfter[i] = countsAfter[n]
	}

	statsBefore := metrics.ComputeIntStats(perBefore)
	statsAfter := metrics.ComputeIntStats(perAfter)

	movedRatio := float64(moved) / float64(total)

	out, w, err := createCSVWriter(outPath)
	if err != nil {
		return err
	}
	defer out.Close()
	defer w.Flush()

	// Header
	if err := w.Write([]string{"node_id", "count_before", "count_after"}); err != nil {
		return fmt.Errorf("write header: %w", err)
	}
	// Rows
	for i, n := range nodeList {
		if err := w.Write([]string{
			n,
			fmt.Sprintf("%d", perBefore[i]),
			fmt.Sprintf("%d", perAfter[i]),
		}); err != nil {
			return fmt.Errorf("write row: %w", err)
		}
	}

	// Summary meta rows
	summaryRows := [][]string{
		{"#mode", "churn"},
		{"#algo", algoName},
		{"#churn_op", churnOp},
		{"#nodes_before", fmt.Sprintf("%d", len(nodesBefore))},
		{"#nodes_after", fmt.Sprintf("%d", len(nodesAfter))},
		{"#keys", fmt.Sprintf("%d", total)},
		{"#moved", fmt.Sprintf("%d", moved)},
		{"#moved_ratio", fmt.Sprintf("%.6f", movedRatio)},
		{"#zipf_s", fmt.Sprintf("%.3f", zipfS)},
		{"#table_size", fmt.Sprintf("%d", opts.TableSize)},
		{"#load_factor", fmt.Sprintf("%.3f", opts.LoadFactor)},
		{"#vnodes", fmt.Sprintf("%d", opts.Vnodes)},
		{"#walk_threshold", fmt.Sprintf("%d", opts.WalkThreshold)},
		{"#seed", fmt.Sprintf("%d", seed)},
		{"#mean_before", fmt.Sprintf("%.3f", statsBefore.Mean)},
		{"#max_before", fmt.Sprintf("%d", statsBefore.Max)},
		{"#cv_before", fmt.Sprintf("%.5f", statsBefore.CV)},
		{"#mean_after", fmt.Sprintf("%.3f", statsAfter.Mean)},
		{"#max_after", fmt.Sprintf("%d", statsAfter.Max)},
		{"#cv_after", fmt.Sprintf("%.5f", statsAfter.CV)},
	}

	for _, row := range summaryRows {
		if err := w.Write(row); err != nil {
			return fmt.Errorf("write summary row: %w", err)
		}
	}

	log.Printf("mode=churn algo=%s churn_op=%s nodes_before=%d nodes_after=%d keys=%d moved=%d moved_ratio=%.4f",
		algoName, churnOp, len(nodesBefore), len(nodesAfter), total, moved, movedRatio)

	return nil
}

func createCSVWriter(outPath string) (*os.File, *csv.Writer, error) {
	var out *os.File
	if outPath == "" {
		out = os.Stdout
	} else {
		f, err := os.Create(outPath)
		if err != nil {
			return nil, nil, fmt.Errorf("create output file: %w", err)
		}
		out = f
	}
	w := csv.NewWriter(out)
	return out, w, nil
}
