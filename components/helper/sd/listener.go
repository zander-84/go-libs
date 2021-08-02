package sd

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

type Listener struct {
	name       string
	data       map[string]int
	version    uint64
	lock       sync.RWMutex
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewListener(name string) *Listener {
	l := new(Listener)
	l.version = 0
	l.data = make(map[string]int)
	l.name = name
	l.ctx, l.cancelFunc = context.WithCancel(context.Background())
	return l
}
func (l *Listener) Name() string {
	return l.name
}

func (l *Listener) Close() {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.cancelFunc()
	atomic.StoreUint64(&l.version, 0)
	l.data = make(map[string]int)

}

func (l *Listener) Context() context.Context {
	return l.ctx
}

func (l *Listener) Exist(addr string) bool {
	l.lock.RLock()
	defer l.lock.RUnlock()
	_, ok := l.data[addr]
	return ok
}

// 全量设置
func (l *Listener) Set(data map[string]int) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.data = data
	atomic.AddUint64(&l.version, 1)
}

// 增量添加
func (l *Listener) Add(addr string) {
	l.AddWeight(addr, 1)
}

func (l *Listener) AddWeight(addr string, weight int) {
	l.lock.Lock()
	defer l.lock.Unlock()
	if weight == 0 {
		weight = 1
	}
	l.data[addr] = weight
	atomic.AddUint64(&l.version, 1)
}

func (l *Listener) Remove(addr string) {
	l.lock.Lock()
	defer l.lock.Unlock()
	delete(l.data, addr)
	atomic.AddUint64(&l.version, 1)
}

func (l *Listener) GetVersion() uint64 {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return atomic.LoadUint64(&l.version)
}

func (l *Listener) Get() (map[string]int, []string, uint64) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	var dataMap = make(map[string]int, 0)
	var dataSlice = make([]string, 0)
	for k, v := range l.data {
		dataMap[k] = v
		dataSlice = append(dataSlice, k)
	}

	return dataMap, dataSlice, atomic.LoadUint64(&l.version)
}

func (l *Listener) SpiltByHyphen(rowAddr string) (addr string, weight int) {
	tmpArr := strings.Split(rowAddr, "-")
	addr = tmpArr[0]

	if len(tmpArr) > 1 {
		n, err := strconv.Atoi(tmpArr[1])
		if err != nil {
			weight = 1
		} else {
			weight = n
		}
	}

	if weight < 1 {
		weight = 1
	}
	return addr, weight
}

func (l *Listener) Println() {
	addrs, _, version := l.Get()
	for k, v := range addrs {
		fmt.Printf("地址:%s  权重:%d  版本:%d\n", k, v, version)
	}
	fmt.Println()
}
