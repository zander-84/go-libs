package robfig

import (
	"github.com/robfig/cron/v3"
	cron2 "github.com/zander-84/go-libs/components/cron"
	"github.com/zander-84/go-libs/components/helper"
	"github.com/zander-84/go-libs/think"
	"sync"
	"sync/atomic"
)

type Robfig struct {
	engine *cron.Cron
	conf   Conf
	jobs   map[string]*cron2.Job
	err    error
	lock   sync.Mutex
	once   int64
	time   *helper.Time
}

func NewRobfig(conf Conf) *Robfig {
	this := new(Robfig)
	this.init(conf)
	return this
}

func (this *Robfig) init(conf Conf) {
	this.conf = conf.SetDefault()
	this.time = helper.NewTime()
	this.err = think.ErrInstanceUnDone
	this.jobs = make(map[string]*cron2.Job)
	atomic.StoreInt64(&this.once, 0)
}

func (this *Robfig) Start() error {
	this.lock.Lock()
	defer this.lock.Unlock()
	if atomic.CompareAndSwapInt64(&this.once, 0, 1) {
		this.engine = cron.New(cron.WithLocation(this.time.Location()), cron.WithSeconds())
		this.err = nil
	}
	return this.err
}

func (this *Robfig) Stop() error {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.engine != nil {
		this.engine.Stop()
	}

	this.engine = nil
	atomic.StoreInt64(&this.once, 0)
	this.err = think.ErrInstanceUnDone
	this.jobs = make(map[string]*cron2.Job)
	return nil
}

func (this *Robfig) Restart(conf Conf) error {
	this.Stop()
	this.init(conf)
	return this.Start()
}

func (this *Robfig) Engine() *cron.Cron {
	return this.engine
}
