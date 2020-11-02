package helper

import (
	"testing"
	"time"
)

// go test -v  -run TestGoGroup_Do
func TestGoGroup_Do(t *testing.T) {
	a := NewGoGroup()
	name := "zander"
	age := 18
	a.AsyncFunc(func() {
		time.Sleep(2 * time.Second)
		t.Log(name)
	}, func() {
		t.Log(age)
	})
}
