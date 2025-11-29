package router

import (
	"errors"

	chbl "github.com/bhusalashish/consistent-hashing-bounded-loads.git/pkg/router/chbl"
	jump "github.com/bhusalashish/consistent-hashing-bounded-loads.git/pkg/router/jump"
	maglev "github.com/bhusalashish/consistent-hashing-bounded-loads.git/pkg/router/maglev"
	routercore "github.com/bhusalashish/consistent-hashing-bounded-loads.git/pkg/routercore"
)

// Mapper is the common interface for all routing algorithms.
//
// Implementations (Jump, Maglev, CH-BL) must be able to:
//   - Add nodes (e.g., when a new backend comes online),
//   - Remove nodes (e.g., when a backend is drained),
//   - Pick a node for a given key.
//
// This makes it easy to swap algorithms without changing
// the calling code.
type Mapper interface {
	// Add registers one or more nodes with the mapper.
	// Re-adding an existing node MUST be safe (idempotent).
	Add(nodes ...string)

	// Remove unregisters one or more nodes.
	// Removing a node that does not exist MUST be safe (no-op).
	Remove(nodes ...string)

	// Pick returns the node chosen for the given key.
	//
	// Implementations may panic if there are zero nodes;
	// the caller is responsible for ensuring the node set
	// is non-empty before calling Pick.
	Pick(key []byte) string
}

// Algo is an enum-like type for supported algorithms.
type Algo string

const (
	AlgoJump   Algo = "jump"   // Jump consistent hashing (baseline)
	AlgoMaglev Algo = "maglev" // Maglev permutation table
	AlgoCHBL   Algo = "chbl"   // Consistent Hashing with Bounded Loads
)

// Options configures algorithm-specific parameters and shared parameters.
//
// Not all fields are used by all algorithms; unused fields are ignored.
type Options struct {
	// TableSize is the Maglev lookup table size (M).
	// If zero, the Maglev implementation should pick a sensible default.
	TableSize int

	// LoadFactor is the "c" parameter for CH-BL. It bounds the maximum
	// per-node load to at most c * (total_keys / num_nodes).
	// Typical values are in [1.1, 1.5].
	LoadFactor float64

	// Vnodes is the number of virtual nodes per physical node for CH-BL.
	// More vnodes -> smoother distribution but more memory and rebuild cost.
	Vnodes int

	// WalkThreshold is the number of linear steps CH-BL will take on the
	// ring before invoking the two-choice fallback to avoid long walks.
	WalkThreshold int

	// HashSeed is a shared seed for the hash function. Using a fixed seed
	// makes the mapping deterministic across processes for the same
	// node set and keys, which is important for reproducible experiments.
	HashSeed uint64
}

// ErrUnknownAlgo is returned by New when the requested Algo is not supported.
var ErrUnknownAlgo = errors.New("router: unknown algorithm")

// New constructs a Mapper implementation for the given algorithm and options.
//
// The nodes slice is the initial set of backend IDs.
// Implementations MUST treat node IDs as opaque strings but stable identifiers.
func New(algo routercore.Algo, opts routercore.Options, nodes []string) (routercore.Mapper, error) {
	switch algo {
	case routercore.AlgoJump:
		return jump.NewJump(nodes, opts)
	case routercore.AlgoMaglev:
		return maglev.NewMaglev(nodes, opts)
	case routercore.AlgoCHBL:
		return chbl.NewCHBL(nodes, opts)
	default:
		return nil, routercore.ErrUnknownAlgo
	}
}
