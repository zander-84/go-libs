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
	c1 := Event{
		Ctx: context.Background(),
		Handler: func(e *Event) error {
			data := e.Input.(int)
			e.Output = data * 9
			return nil
		},
		Input: 11,
		Error: nil,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	c2 := Event{
		Ctx: ctx,
		Handler: func(e *Event) error {
			time.Sleep(2 * time.Second)
			data := e.Input.(int)
			e.Output = data * 9
			return nil
		},
		Input: 22,
		Error: nil,
	}

	t.Log(a.Do(done, c1, c2))
	time.Sleep(10 * time.Second)
}
