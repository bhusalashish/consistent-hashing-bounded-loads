package chbl

import (
	"math"
	"sync"

	"github.com/bhusalashish/consistent-hashing-bounded-loads.git/internal/ring"
	"github.com/bhusalashish/consistent-hashing-bounded-loads.git/pkg/hash"
	"github.com/bhusalashish/consistent-hashing-bounded-loads.git/pkg/routercore"
)

const (
	defaultVnodes        = 100
	defaultLoadFactor    = 1.25
	defaultWalkThreshold = 8
)

type mapper struct {
	mu sync.Mutex

	nodes []string
	ring  *ring.Ring

	// per-node load and capacity
	load            []int
	capacityPerNode int

	// parameters
	vnodes        int
	loadFactor    float64
	walkThreshold int
	expectedKeys  int

	// hash seeds for first and second candidate
	seed1 uint64
	seed2 uint64
}

// NewCHBL constructs a bounded-load consistent-hashing mapper.
//
// It uses a vnode-based ring and enforces a per-node capacity:
//
//	C = ceil(c * ExpectedKeys / numNodes)
//
// where c = opts.LoadFactor (default 1.25).
// ExpectedKeys must be set by the caller for capacity guarantees to hold.
func NewCHBL(nodes []string, opts routercore.Options) (routercore.Mapper, error) {
	m := &mapper{
		vnodes:        defaultOrInt(opts.Vnodes, defaultVnodes),
		loadFactor:    defaultOrFloat(opts.LoadFactor, defaultLoadFactor),
		walkThreshold: defaultOrInt(opts.WalkThreshold, defaultWalkThreshold),
		expectedKeys:  opts.ExpectedKeys,
		seed1:         opts.HashSeed,
	}

	// derive a distinct second seed for two-choice fallback
	if m.seed1 == 0 {
		m.seed1 = 1 // avoid zero seed
	}
	m.seed2 = m.seed1 ^ 0x9e3779b97f4a7c15

	m.rebuild(nodes)
	return m, nil
}

func defaultOrInt(v, def int) int {
	if v <= 0 {
		return def
	}
	return v
}

func defaultOrFloat(v, def float64) float64 {
	if v <= 0 {
		return def
	}
	return v
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

// rebuild rebuilds the ring and resets load/capacity for the given nodes.
func (m *mapper) rebuild(nodes []string) {
	if len(nodes) == 0 {
		m.nodes = nil
		m.ring = nil
		m.load = nil
		m.capacityPerNode = 0
		return
	}

	// de-duplicate nodes, preserve order
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

	// rebuild ring
	m.ring = ring.New(m.nodes, m.vnodes, m.seed1)

	// compute capacity C = ceil(c * m / n)
	n := len(m.nodes)
	if m.expectedKeys <= 0 {
		// if ExpectedKeys is not set, we still define some capacity so that
		// the algorithm behaves reasonably; we default to avg * c for m = n.
		m.expectedKeys = n
	}
	avg := float64(m.expectedKeys) / float64(n)
	m.capacityPerNode = int(math.Ceil(m.loadFactor * avg))

	m.load = make([]int, n)
}

// Pick assigns the key to a node, enforcing the per-node capacity C and
// using a two-choice fallback if the linear walk gets too long.
//
// NOTE: This mapper is stateful over Pick calls (it tracks load).
func (m *mapper) Pick(key []byte) string {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.nodes) == 0 {
		panic("chbl: no nodes registered")
	}
	if m.ring == nil || len(m.ring.Tokens) == 0 {
		panic("chbl: ring not initialized")
	}

	h1 := hash.XXH64(key, m.seed1)
	idx := m.ring.SuccessorIndex(h1)
	startIdx := idx
	steps := 0

	for {
		token := m.ring.Tokens[idx]
		nodeIdx := token.NodeIdx

		if m.load[nodeIdx] < m.capacityPerNode {
			m.load[nodeIdx]++
			return m.nodes[nodeIdx]
		}

		steps++
		// two-choice fallback if walk becomes too long
		if steps == m.walkThreshold {
			chosen := m.twoChoiceFallback(key, nodeIdx)
			if chosen >= 0 {
				m.load[chosen]++
				return m.nodes[chosen]
			}
			// else continue walking from nodeIdx with smaller load
		}

		// move to next token on the ring
		idx++
		if idx == len(m.ring.Tokens) {
			idx = 0
		}
		if idx == startIdx {
			// We've looped around the whole ring and found no capacity.
			// This indicates that ExpectedKeys * LoadFactor is too low
			// for the actual call volume.
			panic("chbl: all nodes at capacity; increase ExpectedKeys or LoadFactor")
		}
	}
}

// twoChoiceFallback hashes the key again to get a second candidate and
// returns the index of the better node (less loaded and with capacity),
// or -1 if neither candidate has capacity.
func (m *mapper) twoChoiceFallback(key []byte, primaryIdx int) int {
	h2 := hash.XXH64(key, m.seed2)
	idx2 := m.ring.SuccessorIndex(h2)
	nodeIdx2 := m.ring.Tokens[idx2].NodeIdx

	// primary nodeIdx is the one we were walking from
	nodeIdx1 := primaryIdx

	has1 := m.load[nodeIdx1] < m.capacityPerNode
	has2 := m.load[nodeIdx2] < m.capacityPerNode

	if !has1 && !has2 {
		return -1
	}
	if has1 && !has2 {
		return nodeIdx1
	}
	if !has1 && has2 {
		return nodeIdx2
	}
	// both have capacity: pick the less loaded one
	if m.load[nodeIdx1] <= m.load[nodeIdx2] {
		return nodeIdx1
	}
	return nodeIdx2
}
