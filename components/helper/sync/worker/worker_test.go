package worker

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// go test -v  -run TestGoGroup_Do
func TestWorker_DO(t *testing.T) {
	fmt.Println("start:", time.Now().Format("2006-01-02 15:04:05"))
	p := NewWorker()
	for i := 0; i < 20000000; i++ {
		p.AddJob(context.Background(), fmt.Sprintf("%d", i), i, func(in interface{}) (interface{}, error) {
			//time.Sleep(2*time.Second)
			return in, nil
		})
	}

	c1, cacel2 := context.WithTimeout(context.Background(), 1*time.Second)
	defer cacel2()

	p.AddJob(c1, "student", "student", func(in interface{}) (interface{}, error) {
		fmt.Println("before Sleep 3")
		time.Sleep(5 * time.Second)
		fmt.Println("after Sleep 3")

		return in, nil
	})
	c1, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	jobs, err := p.Do(c1, 20)
	fmt.Println("fin Do")
	jobs.GetMap()
	fmt.Println(len(jobs.dataMap))
	//fmt.Println(jobs.dataMap)
	//fmt.Println(student)
	//time.Sleep(3 * time.Second)
	if err != nil {
		fmt.Println("do err: ", err.Error())
	} else {
		fmt.Println("no err")
		//for k, e := range a.GetSlice() {
		//	fmt.Println(k, "==>", e.title, e.error)
		//}
		//for k, e := range a.GetMap() {
		//	fmt.Println(k,"==>",fmt.Sprintf("%+v", e))
		//}
	}
	//time.Sleep(3 * time.Second)
	fmt.Println("end:", time.Now().Format("2006-01-02 15:04:05"))
}

//BenchmarkDOWorker-8   	  236150	      4573 ns/op
func BenchmarkDOWorker(b *testing.B) {
	for i := 0; i < b.N; i++ {
		func() {
			p := NewWorker()
			for ii := 0; ii < 3; ii++ {
				p.AddJob(context.Background(), fmt.Sprintf("%d", ii), ii, func(in interface{}) (interface{}, error) {
					return in, nil
				})
			}
			_, err := p.Do(context.Background(), 10)
			if err != nil {
				fmt.Println(err.Error())
			}
		}()
	}
}
