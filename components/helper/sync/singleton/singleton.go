package singleton

import (
	"context"
	"errors"
	"fmt"
)

//单人工作模式
type Singleton struct {
	jobs chan *Job
	stop chan struct{}
	Data []*Job
}

func NewSingleton(taskNum int) *Singleton {
	return &Singleton{
		jobs: make(chan *Job, taskNum),
		stop: nil,
		Data: make([]*Job, 0),
	}
}

func (thisSingleton *Singleton) AddJob(job *Job) error {
	if cap(thisSingleton.jobs) == len(thisSingleton.jobs) {
		return errors.New("overload")
	}
	select {
	case thisSingleton.jobs <- job:
		return nil
	case <-job.Ctx.Done():
		return job.Ctx.Err()
	}
}
func (thisSingleton *Singleton) Run() {
	for job := range thisSingleton.jobs {
		job.Error = func() (err error) {
			defer func() {
				if err := recover(); err != nil {
					err = fmt.Errorf("panic Handler: %v", err)
				}
			}()

			return job.Handler(job)
		}()
		thisSingleton.Data = append(thisSingleton.Data, job)
	}
	thisSingleton.stop <- struct{}{}
}
func (thisSingleton *Singleton) Shutdown(ctx context.Context) {
	close(thisSingleton.jobs)
	select {
	case <-thisSingleton.stop:
	case <-ctx.Done():
	}
}
