package sd

import (
	"math/rand"
	"sync"
	"sync/atomic"
)

type random struct {
	listener   *Listener
	nodes      []string
	version    uint64
	lock       sync.RWMutex
	recordLock sync.RWMutex
	isRecord   bool
	used       map[string]int64
}

func NewRandom(l *Listener, isRecord bool) Balancer {
	rd := &random{
		listener: l,
		isRecord: isRecord,
		used:     make(map[string]int64),
	}
	rd.Update()
	return rd
}

func (rd *random) Update() {
	rd.lock.Lock()
	defer rd.lock.Unlock()
	if atomic.LoadUint64(&rd.version) == rd.listener.GetVersion() {
		return
	}

	_, addrSlice, version := rd.listener.Get()
	rd.nodes = addrSlice
	atomic.StoreUint64(&rd.version, version)
}

func (rd *random) Next() (string, error) {
	if atomic.LoadUint64(&rd.version) != rd.listener.GetVersion() {
		rd.Update()
	}
	rd.lock.RLock()
	defer rd.lock.RUnlock()

	listenErr := rd.listener.Err()
	if listenErr != nil {
		return "", listenErr
	}

	if len(rd.nodes) <= 0 {
		return "", ErrNoNode
	}

	res := rd.nodes[rand.Intn(len(rd.nodes))]

	// 保存使用数据
	if rd.isRecord {
		rd.record(res)
	}
	return res, nil
}
func (rd *random) Get(uid string) (string, error) {
	return rd.Next()
}

func (rd *random) All() ([]string, error) {
	if atomic.LoadUint64(&rd.version) != rd.listener.GetVersion() {
		rd.Update()
	}
	rd.lock.RLock()
	defer rd.lock.RUnlock()
	if len(rd.nodes) <= 0 {
		return nil, ErrNoNode
	}
	return rd.nodes, nil
}
func (rd *random) Used() map[string]int64 {
	rd.recordLock.RLock()
	defer rd.recordLock.RUnlock()
	return rd.used
}

func (rd *random) record(data string) {
	rd.recordLock.Lock()
	if tmp, ok := rd.used[data]; ok {
		rd.used[data] = tmp + 1
	} else {
		rd.used[data] = 1
	}
	rd.recordLock.Unlock()
}
