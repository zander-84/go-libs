package helper

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestUnique(t *testing.T) {
	u := NewUnique("LF", "A", "-")
	var c = make(map[string]int, 0)
	var l sync.Mutex
	t1 := time.Now()
	u.ID()
	w := sync.WaitGroup{}
	for i := 0; i < 2000010; i++ {
		w.Add(1)
		go func(i int) {
			defer w.Done()
			//_ = u.ID()
			u.ID()
			u.ID()
			u.ID()
			u.ID()
			u.ID()
			u.ID()
			dd := u.ID()
			l.Lock()
			if vv, ok := c[dd]; ok {

				t.Error(i, dd+"存在", vv)
			} else {
				c[dd] = i
			}
			l.Unlock()
		}(i)
	}
	t2 := time.Now()
	fmt.Println(t2.Unix() - t1.Unix())
	fmt.Println("fgggg")
	w.Wait()
	go func() {
		fmt.Println(u.ID())
	}()
	l.Lock()
	fmt.Println(len(c))
	l.Unlock()
	time.Sleep(2 * time.Second)
}

//BenchmarkUnique-16    	  904142	      1280 ns/op
func BenchmarkUnique(b *testing.B) {
	u := NewUnique("LF", "A", "-")
	for i := 0; i < b.N; i++ {
		u.IDWithTag(u.CharLower(5))
	}
}
