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
}

var ErrUnknownAlgo = errors.New("router: unknown algorithm")
