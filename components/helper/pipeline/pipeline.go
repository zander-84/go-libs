package pipeline

import (
	"context"
	"sync"
)

type Pipeline struct{}

func NewPipeline() *Pipeline {
	return new(Pipeline)
}

func (this *Pipeline) Do(done chan struct{}, events ...Event) []Event {
	eventChans := make([]<-chan Event, 0)
	for _, e := range events {
		eventChans = append(eventChans, consumer(done, producer(e)))
	}
	var res = make([]Event, 0)
	for e := range mergeEvent(done, eventChans...) {
		res = append(res, e)
	}
	return res
}

type Event struct {
	Ctx     context.Context
	Name    string
	Handler func(input interface{}) (interface{}, error)
	Input   interface{}
	Output  interface{}
	Error   error
}

func producer(events ...Event) <-chan Event {
	eventsChan := make(chan Event, len(events))
	go func() {
		defer close(eventsChan)
		for _, e := range events {
			eventsChan <- e
		}
	}()
	return eventsChan
}

func consumer(done <-chan struct{}, inEvent <-chan Event) <-chan Event {
	outEvent := make(chan Event)
	go func() {
		defer close(outEvent)
		for e := range inEvent {
			select {
			case <-done:
				return
			default:
				func() {
					fin := make(chan bool, 0)
					go func() {
						defer func() {
							if err := recover(); err != nil {
								return
							}
						}()
						e.Output, e.Error = e.Handler(e.Input)
						fin <- true
					}()

					select {
					case <-e.Ctx.Done():
						e.Error = e.Ctx.Err()
						outEvent <- e
					case <-fin:
						outEvent <- e
					}

					//e.Output, e.Error = e.Handler(e.Input)
					//outEvent <- e
				}()
			}
		}
	}()
	return outEvent
}

func mergeEvent(done <-chan struct{}, events ...<-chan Event) <-chan Event {
	var wg sync.WaitGroup
	outEvent := make(chan Event)

	output := func(es <-chan Event) {
		for e := range es {
			select {
			case outEvent <- e:
			case <-done:
			}
		}
		wg.Done()
	}
	wg.Add(len(events))

	for _, e := range events {
		go output(e)
	}

	go func() {
		wg.Wait()
		close(outEvent)
	}()

	return outEvent
}
