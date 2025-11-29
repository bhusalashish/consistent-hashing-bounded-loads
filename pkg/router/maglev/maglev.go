package maglev

import (
	"sync"

	"github.com/bhusalashish/consistent-hashing-bounded-loads.git/pkg/hash"
	"github.com/bhusalashish/consistent-hashing-bounded-loads.git/pkg/routercore"
)

const defaultTableSize = 65537 // a prime, good default for Maglev

// mapper implements routercore.Mapper using the Maglev algorithm.
type mapper struct {
	mu    sync.RWMutex
	nodes []string // node IDs, indexable by table entries
	table []int    // slot -> node index
	m     int      // table size
	seed  uint64   // base seed for hashing
}

// NewMaglev constructs a new Maglev mapper.
//
// opts.TableSize controls M (table size). If zero or negative, a sensible
// default (defaultTableSize) is chosen. opts.HashSeed controls hashing.
func NewMaglev(nodes []string, opts routercore.Options) (routercore.Mapper, error) {
	m := &mapper{
		seed: opts.HashSeed,
	}

	if opts.TableSize > 0 {
		m.m = opts.TableSize
	} else {
		m.m = defaultTableSize
	}

	// initial build
	m.rebuild(nodes)
	return m, nil
}

func (m *mapper) Add(nodes ...string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.rebuild(append(m.nodes, nodes...))
}

func (m *mapper) Remove(nodes ...string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.nodes) == 0 || len(nodes) == 0 {
		return
	}
	removeSet := make(map[string]struct{}, len(nodes))
	for _, n := range nodes {
		removeSet[n] = struct{}{}
	}

	var kept []string
	for _, n := range m.nodes {
		if _, drop := removeSet[n]; !drop {
			kept = append(kept, n)
		}
	}

	m.rebuild(kept)
}

// Pick selects a node for the given key by hashing into the Maglev table.
func (m *mapper) Pick(key []byte) string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.nodes) == 0 {
		panic("maglev: no nodes registered")
	}
	if len(m.table) == 0 || m.m == 0 {
		panic("maglev: table not initialized")
	}

	h := hash.XXH64(key, m.seed)
	slot := int(h % uint64(m.m))
	nodeIdx := m.table[slot]

	if nodeIdx < 0 || nodeIdx >= len(m.nodes) {
		panic("maglev: invalid table entry; rebuild required")
	}

	return m.nodes[nodeIdx]
}

// rebuild rebuilds the Maglev lookup table for the given node list.
//
// It deduplicates nodes, computes per-node permutations, and fills
// the table so that each slot maps to exactly one node index.
func (m *mapper) rebuild(nodes []string) {
	// deduplicate nodes while preserving order
	seen := make(map[string]struct{}, len(nodes))
	var uniq []string
	for _, n := range nodes {
		if _, ok := seen[n]; ok {
			continue
		}
		seen[n] = struct{}{}
		uniq = append(uniq, n)
	}
	m.nodes = uniq

	// if no nodes, clear the table
	if len(m.nodes) == 0 {
		m.table = nil
		return
	}

	M := m.m
	table := make([]int, M)
	for i := range table {
		table[i] = -1
	}

	type permState struct {
		offset int
		skip   int
		next   int
	}

	perms := make([]permState, len(m.nodes))

	// Compute offset and skip per node using two hash streams.
	// We derive them from the same base seed with different mixes.
	const altSeed = 0x9e3779b97f4a7c15 // arbitrary odd constant for variation

	for i, id := range m.nodes {
		// h1 chooses starting offset
		h1 := hash.XXH64String(id, m.seed)
		// h2 chooses skip; ensure 1 <= skip <= M-1
		h2 := hash.XXH64String(id, m.seed^altSeed)

		offset := int(h1 % uint64(M))
		skip := int(h2%(uint64(M-1))) + 1

		perms[i] = permState{
			offset: offset,
			skip:   skip,
			next:   0,
		}
	}

	filled := 0
	for filled < M {
		for i := range perms {
			pos := (perms[i].offset + perms[i].next*perms[i].skip) % M
			perms[i].next++

			if table[pos] == -1 {
				table[pos] = i
				filled++
				if filled == M {
					break
				}
			}
		}
	}

	m.table = table
}
