package sd

import (
	"fmt"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type Hash func(data []byte) uint32

type UInt32Slice []uint32

func (s UInt32Slice) Len() int {
	return len(s)
}

func (s UInt32Slice) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s UInt32Slice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type consistentHash struct {
	listener   *Listener
	lock       sync.RWMutex
	hash       Hash
	replicas   int               //复制因子
	keys       UInt32Slice       //已排序的节点hash切片
	hashMap    map[uint32]string //节点哈希和Key的map,键是hash值，值是节点key
	recordLock sync.RWMutex
	isRecord   bool
	used       map[string]int64
	version    uint64
}

var DefaultHash = crc32.ChecksumIEEE

func NewConsistentHash(listener *Listener, isRecord bool) *consistentHash {
	m := &consistentHash{
		listener: listener,
		replicas: 10000,
		hashMap:  make(map[uint32]string),
		keys:     make(UInt32Slice, 0),
		isRecord: isRecord,
		used:     make(map[string]int64),
	}

	return m
}

func (ch *consistentHash) Update() {
	ch.lock.Lock()
	defer ch.lock.Unlock()
	if atomic.LoadUint64(&ch.version) == ch.listener.GetVersion() {
		return
	}

	_, addrSlice, version := ch.listener.Get()
	hashMap := make(map[uint32]string)
	keys := make(UInt32Slice, 0)
	for _, node := range addrSlice {
		for i := 0; i < ch.replicas; i++ {
			hash := DefaultHash([]byte(strconv.Itoa(i) + node))
			keys = append(keys, hash)
			hashMap[hash] = node
		}
	}
	sort.Sort(keys)
	ch.keys = keys
	ch.hashMap = hashMap
	atomic.StoreUint64(&ch.version, version)
}

func (ch *consistentHash) Next() (string, error) {
	return ch.Get(fmt.Sprintf("%d", time.Now().UnixNano()))
}

func (ch *consistentHash) Get(uid string) (string, error) {
	if atomic.LoadUint64(&ch.version) != ch.listener.GetVersion() {
		ch.Update()
	}
	ch.lock.RLock()
	defer ch.lock.RUnlock()

	listenErr := ch.listener.Err()
	if listenErr != nil {
		return "", listenErr
	}

	if len(ch.hashMap) <= 0 {
		return "", ErrNoNode
	}
	hash := DefaultHash([]byte(uid))
	// 通过二分查找获取最优节点，第一个"服务器hash"值大于"数据hash"值的就是最优"服务器节点"
	idx := sort.Search(len(ch.keys), func(i int) bool { return ch.keys[i] >= hash })
	// 如果查找结果 大于 服务器节点哈希数组的最大索引，表示此时该对象哈希值位于最后一个节点之后，那么放入第一个节点中
	if idx == len(ch.keys) {
		idx = 0
	}

	res := ch.hashMap[ch.keys[idx]]
	if ch.isRecord {
		ch.record(res)
	}

	return res, nil
}

func (ch *consistentHash) record(data string) {
	ch.recordLock.Lock()
	if tmp, ok := ch.used[data]; ok {
		ch.used[data] = tmp + 1
	} else {
		ch.used[data] = 1
	}
	ch.recordLock.Unlock()
}

func (ch *consistentHash) All() ([]string, error) {
	if atomic.LoadUint64(&ch.version) != ch.listener.GetVersion() {
		ch.Update()
	}
	ch.lock.RLock()
	defer ch.lock.RUnlock()
	if len(ch.hashMap) <= 0 {
		return nil, ErrNoNode
	}

	var res = make([]string, 0)
	for _, n := range ch.hashMap {
		res = append(res, n)
	}

	return res, nil
}

func (ch *consistentHash) Used() map[string]int64 {
	ch.recordLock.RLock()
	defer ch.recordLock.RUnlock()
	return ch.used
}
