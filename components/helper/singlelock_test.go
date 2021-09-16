package helper

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewLock(t *testing.T) {
	lock := NewMemoryLock()
	var data int64
	w := sync.WaitGroup{}

	go func() {
		for {
			time.Sleep(time.Second)
			fmt.Println("running at: ", atomic.LoadInt64(&data))
		}
	}()

	for i := 0; i < 20000000; i++ {
		w.Add(1)
		go func(i int) {
			defer w.Done()

			lock.Lock(context.Background(), fmt.Sprintf("%d", i), fmt.Sprintf("%d", i), time.Second*2)
			lock.UnLock(context.Background(), fmt.Sprintf("%d", i), fmt.Sprintf("%d", i))

			if err := lock.Lock(context.Background(), fmt.Sprintf("%d", 1), fmt.Sprintf("%d", 1), time.Second*2); err == nil {
				fmt.Println("lock at:", i, " ", defaultTime.FormatHyphenFromNow())
			}
			atomic.AddInt64(&data, 1)
		}(i)
	}

	//for i := 0; i < 100; i++ {
	//	w.Add(1)
	//	go func(i int) {
	//		defer w.Done()
	//		if err := lock.Lock(context.Background(), fmt.Sprintf("1-%d", i), fmt.Sprintf("val-%d", i), time.Second*time.Duration(i)); err == nil {
	//		}
	//		atomic.AddInt64(&data, 1)
	//	}(i)
	//}

	fmt.Println("Wait: ", defaultTime.FormatHyphenFromNow())
	w.Wait()
	time.Sleep(2 * time.Second)

	//lock.Println()
	fmt.Println("finish: ", defaultTime.FormatHyphenFromNow())

}

//BenchmarkNewLock-16    	 1000000	      1609 ns/op
func BenchmarkNewLock(b *testing.B) {
	lock := NewMemoryLock()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			lock.Lock(context.Background(), fmt.Sprintf("%d", 1), fmt.Sprintf("%d", 1), time.Second*2)
		}
	})
}
