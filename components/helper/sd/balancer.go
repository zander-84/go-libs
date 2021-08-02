package sd

import (
	"errors"
)

// ErrNoNode is returned when no qualifying node are available.
var ErrNoNode = errors.New("no node available")
var ErrNotImplemented = errors.New("not implemented")

// Balancer yields endpoints according to some heuristic.
type Balancer interface {
	Update()
	Next() (string, error)
	All() ([]string, error)
	Used() map[string]int64
	Get(uid string) (string, error) //用于hash一致性
}

type Policy int

const (
	RoundRobin Policy = iota
	WeightRoundRobin
	ConsistentHash
	Random
)

func NewBalancer(listener *Listener, record bool, p Policy) Balancer {
	switch p {
	case RoundRobin:
		return NewRoundRobin(listener, record)
	case WeightRoundRobin:
		return NewWeightRoundRobin(listener, record)
	case ConsistentHash:
		return NewConsistentHash(listener, record)
	case Random:
		return NewRandom(listener, record)
	default:
		return NewRoundRobin(listener, record)
	}
}
