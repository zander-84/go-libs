package robfig

import (
	"fmt"
	"github.com/zander-84/go-libs/components/cron"
	"testing"
	"time"
)

// go test -v  -run TestRobfig
func TestRobfig(t *testing.T) {
	c := NewRobfig(Conf{})
	if err := c.Start(); err != nil {
		t.Fatal("start cron err: ", err.Error())
	}

	testAdd(t, c)
	testRemove(t, c)
	testAdd(t, c)

	if err := c.StartJobs(); err != nil {
		t.Fatal("StartJobs error ", err.Error())
	}
	if err := c.StartJobs(); err != nil {
		t.Fatal("StartJobs error ", err.Error())
	}

	time.Sleep(2 * time.Minute)

	t.Log("success")
}

func testAdd(t *testing.T, c *Robfig) {
	if err := c.AddJob(&cron.Job{
		ID:   "test1",
		Desc: "测试1",
		Spec: "* * * * * *",
		Cmd: cron.CmdFunc(func() {
			fmt.Println("hello world")
		}),
		Obj: nil,
	}); err != nil {
		t.Fatal("add error ", err.Error())
	}

	if err := c.AddJob(&cron.Job{
		ID:   "test2",
		Desc: "测试2",
		Spec: "* * * * * *",
		Cmd: cron.CmdFunc(func() {
			fmt.Println("hello world 2222")
		}),
		Obj: nil,
	}); err != nil {
		t.Fatal("add error ", err.Error())
	}
}
func testRemove(t *testing.T, c *Robfig) {
	if err := c.RemoveJob("test1"); err != nil {
		t.Fatal("add error ", err.Error())
	}
	if err := c.RemoveJob("test2"); err != nil {
		t.Fatal("add error ", err.Error())
	}
}
