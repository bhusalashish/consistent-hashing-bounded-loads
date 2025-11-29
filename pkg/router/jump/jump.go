package jump

import (
	"sync"

	"github.com/bhusalashish/consistent-hashing-bounded-loads.git/pkg/hash"
	"github.com/bhusalashish/consistent-hashing-bounded-loads.git/pkg/routercore"
)

const MAGIC_NUMBER = 2862933555777941757

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

	// Compute 64-bit hash using our standard xxhash implementation
	h := hash.XXH64(key, 0) // seed = 0 for Jump (standard practice)

	// Jump Consistent Hash algorithm (Google)
	numBuckets := len(m.nodes)
	b := -1
	j := 0

	for j < numBuckets {
		b = j
		h = h*MAGIC_NUMBER + 1
		j = int(float64(b+1) * (float64(1<<31) / float64((h>>33)+1)))
	}

	return m.nodes[b]
}
