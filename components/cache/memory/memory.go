package memory

import (
	"github.com/patrickmn/go-cache"
	"github.com/zander-84/go-libs/components/errs"
	"sync"
	"sync/atomic"
	"time"
)

type Memory struct {
	engine *cache.Cache
	conf   Conf
	once   int64
	err    error
	lock   sync.Mutex
}

func NewMemory(conf Conf) *Memory {
	this := &Memory{}
	this.init(conf)
	return this
}

func (this *Memory) init(conf Conf) {
	this.conf = conf.SetDefault()
	this.err = errs.UninitializedError
	this.once = 0
	this.engine = nil
}

func (this *Memory) Start() error {
	this.lock.Lock()
	defer this.lock.Unlock()

	atomic.AddInt64(&this.once, 1)
	if this.once != 1 {
		return this.err
	}
	this.engine = cache.New(time.Duration(this.conf.Expiration)*time.Minute, time.Duration(this.conf.CleanupInterval)*time.Minute)
	this.err = nil
	return nil
}
func (this *Memory) Stop() {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.engine = nil
	this.once = 0
	this.err = errs.UninitializedError
}

func (this *Memory) Restart(conf Conf) error {
	this.Stop()
	this.init(conf)
	return this.Start()
}

func (this *Memory) Engine() *cache.Cache {
	return this.engine
}
