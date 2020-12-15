package pipeline

import (
	"context"
	"fmt"
	"sync"
)

type Pipeline struct {
	ctx     context.Context
	workers chan struct{}
	done    chan struct{}
	jobLen  int
}

func NewPipeline(ctx context.Context, workers int) *Pipeline {
	thisPipeline := new(Pipeline)
	thisPipeline.ctx = ctx
	thisPipeline.done = make(chan struct{})
	if workers >= 0 {
		thisPipeline.workers = make(chan struct{}, workers)
	}
	return thisPipeline
}
func (thisPipeline *Pipeline) Do(events ...Event) ([]Event, error) {
	if thisPipeline.ctx.Done() == nil {
		return thisPipeline.do(events...), nil
	}

	var res []Event
	fin := make(chan struct{})
	go func() {
		res = thisPipeline.do(events...)
		fin <- struct{}{}
	}()

	select {
	case <-thisPipeline.ctx.Done():
		close(thisPipeline.done)
		return nil, thisPipeline.ctx.Err()
	case <-fin:
		return res, nil
	}
}

func (thisPipeline *Pipeline) do(events ...Event) []Event {
	thisPipeline.jobLen = len(events)
	var res = make([]Event, 0)
	for e := range thisPipeline.mergeEvent(thisPipeline.gen(events...)) {
		res = append(res, e)
	}
	return res
}
func (thisPipeline *Pipeline) mergeEvent(events <-chan Event) <-chan Event {
	var wg sync.WaitGroup
	outEvent := make(chan Event, thisPipeline.jobLen)

	for e := range events {
		if thisPipeline.workers != nil {
			thisPipeline.workers <- struct{}{}
		}
		wg.Add(1)
		go func(evt Event) {
			defer func() {
				wg.Done()
				if thisPipeline.workers != nil {
					<-thisPipeline.workers
				}
			}()
			outEvent <- thisPipeline.consumer(evt)
		}(e)
	}

	go func() {
		wg.Wait()
		close(outEvent)
	}()

	return outEvent
}

func (thisPipeline *Pipeline) gen(events ...Event) <-chan Event {
	eventsChan := make(chan Event, len(events))
	go func() {
		defer close(eventsChan)
		for _, e := range events {
			eventsChan <- e
		}
	}()
	return eventsChan
}

func (thisPipeline *Pipeline) consumer(event Event) Event {
	select {
	case <-thisPipeline.done:
		event.Error = thisPipeline.ctx.Err()
		return event
	default:
		if event.Ctx == nil {
			func() {
				defer func() {
					if err := recover(); err != nil {
						event.Error = fmt.Errorf("painc %v", err)
					}
				}()
				event.Error = event.Handler(&event)
			}()
			return event
		}

		fin := make(chan error, 0)
		go func() {
			defer func() {
				if err := recover(); err != nil {
					event.Error = fmt.Errorf("painc %v", err)
					fin <- event.Error
					return
				}
			}()
			event.Error = event.Handler(&event)
			fin <- event.Error
		}()

		select {
		case <-event.Ctx.Done():
			event.Error = event.Ctx.Err()
			return event
		case <-fin:
			return event
		}

	}
}
