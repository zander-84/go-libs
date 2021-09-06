package memory

import (
	"github.com/golang/groupcache/singleflight"
	"github.com/patrickmn/go-cache"
	"github.com/zander-84/go-libs/think"
	"sync"
	"sync/atomic"
	"time"
)

type Memory struct {
	engine       *cache.Cache
	conf         Conf
	once         int64
	err          error
	lock         sync.Mutex
	singleflight singleflight.Group
}

func NewMemory(conf Conf) *Memory {
	this := &Memory{}
	this.init(conf)
	return this
}

func (this *Memory) init(conf Conf) {
	this.conf = conf.SetDefault()
	this.err = think.ErrInstanceUnDone
	atomic.StoreInt64(&this.once, 0)
	this.engine = nil
}

func (this *Memory) Start() error {
	this.lock.Lock()
	defer this.lock.Unlock()

	if atomic.CompareAndSwapInt64(&this.once, 0, 1) {
		this.engine = cache.New(time.Duration(this.conf.Expiration)*time.Minute, time.Duration(this.conf.CleanupInterval)*time.Minute)
		this.err = nil
	}

	return this.err
}

func (this *Memory) Stop() error {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.engine = nil
	atomic.StoreInt64(&this.once, 0)
	this.err = think.ErrInstanceUnDone
	return nil
}

func (this *Memory) Restart(conf Conf) error {
	this.Stop()
	this.init(conf)
	return this.Start()
}

func (this *Memory) Engine() *cache.Cache {
	return this.engine
}
