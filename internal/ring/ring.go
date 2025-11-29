package ring

import (
	"fmt"
	"sort"

	"github.com/bhusalashish/consistent-hashing-bounded-loads.git/pkg/hash"
)

// Token represents a point on the hash ring belonging to a particular node.
type Token struct {
	H       uint64 // hash position on the ring
	NodeIdx int    // index into Ring.Nodes
}

// Ring is a vnode-based consistent hash ring.
type Ring struct {
	Tokens []Token  // sorted by H
	Nodes  []string // node index -> node ID
}

// New constructs a ring from the given node IDs and vnode count.
//
// Each node gets 'vnodes' virtual tokens placed around the ring.
// 'seed' is used to produce deterministic token positions.
func New(nodes []string, vnodes int, seed uint64) *Ring {
	if vnodes <= 0 {
		panic("ring: vnodes must be > 0")
	}
	if len(nodes) == 0 {
		return &Ring{}
	}

	r := &Ring{
		Nodes: append([]string(nil), nodes...),
	}

	var tokens []Token
	for i, id := range r.Nodes {
		for v := 0; v < vnodes; v++ {
			key := []byte(fmt.Sprintf("%s#%d-%d", id, v, seed))
			h := hash.XXH64(key, seed)
			tokens = append(tokens, Token{
				H:       h,
				NodeIdx: i,
			})
		}
	}

	sort.Slice(tokens, func(i, j int) bool {
		return tokens[i].H < tokens[j].H
	})

	r.Tokens = tokens
	return r
}

// SuccessorIndex returns the index in r.Tokens of the first token
// whose H >= h, wrapping to 0 if necessary.
//
// Panics if the ring is empty.
func (r *Ring) SuccessorIndex(h uint64) int {
	if len(r.Tokens) == 0 {
		panic("ring: SuccessorIndex on empty ring")
	}
	i := sort.Search(len(r.Tokens), func(i int) bool {
		return r.Tokens[i].H >= h
	})
	if i == len(r.Tokens) {
		return 0
	}
	return i
}
