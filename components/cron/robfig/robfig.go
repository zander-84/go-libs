package robfig

import (
	"github.com/robfig/cron/v3"
	cron2 "github.com/zander-84/go-libs/components/cron"
	"github.com/zander-84/go-libs/components/errs"
	"github.com/zander-84/go-libs/components/helper"
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
	this.time = helper.NewTime(this.conf.TimeZone)
	this.err = errs.UninitializedError
	this.jobs = make(map[string]*cron2.Job)
	this.once = 0
}

func (this *Robfig) Start() error {
	this.lock.Lock()
	defer this.lock.Unlock()

	atomic.AddInt64(&this.once, 1)
	if atomic.LoadInt64(&this.once) != 1 {
		atomic.StoreInt64(&this.once, 2)
		return this.err
	}

	this.engine = cron.New(cron.WithLocation(this.time.Location()), cron.WithSeconds())
	this.err = nil
	return nil
}

func (this *Robfig) Stop() {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.engine != nil {
		this.engine.Stop()
	}

	this.engine = nil
	this.once = 0
	this.err = errs.UninitializedError
	this.jobs = make(map[string]*cron2.Job)
}

func (this *Robfig) Restart(conf Conf) error {
	this.Stop()
	this.init(conf)
	return this.Start()
}

func (this *Robfig) Engine() *cron.Cron {
	return this.engine
}
