package worker

import (
	"sync/atomic"
	"testing"
	"time"
)

// go test -v  -run TestDispatcher
func TestDispatcher(t *testing.T) {
	d := NewDispatcher(Conf{
		MaxWorkers: 10,
		MaxQueues:  1000,
	})

	if err := d.Start(); err != nil {
		t.Fatal("start Dispatcher err: ", err.Error())
	}
	var a int64
	cnt := 100000
	for i := 0; i < cnt; i++ {
		d.AddJob(JobFunc(func() error {
			atomic.AddInt64(&a, 1)
			return nil
		}))
	}

	time.Sleep(2 * time.Second)
	if a == int64(cnt) {
		t.Log("success")
	} else {
		t.Fatal("error")
	}

}
