package jump

import (
	"testing"

	"github.com/bhusalashish/consistent-hashing-bounded-loads.git/pkg/routercore"
)

func TestJumpDeterministic(t *testing.T) {
	nodes := []string{"A", "B", "C", "D"}

	m, _ := NewJump(nodes, routercore.Options{})

	k := []byte("user-123")
	r1 := m.Pick(k)
	r2 := m.Pick(k)

	if r1 != r2 {
		t.Fatalf("expected deterministic result, got %s vs %s", r1, r2)
	}
}

func TestJumpChangesMinimalOnAdd(t *testing.T) {
	nodes := []string{"A", "B", "C"}

	m1, _ := NewJump(nodes, routercore.Options{})
	m2, _ := NewJump([]string{"A", "B", "C", "D"}, routercore.Options{})

	moved := 0
	total := 10000

	for i := 0; i < total; i++ {
		key := []byte("k-" + string(rune(i)))
		if m1.Pick(key) != m2.Pick(key) {
			moved++
		}
	}

	ratio := float64(moved) / float64(total)
	if ratio > 0.40 { // should be around ~1/4 = 0.25 ideally
		t.Fatalf("too many keys moved: %.2f", ratio)
	}
}
