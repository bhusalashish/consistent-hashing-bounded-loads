package hash

import "testing"

func TestXXH64Deterministic(t *testing.T) {
	a := XXH64([]byte("hello"), 123)
	b := XXH64([]byte("hello"), 123)
	if a != b {
		t.Fatalf("expected deterministic hash, got %d and %d", a, b)
	}
}
