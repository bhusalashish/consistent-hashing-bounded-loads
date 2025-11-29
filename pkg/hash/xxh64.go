package hash

import (
	"encoding/binary"
	"github.com/cespare/xxhash/v2"
)

// XXH64 hashes arbitrary bytes using xxHash64 with a custom seed.
//
// Why we combine seed + key manually:
//   - xxHash64 itself does not take a seed in the constructor.
//   - Prepending the seed as 8 bytes to the input is a common,
//     production-safe pattern to ensure deterministic variation.
//
// Why 64-bit:
//   - All our routing algorithms (Jump, Maglev, CH-BL) benefit from
//     a uniformly distributed 64-bit space.
func XXH64(data []byte, seed uint64) uint64 {
	b := make([]byte, 8+len(data))
	binary.LittleEndian.PutUint64(b[:8], seed)
	copy(b[8:], data)
	return xxhash.Sum64(b)
}

// XXH64String is a convenience wrapper around XXH64 for string keys.
// Using this avoids repeated []byte(key) allocations all over the code.
func XXH64String(s string, seed uint64) uint64 {
	b := make([]byte, 8+len(s))
	binary.LittleEndian.PutUint64(b[:8], seed)
	copy(b[8:], s)
	return xxhash.Sum64(b)
}
