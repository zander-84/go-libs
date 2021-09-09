package helper

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestUnique(t *testing.T) {
	u := NewUnique("LF", "A", "-")
	var c = make(map[string]bool, 0)
	var l sync.Mutex
	t1 := time.Now()
	for i := 0; i < 2000010; i++ {
		go func(i int) {
			_ = u.ID()
			//dd:=u.ID()
			//l.Lock()
			//if _,ok:=c[dd];ok{
			//	t.Error(dd+"存在")
			//}else {
			//	c[dd] = true
			//}
			//l.Unlock()
		}(i)
	}
	t2 := time.Now()
	fmt.Println(t2.Unix() - t1.Unix())
	fmt.Println("fgggg")

	time.Sleep(4 * time.Second)
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
