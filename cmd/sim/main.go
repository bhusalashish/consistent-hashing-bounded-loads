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
	algo := flag.String("algo", "jump", "routing algorithm: jump | maglev | chbl")
	nodesN := flag.Int("nodes", 8, "number of nodes")
	keysN := flag.Int("keys", 100000, "number of keys to simulate")
	zipfS := flag.Float64("zipf-s", 0.0, "Zipf skew parameter s (0 = uniform)")
	tableSize := flag.Int("table-size", 65537, "Maglev table size (M)")
	loadFactor := flag.Float64("load-factor", 1.25, "CH-BL load factor c (>=1.0)")
	vnodes := flag.Int("vnodes", 100, "CH-BL virtual nodes per physical node")
	walkThreshold := flag.Int("walk-threshold", 8, "CH-BL walk threshold before two-choice fallback")
	seed := flag.Int64("seed", time.Now().UnixNano(), "random seed")
	outPath := flag.String("out", "", "output CSV file path (default stdout)")

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

	// ----- Nodes -----
	nodes := make([]string, *nodesN)
	for i := 0; i < *nodesN; i++ {
		nodes[i] = fmt.Sprintf("node-%d", i)
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

	// ----- Construct mapper via factory -----
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

	mapper, err := router.New(algoEnum, opts, nodes)
	if err != nil {
		log.Fatalf("failed to construct mapper: %v", err)
	}

	// ----- Workload generation -----
	rng := rand.New(rand.NewSource(*seed))

	// Per-node counts
	counts := make(map[string]int, *nodesN)

	if *zipfS <= 0.0 {
		// Uniform keys: sequential labels "key-0"..."key-N-1"
		for i := 0; i < *keysN; i++ {
			key := []byte(fmt.Sprintf("key-%d", i))
			n := mapper.Pick(key)
			counts[n]++
		}
	} else {
		// Zipf-skewed selection of key indices [0, keysN-1]
		// v parameter usually 1.0; s controls skew.
		zipf := rand.NewZipf(rng, *zipfS, 1.0, uint64(*keysN-1))
		for i := 0; i < *keysN; i++ {
			kIdx := zipf.Uint64()
			key := []byte(fmt.Sprintf("key-%d", kIdx))
			n := mapper.Pick(key)
			counts[n]++
		}
	}

	// ----- Build stats -----
	perNode := make([]int, 0, len(nodes))
	for _, id := range nodes {
		perNode = append(perNode, counts[id]) // keep consistent order
	}
	stats := metrics.ComputeIntStats(perNode)

	// ----- Prepare CSV output -----
	var out *os.File
	if *outPath == "" {
		out = os.Stdout
	} else {
		f, err := os.Create(*outPath)
		if err != nil {
			log.Fatalf("failed to create output file: %v", err)
		}
		defer f.Close()
		out = f
	}

	w := csv.NewWriter(out)
	defer w.Flush()

	// Header
	if err := w.Write([]string{
		"node_id",
		"count",
	}); err != nil {
		log.Fatalf("failed to write header: %v", err)
	}

	// Per-node rows
	for i, id := range nodes {
		row := []string{
			id,
			fmt.Sprintf("%d", perNode[i]),
		}
		if err := w.Write(row); err != nil {
			log.Fatalf("failed to write row: %v", err)
		}
	}

	// Summary as commented rows (prefixed with #), still valid CSV-ish.
	summaryRows := [][]string{
		{"#algo", *algo},
		{"#nodes", fmt.Sprintf("%d", *nodesN)},
		{"#keys", fmt.Sprintf("%d", *keysN)},
		{"#zipf_s", fmt.Sprintf("%.3f", *zipfS)},
		{"#table_size", fmt.Sprintf("%d", *tableSize)},
		{"#load_factor", fmt.Sprintf("%.3f", *loadFactor)},
		{"#vnodes", fmt.Sprintf("%d", *vnodes)},
		{"#walk_threshold", fmt.Sprintf("%d", *walkThreshold)},
		{"#seed", fmt.Sprintf("%d", *seed)},
		{"#mean", fmt.Sprintf("%.3f", stats.Mean)},
		{"#max", fmt.Sprintf("%d", stats.Max)},
		{"#std", fmt.Sprintf("%.3f", stats.Std)},
		{"#cv", fmt.Sprintf("%.5f", stats.CV)},
	}

	for _, row := range summaryRows {
		if err := w.Write(row); err != nil {
			log.Fatalf("failed to write summary row: %v", err)
		}
	}

	// Also log a one-line summary to stderr for quick inspection.
	log.Printf("algo=%s nodes=%d keys=%d zipf_s=%.2f mean=%.2f max=%d cv=%.4f",
		*algo, *nodesN, *keysN, *zipfS, stats.Mean, stats.Max, stats.CV)
}
