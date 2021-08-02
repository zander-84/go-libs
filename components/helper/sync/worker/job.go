package worker

import (
	"context"
	"errors"
	"sync"
)

var ErrUndone = errors.New("event undone")
var ErrRecordNotFound = errors.New("event not found")

type Jobs struct {
	dataSlice []*Job
	dataMap   map[string]*Job
	once      sync.Once
}

func (es *Jobs) GetSlice() []*Job {
	return es.dataSlice
}

func (es *Jobs) GetMap() map[string]*Job {
	es.once.Do(func() {
		es.dataMap = make(map[string]*Job)
		for _, v := range es.dataSlice {
			es.dataMap[v.title] = v
		}
	})
	return es.dataMap
}

func (es *Jobs) GetByTitle(title string) (*Job, error) {
	d := es.GetMap()
	e, ok := d[title]
	if ok {
		return e, nil
	} else {
		return nil, ErrRecordNotFound
	}
}

type Job struct {
	ctx     context.Context
	title   string
	handler func(in interface{}) (interface{}, error)
	input   interface{}
	output  interface{}
	error   error
}

func newJob(ctx context.Context, title string, input interface{}, handler func(in interface{}) (interface{}, error)) *Job {
	e := new(Job)
	e.ctx = ctx
	e.title = title
	e.handler = handler
	e.input = input
	e.error = ErrUndone
	return e
}
func (e *Job) Title() string {
	return e.title
}

func (e *Job) Result() (interface{}, error) {
	return e.output, e.error
}
