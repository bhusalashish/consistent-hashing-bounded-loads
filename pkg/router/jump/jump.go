package jump

import (
	"sync"

	"github.com/bhusalashish/consistent-hashing-bounded-loads.git/pkg/router"
)

// mapper is a placeholder implementation to satisfy router.Mapper.
// We will replace this with a real Jump consistent hashing implementation.
type mapper struct {
	mu    sync.RWMutex
	nodes []string
}

func newJump(nodes []string, opts router.Options) (router.Mapper, error) {
	m := &mapper{}
	m.Add(nodes...)
	return m, nil
}

func (m *mapper) Add(nodes ...string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.nodes = append(m.nodes, nodes...)
}

func (m *mapper) Remove(nodes ...string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(nodes) == 0 || len(m.nodes) == 0 {
		return
	}
	rem := make(map[string]struct{}, len(nodes))
	for _, n := range nodes {
		rem[n] = struct{}{}
	}
	var out []string
	for _, n := range m.nodes {
		if _, ok := rem[n]; !ok {
			out = append(out, n)
		}
	}
	m.nodes = out
}

func (m *mapper) Pick(key []byte) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	// TEMPORARY stub: always return first node.
	// We will replace this logic with real Jump hashing in a later step.
	if len(m.nodes) == 0 {
		panic("jump: no nodes registered")
	}
	return m.nodes[0]
}
