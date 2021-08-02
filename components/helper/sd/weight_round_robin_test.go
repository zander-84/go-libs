package sd

import (
	"testing"
	"time"
)

func TestWeightRoundRobin(t *testing.T) {
	l := NewListener("test")
	l.Set(map[string]int{
		"127.0.0.1:8081": 1,
		"127.0.0.1:8082": 2,
		"127.0.0.1:8083": 3,
	})
	rr := NewWeightRoundRobin(l, true)

	for i := 0; i < 10000000; i++ {
		go func(i int) {
			if _, err := rr.Next(); err != nil {
				t.Log(err.Error())
			}
		}(i)
	}
	//go func() {
	//	l.Remove("127.0.0.1:8082")
	//}()
	time.Sleep(5 * time.Second)
	t.Log(rr.Used())
}

//BenchmarkRoundRobin-8   	11099640	       115 ns/op  保存数据
//BenchmarkRoundRobin-8   	20444817	        58.2 ns/op  不存数据
func BenchmarkWeightRoundRobin(b *testing.B) {
	l := NewListener("test")
	l.Set(map[string]int{
		"127.0.0.1:8081": 1,
		"127.0.0.1:8082": 1,
		"127.0.0.1:8083": 2,
	})

	rr := NewRoundRobin(l, true)
	for i := 0; i < b.N; i++ {
		if _, err := rr.Next(); err != nil {
			b.Log(err.Error())
		}
	}
}
