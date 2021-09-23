package etcd

import (
	"fmt"
	"github.com/zander-84/go-libs/components/helper/sd"
	"sync"
	"testing"
	"time"
)

func TestNewEtcd(t *testing.T) {
	es := NewEtcd(Conf{
		Endpoints: []string{"127.0.0.1:2379"},
		TlsPem:    ``,
		TlsKey:    ``,
		TlsCa:     ``,
	})

	if err := es.Start(); err != nil {
		t.Fatal(err.Error())
	}
	if err := es.Start(); err != nil {
		t.Fatal(err.Error())
	}
	rs := RegisterServer("/project/endpoint/test/127.0.0.1:80", "/project/endpoint/test/127.0.0.1:80", 1)
	if err := es.RegisterServer(rs); err != nil {
		t.Fatal(err.Error())
	}

	go func() {
		for {

			rs3 := RegisterServer("/project/endpoint/test/127.0.0.1:82", "/project/endpoint/test/127.0.0.1:82", 1)
			if err := es.RegisterServer(rs3); err != nil {
				t.Fatal(err.Error())
			}

			rs2 := RegisterServer("/project/endpoint/test/127.0.0.1:8080-3", "/project/endpoint/test/127.0.0.1:8080", 1)
			if err := es.RegisterServer(rs2); err != nil {
				t.Fatal(err.Error())
			}
			time.Sleep(2 * time.Second)
			es.Deregister(rs2)
			es.Deregister(rs3)
		}

	}()
	lis := sd.NewListener("/project/endpoint/test/", "-")
	es.Watch(lis)

	time.Sleep(2 * time.Second)

	for {
		time.Sleep(time.Second / 2)
		go func() {
			fmt.Println("1")
			lis.Println()
		}()
		go func() {
			fmt.Println("2")
			lis.Println()
		}()
	}
	rr := sd.NewRoundRobin(lis, true)
	w := sync.WaitGroup{}
	for i := 0; i < 10000000; i++ {
		w.Add(1)
		go func(i int) {
			defer w.Done()

			if d, err := rr.Next(); err != nil {
				t.Error(err.Error())
			} else if d == "" {
				t.Error("空数据")
			}
		}(i)
	}

	w.Wait()
	t.Log(rr.Used())
}
