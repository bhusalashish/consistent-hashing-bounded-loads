package ringch

import (
	"sync"

	"github.com/bhusalashish/consistent-hashing-bounded-loads.git/internal/ring"
	"github.com/bhusalashish/consistent-hashing-bounded-loads.git/pkg/hash"
	"github.com/bhusalashish/consistent-hashing-bounded-loads.git/pkg/routercore"
)

const (
	defaultVnodes = 50 // fewer vnodes than CH-BL
)

// mapper implements a bare-minimum consistent-hashing router.
// No load caps, no bounded loads, no two-choice fallback.
// Simply: hash key → ring successor → node.
type mapper struct {
	mu    sync.RWMutex
	nodes []string
	rng   *ring.Ring

	vnodes   int
	hashSeed uint64
}

// NewRingCH constructs a basic CH router.
func NewRingCH(nodes []string, opts routercore.Options) (routercore.Mapper, error) {
	m := &mapper{
		vnodes:   defaultOrInt(opts.Vnodes, defaultVnodes),
		hashSeed: opts.HashSeed,
	}
	m.rebuild(nodes)
	return m, nil
}

// Add adds new nodes and rebuilds ring.
func (m *mapper) Add(nodes ...string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.rebuild(append(m.nodes, nodes...))
}

// Remove removes nodes and rebuilds ring.
func (m *mapper) Remove(nodes ...string) {
	m.mu.Lock()
	defer m.mu.Unlock()

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

func (m *mapper) rebuild(nodes []string) {
	if len(nodes) == 0 {
		m.nodes = nil
		m.rng = nil
		return
	}

	// Deduplicate while preserving order
	seen := make(map[string]struct{}, len(nodes))
	var uniq []string
	for _, n := range nodes {
		if _, exists := seen[n]; exists {
			continue
		}
		seen[n] = struct{}{}
		uniq = append(uniq, n)
	}
	m.nodes = uniq

	// Build ring
	m.rng = ring.New(m.nodes, m.vnodes, m.hashSeed)
}

func (m *mapper) Pick(key []byte) string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.nodes) == 0 {
		panic("ringch: no nodes registered")
	}
	if m.rng == nil {
		panic("ringch: ring not initialized")
	}

	h := hash.XXH64(key, m.hashSeed)
	idx := m.rng.SuccessorIndex(h)
	return m.nodes[m.rng.Tokens[idx].NodeIdx]
}

func defaultOrInt(v, def int) int {
	if v <= 0 {
		return def
	}
	return v
}
