package jump

import (
	"sync"

	"github.com/bhusalashish/consistent-hashing-bounded-loads.git/pkg/routercore"
)

type mapper struct {
	mu    sync.RWMutex
	nodes []string
}

func NewJump(nodes []string, opts routercore.Options) (routercore.Mapper, error) {
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
	rem := make(map[string]struct{})
	for _, n := range nodes {
		rem[n] = struct{}{}
	}
	var out []string
	for _, n := range m.nodes {
		if _, exists := rem[n]; !exists {
			out = append(out, n)
		}
	}
	m.nodes = out
}

func (m *mapper) Pick(key []byte) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if len(m.nodes) == 0 {
		panic("jump: no nodes registered")
	}
	return m.nodes[0] // STUB: will replace with real Jump hashing later
}
