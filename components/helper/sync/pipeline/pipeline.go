package pipeline

import (
	"context"
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
	var res []Event
	fin := make(chan struct{})
	go func() {
		res = thisPipeline.do(events...)
		fin <- struct{}{}
	}()

	select {
	case <-thisPipeline.ctx.Done():
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

func (thisPipeline *Pipeline) consumer(inEvent Event) Event {
	select {
	case <-thisPipeline.done:
		inEvent.Error = thisPipeline.ctx.Err()
		return inEvent
	default:
		if inEvent.Ctx == nil {
			inEvent.Error = inEvent.Handler(&inEvent)
			return inEvent
		}

		fin := make(chan error, 0)
		go func() {
			defer func() {
				if err := recover(); err != nil {
					return
				}
			}()
			inEvent.Error = inEvent.Handler(&inEvent)
			fin <- inEvent.Error
		}()

		select {
		case <-inEvent.Ctx.Done():
			inEvent.Error = inEvent.Ctx.Err()
			return inEvent
		case <-fin:
			return inEvent
		}

	}
}
