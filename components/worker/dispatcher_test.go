package worker

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

// go test -v  -run TestDispatcher
func TestDispatcher(t *testing.T) {
	d := NewDispatcher(Conf{
		MaxWorkers: 100,
		MaxQueues:  1000000,
	})

	if err := d.Start(); err != nil {
		t.Fatal("start Dispatcher err: ", err.Error())
	}
	var a int64
	cnt := 1000000
	for i := 0; i < cnt; i++ {
		if err := d.AddJob(JobFunc(func() error {
			atomic.AddInt64(&a, 1)
			return nil
		})); err != nil {
			fmt.Println(err.Error())
		}
	}
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1000*time.Millisecond)
	defer cancel()
	d.Stop(ctx)

	if a == int64(cnt) {
		t.Log("success")
	} else {
		fmt.Println("a: ", a)
		fmt.Println("cnt: ", cnt)
		t.Fatal("error")
	}

}
