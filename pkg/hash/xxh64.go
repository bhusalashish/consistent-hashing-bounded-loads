package hash

import (
	"github.com/cespare/xxhash/v2"
)

func XXH64(b []byte, seed uint64) uint64 {
	h := xxhash.New()
	// Simple pattern: write seed then bytes
	var buf [8]byte
	for i := 0; i < 8; i++ {
		buf[i] = byte(seed >> (8 * i))
	}
	h.Write(buf[:])
	h.Write(b)
	return h.Sum64()
}
