package maglev

import (
	"testing"

	"github.com/bhusalashish/consistent-hashing-bounded-loads.git/pkg/routercore"
)

func TestMaglevDeterministic(t *testing.T) {
	nodes := []string{"n1", "n2", "n3", "n4"}

	m1, _ := NewMaglev(nodes, routercore.Options{
		TableSize: 65537,
		HashSeed:  42,
	})

	m2, _ := NewMaglev(nodes, routercore.Options{
		TableSize: 65537,
		HashSeed:  42,
	})

	keys := [][]byte{
		[]byte("user-1"),
		[]byte("user-2"),
		[]byte("user-3"),
		[]byte("user-123456"),
	}

	for _, k := range keys {
		if m1.Pick(k) != m2.Pick(k) {
			t.Fatalf("expected deterministic mapping for key %q", string(k))
		}
	}
}

func TestMaglevUsesAllNodes(t *testing.T) {
	nodes := []string{"n1", "n2", "n3", "n4", "n5"}
	m, _ := NewMaglev(nodes, routercore.Options{
		TableSize: 65537,
		HashSeed:  1,
	})

	counts := make(map[string]int)
	total := 10000

	for i := 0; i < total; i++ {
		key := []byte("k-" + string(rune(i)))
		n := m.Pick(key)
		counts[n]++
	}

	for _, n := range nodes {
		if counts[n] == 0 {
			t.Fatalf("node %s never received any keys", n)
		}
	}
}
