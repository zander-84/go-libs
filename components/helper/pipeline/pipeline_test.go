package pipeline

import (
	"context"
	"testing"
	"time"
)

// go test -v  -run TestGoGroup_Do
func TestPipeline_DOPipeline(t *testing.T) {
	a := NewPipeline()
	done := make(chan struct{}, 1)
	defer close(done)
	//c1 := Event{
	//	Ctx: context.Background(),
	//	Handler: func(e *Event) error {
	//		data := e.Input.(int)
	//		e.Output = data * 9
	//		return nil
	//	},
	//	Input: 11,
	//	Error: nil,
	//}
	//ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	//defer cancel()
	//c2 := Event{
	//	Ctx: ctx,
	//	Handler: func(e *Event) error {
	//		time.Sleep(2 * time.Second)
	//		data := e.Input.(int)
	//		e.Output = data * 9
	//		return nil
	//	},
	//	Input: 22,
	//	Error: nil,
	//}
	events := make([]Event, 0)
	for i := 0; i <= 100; i++ {
		events = append(events, Event{
			Ctx: context.Background(),
			Handler: func(e *Event) error {
				time.Sleep(3 * time.Second)
				data := e.Input.(int)
				e.Output = data * 9
				return nil
			},
			Input: 22,
			Error: nil,
		})
	}
	t.Log(a.Do(done, events...))
	//time.Sleep(10 * time.Second)
}
