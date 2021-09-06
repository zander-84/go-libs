package worker

import (
	"context"
	"errors"
	"github.com/zander-84/go-libs/think"
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
	exit       chan struct{}
}

func NewDispatcher(conf Conf) *Dispatcher {
	var this = new(Dispatcher)
	this.init(conf)
	return this
}

func (this *Dispatcher) init(conf Conf) {
	this.conf = conf.SetDefault()
	atomic.StoreInt64(&this.once, 0)
	this.err = think.ErrInstanceUnDone
}

func (this *Dispatcher) Start() error {
	this.lock.Lock()
	defer this.lock.Unlock()

	if atomic.CompareAndSwapInt64(&this.once, 0, 1) {
		this.worker = make([]*worker, this.conf.MaxWorkers)
		this.workerPool = make(chan chan Job, this.conf.MaxWorkers)
		this.jobChannel = make(chan Job, this.conf.MaxQueues)
		this.exit = make(chan struct{}, 1)
		this.run()
		this.err = nil
	}

	return this.err
}

func (this *Dispatcher) Stop(ctx context.Context) error {
	this.lock.Lock()
	defer this.lock.Unlock()
	done := ctx.Done()
	if done == nil {
		this.stop()
		return nil
	}
	fin := make(chan struct{}, 1)
	go func() {
		this.stop()
		fin <- struct{}{}
	}()
	select {
	case <-done:
		return ctx.Err()
	case <-fin:
		return nil
	}
}
func (this *Dispatcher) stop() {
	if atomic.LoadInt64(&this.once) == 0 {
		return
	}

	close(this.jobChannel)
	<-this.exit

	for _, w := range this.worker {
		w.stop()
	}
	atomic.StoreInt64(&this.once, 0)
}

func (this *Dispatcher) run() {
	for i := 0; i < len(this.worker); i++ {
		this.worker[i] = newWorker(this.workerPool)
		this.worker[i].start(i)
	}
	go this.dispatch()
}

func (this *Dispatcher) dispatch() {
	for job := range this.jobChannel {
		jobChannel := <-this.workerPool
		jobChannel <- job
	}
	this.exit <- struct{}{}
}

func (this *Dispatcher) AddJob(job Job) (err error) {
	defer func() {
		if rerr := recover(); rerr != nil {
			err = errors.New("queue already exited")
		}
	}()

	if this.conf.MaxQueues <= len(this.jobChannel) {
		return errors.New("system busyness")
	}
	this.jobChannel <- job
	return nil
}
