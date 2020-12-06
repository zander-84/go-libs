package singleton

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// go test -v  -run TestNewSingleton
func TestNewSingleton(t *testing.T) {
	w := NewSingleton(5)
	go w.Run()

	_ = w.AddJob(&Job{
		Ctx:  context.Background(),
		Name: "job1",
		Handler: func(j *Job) error {
			fmt.Println(j.Input)
			j.Output = j.Input
			return nil
		},
		Input:  "job1",
		Output: nil,
		Error:  nil,
	})

	_ = w.AddJob(&Job{
		Ctx:  context.Background(),
		Name: "job2",
		Handler: func(j *Job) error {
			fmt.Println(j.Input)
			j.Output = j.Input
			return nil
		},
		Input:  "job2",
		Output: nil,
		Error:  nil,
	})
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(2*time.Second))
	defer cancel()

	w.Shutdown(ctx)
	for _, v := range w.Data {
		fmt.Printf("%v", v)
	}

}
