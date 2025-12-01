package ringch

import (
	"testing"

	rc "github.com/bhusalashish/consistent-hashing-bounded-loads.git/pkg/routercore"
)

func TestRingCHDeterministic(t *testing.T) {
	nodes := []string{"n1", "n2", "n3"}

	m, err := NewRingCH(nodes, rc.Options{HashSeed: 42, Vnodes: 50})
	if err != nil {
		t.Fatalf("failed to create ringch mapper: %v", err)
	}

	key := []byte("hello-world")
	r1 := m.Pick(key)
	r2 := m.Pick(key)

	if r1 != r2 {
		t.Fatalf("expected deterministic mapping, got %s vs %s", r1, r2)
	}
}
