package etcd

import (
	"fmt"
	"github.com/zander-84/go-libs/components/helper/sd"
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
	rs := RegisterServer("/project/endpoint/test/127.0.0.1:80", "/project/endpoint/test/127.0.0.1:80", 10)
	if err := es.RegisterServer(rs); err != nil {
		t.Fatal(err.Error())
	}
	rs2 := RegisterServer("/project/endpoint/test/127.0.0.1:8080-2", "/project/endpoint/test/127.0.0.1:8080", 10)
	if err := es.RegisterServer(rs2); err != nil {
		t.Fatal(err.Error())
	}

	lis := sd.NewListener("/project/endpoint/test/")
	es.Watch(lis)
	t.Log("success")
	go func() {
		time.Sleep(10 * time.Second)
		fmt.Println("do Deregister")
		es.Deregister(rs2)
	}()
	for {
		lis.Println()
		time.Sleep(time.Second)
	}
}
