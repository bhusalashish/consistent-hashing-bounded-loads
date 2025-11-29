package maglev

import (
	"sync"

	"github.com/bhusalashish/consistent-hashing-bounded-loads.git/pkg/routercore"
)

type mapper struct {
	mu    sync.RWMutex
	nodes []string
	// table and other fields will be added later
}

func NewMaglev(nodes []string, opts routercore.Options) (routercore.Mapper, error) {
	m := &mapper{}
	m.Add(nodes...)
	return m, nil
}

func (m *mapper) Add(nodes ...string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.nodes = append(m.nodes, nodes...)
	// Later: rebuild the Maglev table here.
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
	// Later: rebuild the Maglev table here.
}

func (m *mapper) Pick(key []byte) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	// TEMPORARY stub: always return first node.
	if len(m.nodes) == 0 {
		panic("maglev: no nodes registered")
	}
	return m.nodes[0]
}
