package sd

import (
	"sync"
	"sync/atomic"
)

type roundRobin struct {
	listener     *Listener
	nodes        []string
	currentIndex uint64
	version      uint64
	indexMax     uint64
	indexLock    sync.RWMutex
	lock         sync.RWMutex
	recordLock   sync.RWMutex
	isRecord     bool
	used         map[string]int64
}

// NewRoundRobin returns a load balancer that returns services in sequence.
func NewRoundRobin(l *Listener, isRecord bool) Balancer {
	rr := &roundRobin{
		listener:     l,
		currentIndex: 0,
		indexMax:     100000000,
		isRecord:     isRecord,
		used:         make(map[string]int64),
	}
	rr.Update()
	return rr
}

func (rr *roundRobin) Update() {
	rr.lock.Lock()
	defer rr.lock.Unlock()
	if atomic.LoadUint64(&rr.version) == rr.listener.GetVersion() {
		return
	}

	_, addrSlice, version := rr.listener.Get()
	rr.nodes = addrSlice
	atomic.StoreUint64(&rr.version, version)
}

func (rr *roundRobin) Next() (string, error) {
	if atomic.LoadUint64(&rr.version) != rr.listener.GetVersion() {
		rr.Update()
	}
	rr.lock.RLock()
	defer rr.lock.RUnlock()
	if len(rr.nodes) <= 0 {
		return "", ErrNoNode
	}

	old := atomic.AddUint64(&rr.currentIndex, 1) - 1
	// once reset
	if old > rr.indexMax {
		rr.indexLock.Lock()
		if atomic.LoadUint64(&rr.currentIndex) > rr.indexMax {
			atomic.StoreUint64(&rr.currentIndex, 0)
		}
		rr.indexLock.Unlock()
	}

	idx := old % uint64(len(rr.nodes))

	// 保存使用数据
	if rr.isRecord {
		rr.record(rr.nodes[idx])
	}
	return rr.nodes[idx], nil
}

func (rr *roundRobin) Get(uid string) (string, error) {
	return rr.Next()
}

func (rr *roundRobin) All() ([]string, error) {
	if atomic.LoadUint64(&rr.version) != rr.listener.GetVersion() {
		rr.Update()
	}
	rr.lock.RLock()
	defer rr.lock.RUnlock()
	if len(rr.nodes) <= 0 {
		return nil, ErrNoNode
	}
	return rr.nodes, nil
}

func (rr *roundRobin) Used() map[string]int64 {
	rr.recordLock.RLock()
	defer rr.recordLock.RUnlock()
	return rr.used
}

func (rr *roundRobin) record(data string) {
	rr.recordLock.Lock()
	if tmp, ok := rr.used[data]; ok {
		rr.used[data] = tmp + 1
	} else {
		rr.used[data] = 1
	}
	rr.recordLock.Unlock()
}
