package hash

import "testing"

func TestXXH64Deterministic(t *testing.T) {
	h1 := XXH64([]byte("hello"), 42)
	h2 := XXH64([]byte("hello"), 42)
	if h1 != h2 {
		t.Fatalf("expected deterministic hash but got %d and %d", h1, h2)
	}
}

func TestXXH64DifferentSeeds(t *testing.T) {
	h1 := XXH64([]byte("hello"), 1)
	h2 := XXH64([]byte("hello"), 2)
	if h1 == h2 {
		t.Fatalf("expected different seeds to produce different hashes")
	}
}
