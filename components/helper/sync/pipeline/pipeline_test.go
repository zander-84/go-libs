package pipeline

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// go test -v  -run TestGoGroup_Do
func TestPipeline_DOPipeline(t *testing.T) {
	events := make([]Event, 0)
	for i := 0; i < 4; i++ {
		events = append(events, Event{
			//Ctx:  context.Background(),
			Name: fmt.Sprintf("%d", i),
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
	c1, cancel := context.WithTimeout(context.Background(), 200*time.Second)
	defer cancel()

	fmt.Println("start:", time.Now().Format("2006-01-02 15:04:05"))
	p := NewPipeline(c1, 2)
	a, err := p.Do(events...)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		for _, e := range a {
			fmt.Println(fmt.Sprintf("%+v", e))
		}
	}
	fmt.Println("end:", time.Now().Format("2006-01-02 15:04:05"))
}
