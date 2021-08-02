package sd

import (
	"sync"
	"sync/atomic"
)

type weightRoundRobin struct {
	listener *Listener
	version  uint64
	lock     sync.RWMutex

	currentIndex uint64
	nodes        []*weightNode

	recordLock sync.RWMutex
	isRecord   bool
	used       map[string]int64
}

type weightNode struct {
	node            string
	weight          int //权重值
	currentWeight   int //节点当前权重
	effectiveWeight int //有效权重
}

func newWeightNode(node string, weight int) *weightNode {
	return &weightNode{
		node:            node,
		weight:          weight,
		currentWeight:   0,
		effectiveWeight: 0,
	}
}
func NewWeightRoundRobin(l *Listener, isRecord bool) Balancer {
	wrr := &weightRoundRobin{
		listener:     l,
		nodes:        make([]*weightNode, 0),
		currentIndex: 0,
		isRecord:     isRecord,
		used:         make(map[string]int64),
	}
	return wrr
}

func (wrr *weightRoundRobin) Update() {
	wrr.lock.Lock()
	defer wrr.lock.Unlock()
	if atomic.LoadUint64(&wrr.version) == wrr.listener.GetVersion() {
		return
	}
	ns, _, version := wrr.listener.Get()
	tmp := make([]*weightNode, 0)
	for n, w := range ns {
		tmp = append(tmp, newWeightNode(n, w))
	}
	wrr.nodes = tmp
	atomic.StoreUint64(&wrr.version, version)
}

func (wrr *weightRoundRobin) Next() (string, error) {
	if atomic.LoadUint64(&wrr.version) != wrr.listener.GetVersion() {
		wrr.Update()
	}
	wrr.lock.Lock()
	defer wrr.lock.Unlock()

	if len(wrr.nodes) <= 0 {
		return "", ErrNoNode
	}

	total := 0
	var best *weightNode

	for i := 0; i < len(wrr.nodes); i++ {
		w := wrr.nodes[i]
		//step 1 统计所有有效权重之和
		total += w.effectiveWeight

		//step 2 变更节点临时权重为的节点临时权重+节点有效权重
		w.currentWeight += w.effectiveWeight

		//step 3 有效权重默认与权重相同，通讯异常时-1, 通讯成功+1，直到恢复到weight大小
		if w.effectiveWeight < w.weight {
			w.effectiveWeight++
		}
		//step 4 选择最大临时权重点节点
		if best == nil || w.currentWeight > best.currentWeight {
			best = w
		}
	}

	if best == nil {
		return "", ErrNoNode
	}
	best.currentWeight -= total
	if wrr.isRecord {
		wrr.record(best.node)
	}
	return best.node, nil
}

func (wrr *weightRoundRobin) Get(uid string) (string, error) {
	return wrr.Next()
}

func (wrr *weightRoundRobin) All() ([]string, error) {
	if atomic.LoadUint64(&wrr.version) != wrr.listener.GetVersion() {
		wrr.Update()
	}

	wrr.lock.RLock()
	defer wrr.lock.RUnlock()

	if len(wrr.nodes) <= 0 {
		return nil, ErrNoNode
	}

	var res = make([]string, 0)
	for _, n := range wrr.nodes {
		res = append(res, n.node)
	}

	return res, nil
}

func (wrr *weightRoundRobin) Used() map[string]int64 {
	wrr.recordLock.RLock()
	defer wrr.recordLock.RUnlock()
	return wrr.used
}

func (wrr *weightRoundRobin) record(data string) {
	wrr.recordLock.Lock()
	if tmp, ok := wrr.used[data]; ok {
		wrr.used[data] = tmp + 1
	} else {
		wrr.used[data] = 1
	}
	wrr.recordLock.Unlock()
}
