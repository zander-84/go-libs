package goredis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/golang/groupcache/singleflight"
	"github.com/zander-84/go-libs/components/errs"
	"sync"
	"sync/atomic"
	"time"
)

type Rdb struct {
	engine       *redis.Client
	conf         Conf
	once         int64
	err          error
	lock         sync.Mutex
	context      context.Context
	singleflight singleflight.Group
}

func NewRdb(conf Conf) *Rdb {
	this := new(Rdb)
	this.init(conf)
	return this
}

func (this *Rdb) init(conf Conf) {
	this.conf = conf.SetDefault()
	this.err = errs.UninitializedError
	this.context = context.Background()
	this.engine = nil
	this.once = 0
}

func (this *Rdb) Start() error {
	this.lock.Lock()
	defer this.lock.Unlock()

	atomic.AddInt64(&this.once, 1)
	if atomic.LoadInt64(&this.once) != 1 {
		atomic.StoreInt64(&this.once, 2)
		return this.err
	}

	this.engine = redis.NewClient(&redis.Options{
		Addr:         this.conf.Addr,
		Password:     this.conf.Password,
		DB:           this.conf.Db,
		PoolSize:     this.conf.PoolSize,
		IdleTimeout:  time.Duration(this.conf.IdleTimeout) * time.Second,
		MinIdleConns: this.conf.MinIdle,
	})

	this.err = this.engine.Ping(context.Background()).Err()
	return this.err
}

func (this *Rdb) Stop() {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.engine != nil {
		_ = this.engine.Close()
	}
	this.engine = nil
	this.once = 0
	this.err = errs.UninitializedError
}

func (this *Rdb) Restart(conf Conf) error {
	this.Stop()
	this.init(conf)
	return this.Start()
}

func (this *Rdb) Engine() *redis.Client {
	return this.engine
}
