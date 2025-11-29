package routercore

import "errors"

type Mapper interface {
	Add(nodes ...string)
	Remove(nodes ...string)
	Pick(key []byte) string
}

type Algo string

const (
	AlgoJump   Algo = "jump"
	AlgoMaglev Algo = "maglev"
	AlgoCHBL   Algo = "chbl"
)

type Options struct {
	TableSize     int
	LoadFactor    float64
	Vnodes        int
	WalkThreshold int
	HashSeed      uint64

	// ExpectedKeys is the expected total number of keys/requests that
	// will be assigned when using CH-BL. It is used to compute the
	// per-node capacity C = ceil(c * ExpectedKeys / numNodes).
	//
	// For Jump and Maglev this field is ignored.
	ExpectedKeys int
}

var ErrUnknownAlgo = errors.New("router: unknown algorithm")
