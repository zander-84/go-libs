package sd

import (
	"fmt"
	"testing"
	"time"
)

func TestConsistentHash(t *testing.T) {
	l := NewListener("test", "")
	l.Set(map[string]int{
		"127.0.0.1:8081": 1,
		"127.0.0.1:8082": 2,
		"127.0.0.1:8083": 3,
	})
	rr := NewConsistentHash(l, true)

	for i := 0; i < 100; i++ {
		go func(i int) {
			if d, err := rr.Get(fmt.Sprintf("%d", i)); err != nil {
				t.Log(err.Error())
			} else {
				t.Log(d)
			}
			//if d,err:=rr.Get("123");err!=nil{
			// t.Log(err.Error())
			//}else {
			// t.Log(d)
			//}
		}(i)
	}

	time.Sleep(3 * time.Second)
	t.Log(rr.Used())
}

//BenchmarkConsistentHash-8   2638166	456 ns/op  保存数据
//BenchmarkRoundRobin-8   	  2944742	        393 ns/op  不存数据
func BenchmarkConsistentHash(b *testing.B) {
	l := NewListener("test", "")
	l.Set(map[string]int{
		"127.0.0.1:8081": 1,
		"127.0.0.1:8082": 1,
		"127.0.0.1:8083": 2,
	})

	rr := NewConsistentHash(l, false)
	for i := 0; i < b.N; i++ {
		if _, err := rr.Get(fmt.Sprintf("%d", i)); err != nil {
			b.Log(err.Error())
		}
	}
}
