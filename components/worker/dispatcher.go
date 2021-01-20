package worker

import (
	"errors"
	"github.com/zander-84/go-libs/components/errs"
	"sync"
	"sync/atomic"
)

type Dispatcher struct {
	worker     []*worker     // 工人
	workerPool chan chan Job // 真正任务派发队列
	jobChannel chan Job
	conf       Conf
	err        error
	lock       sync.Mutex
	once       int64
	status     bool
	tmp        int64
}

func NewDispatcher(conf Conf) *Dispatcher {
	var this = new(Dispatcher)
	this.init(conf)
	return this
}
func (this *Dispatcher) init(conf Conf) {
	this.conf = conf.SetDefault()
	this.once = 0
	this.err = errs.UninitializedError
}

func (this *Dispatcher) Start() error {
	this.lock.Lock()
	defer this.lock.Unlock()

	atomic.AddInt64(&this.once, 1)
	if atomic.LoadInt64(&this.once) != 1 {
		atomic.StoreInt64(&this.once, 2)
		return this.err
	}

	this.worker = make([]*worker, this.conf.MaxWorkers)
	this.workerPool = make(chan chan Job, this.conf.MaxWorkers)
	this.jobChannel = make(chan Job, this.conf.MaxQueues)
	this.run()
	this.err = nil
	return this.err
}

func (this *Dispatcher) Stop() {
	this.lock.Lock()
	defer this.lock.Unlock()
	for _, w := range this.worker {
		w.stop()
	}
	this.once = 0
}

func (this *Dispatcher) run() {
	for i := 0; i < len(this.worker); i++ {
		this.worker[i] = newWorker(this.workerPool)
		this.worker[i].start(i)
	}
	go this.dispatch()
}

func (this *Dispatcher) dispatch() {
	for {
		select {
		case job := <-this.jobChannel:
			jobChannel := <-this.workerPool
			jobChannel <- job
		}
	}
}

func (this *Dispatcher) AddJob(job Job) error {
	if this.conf.MaxQueues <= len(this.jobChannel) {
		return errors.New("system busyness")
	}
	this.jobChannel <- job
	return nil
}
