package sd

import (
	"testing"
	"time"
)

func TestRoundRobin(t *testing.T) {
	l := NewListener("test", "")
	l.Set(map[string]int{
		"127.0.0.1:8081": 1,
		"127.0.0.1:8082": 1,
		"127.0.0.1:8083": 1,
	})
	rr := NewRoundRobin(l, true)

	for i := 0; i < 100000; i++ {
		go func(i int) {
			rr.Next()
		}(i)
	}
	go func() {
		l.Remove("127.0.0.1:8082")
	}()
	time.Sleep(10 * time.Second)
	t.Log(rr.Used())
}

//BenchmarkRoundRobin-8   	11099640	       114 ns/op  保存数据
//BenchmarkRoundRobin-8   	20444817	        57.6 ns/op  不存数据
func BenchmarkRoundRobin(b *testing.B) {
	l := NewListener("test", "")
	l.Set(map[string]int{
		"127.0.0.1:8081": 1,
		"127.0.0.1:8082": 1,
		"127.0.0.1:8083": 1,
	})

	rr := NewRoundRobin(l, false)
	for i := 0; i < b.N; i++ {
		rr.Next()
	}
}
