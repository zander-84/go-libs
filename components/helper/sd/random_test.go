package sd

import (
	"testing"
	"time"
)

func TestRandom(t *testing.T) {
	l := NewListener("test", "")
	l.Set(map[string]int{
		"127.0.0.1:8081": 1,
		"127.0.0.1:8082": 1,
		"127.0.0.1:8083": 1,
	})
	rr := NewRandom(l, true)

	for i := 0; i < 1000000; i++ {
		go func(i int) {
			if d, err := rr.Next(); err != nil {
				t.Error(err.Error())
			} else if d == "" {
				t.Error("空数据")
			}
		}(i)
	}
	go func() {
		l.Remove("127.0.0.1:8082")
	}()
	time.Sleep(3 * time.Second)
	t.Log(rr.Used())
}

//BenchmarkRandom-8   	 9791858	       123 ns/op  保存数据
//BenchmarkRandom-8   	18463314	        63.7 ns/op  不存数据
func BenchmarkRandom(b *testing.B) {
	l := NewListener("test", "")
	l.Set(map[string]int{
		"127.0.0.1:8081": 1,
		"127.0.0.1:8082": 1,
		"127.0.0.1:8083": 1,
	})

	rr := NewRandom(l, true)
	for i := 0; i < b.N; i++ {
		rr.Next()
	}
}
