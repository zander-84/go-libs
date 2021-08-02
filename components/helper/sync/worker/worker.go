package worker

import (
	"context"
	"sync"
	"sync/atomic"
)

// Worker 此模型具体实现 缺点：比较重，优点：能够精确控制整个时间
type Worker struct {
	jobs       []*Job        // 事件中转
	maxWorkers chan struct{} // 最大线程数

	isErrQuit   bool
	errQuitChan chan error
	onceErrQuit int32
}

func NewWorker() *Worker {
	return newWorker(false)
}

// NewFastWorker 遇到错误就退出
func NewFastWorker() *Worker {
	return newWorker(true)
}

func newWorker(errQuit bool) *Worker {
	thisWorker := new(Worker)
	thisWorker.jobs = make([]*Job, 0)

	thisWorker.errQuitChan = make(chan error, 1)
	thisWorker.isErrQuit = errQuit
	thisWorker.onceErrQuit = 0

	return thisWorker
}

func (w *Worker) doErrQuite(err error) {
	if w.isErrQuit && err != nil && atomic.CompareAndSwapInt32(&w.onceErrQuit, 0, 1) {
		w.errQuitChan <- err
	}
}

func (w *Worker) hasErrQuite() bool {
	if w.isErrQuit && atomic.LoadInt32(&w.onceErrQuit) > 0 {
		return true
	}
	return false
}

func (w *Worker) Do(ctx context.Context, maxWorkers int) (*Jobs, error) {
	if maxWorkers < 1 {
		maxWorkers = 1
	}
	w.maxWorkers = make(chan struct{}, maxWorkers)

	var fin = make(chan struct{}, 1)
	var res = &Jobs{}

	go func() {
		w.do()
		fin <- struct{}{}
	}()

	done := ctx.Done()
	if done == nil {
		select {
		case quitErr := <-w.errQuitChan:
			return nil, quitErr
		case <-fin:
			res.dataSlice = w.jobs
			return res, nil
		}
	} else {
		select {
		case quitErr := <-w.errQuitChan:
			return nil, quitErr
		case <-done:
			return res, ctx.Err()
		case <-fin:
			res.dataSlice = w.jobs
			return res, nil
		}
	}
}

// AddJob 为什么没有用chan接受数据，是为了减少goroutine创建，返回值的使用必须等Do执行完毕后
func (w *Worker) AddJob(ctx context.Context, title string, in interface{}, f func(in interface{}) (interface{}, error)) *Job {
	j := newJob(ctx, title, in, f)
	w.jobs = append(w.jobs, j)
	return j
}

//func (thisWorker *Worker) fanInOut(in []*Event) <-chan *Event {
//	eventsChan := make(chan *Event)
//	go func() {
//		defer close(eventsChan)
//		for _, e := range in {
//			eventsChan <- e
//		}
//	}()
//	return eventsChan
//}

func (w *Worker) do() {
	var wg sync.WaitGroup
	//for e := range thisWorker.fanInOut(thisWorker.eventsSlice) {
	for _, j := range w.jobs {
		w.maxWorkers <- struct{}{}
		wg.Add(1)
		go func(j *Job) {
			defer func() {
				wg.Done()
				<-w.maxWorkers
			}()
			w.consumer(j)
		}(j)
	}

	wg.Wait()
}

func (w *Worker) consumer(j *Job) *Job {
	if w.hasErrQuite() {
		return j
	}

	done := j.ctx.Done()
	if done == nil {
		j.output, j.error = j.handler(j.input)
		w.doErrQuite(j.error)
		return j
	}

	fin := make(chan error, 1)

	// 引入 output err 局部变量为了去锁
	var output interface{}
	var err error

	go func() {
		output, err = j.handler(j.input)
		fin <- err
	}()

	select {
	case <-done:
		j.error = j.ctx.Err()
		w.doErrQuite(j.error)
		return j
	case <-fin:
		j.output, j.error = output, err
		w.doErrQuite(j.error)
		return j
	}
}
