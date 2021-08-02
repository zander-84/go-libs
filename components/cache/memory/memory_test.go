package memory

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

// go test -v  -run TestMemory
func TestMemory(t *testing.T) {
	memory := NewMemory(Conf{})
	if err := memory.Start(); err != nil {
		t.Fatal("start memory err: ", err.Error())
	}
	defer memory.Stop()
	testMemoryString(t, memory)
	testMemorySlice(t, memory)
	testRdbStruct(t, memory)
	testRdbExists(t, memory)
	testRdbGetOrSet(t, memory)
	testRdbFlushDB(t, memory)

}

func testRdbFlushDB(t *testing.T, memory *Memory) {
	if err := memory.FlushDB(context.Background()); err != nil {
		t.Fatal("FlushDB  err: ", err.Error())
	}
}

func testRdbGetOrSet(t *testing.T, memory *Memory) {
	key := "str"
	val := "aaaa"
	get := ""

	if err := memory.GetOrSet(context.Background(), key, &get, 100*time.Second, func() (interface{}, error) {
		return val, nil
	}); err != nil {
		t.Fatal("testRdbGetOrSet  err: ", err.Error())
	}

	if val != get {
		t.Fatal("testRdbGetOrSet err: ", get)
	}

}
func testRdbExists(t *testing.T, memory *Memory) {
	key := "maps"
	if ok, err := memory.Exists(context.Background(), key, key); err != nil {
		t.Fatal("Exists  err: ", err.Error())
	} else if !ok {

	}
	testMemoryString(t, memory)
	if ok, err := memory.Exists(context.Background(), "string", "string"); err != nil {
		t.Fatal("Exists  err: ", err.Error())
	} else if !ok {
		t.Fatal("Exists   err ok: ")
	} else if ok {
	}

}
func testRdbStruct(t *testing.T, memory *Memory) {
	key := "struct"
	type class struct {
		term int
	}
	type Student struct {
		Name  string
		Age   int
		Class *class
	}
	val := &Student{
		Name: "zander",
		Age:  10,
		Class: &class{
			term: 1,
		},
	}
	if err := memory.Set(context.Background(), key, val, 100*time.Second); err != nil {
		t.Fatal("set struct err: ", err.Error())
	}
	var get = new(Student)
	if err := memory.Get(context.Background(), key, get); err != nil {
		t.Fatal("get struct err: ", err.Error())
	}
	if get.Name != val.Name || get.Age != val.Age {
		t.Fatal("set get struct err: ", get)
	}
	//fmt.Println(get)
	//get.Age=18
	//get.Class.term=99
	var get2 = new(Student)
	if err := memory.Get(context.Background(), key, get2); err != nil {
		t.Fatal("get struct err: ", err.Error())
	}
	//fmt.Println(get2.Class.term)

}

func testMemoryString(t *testing.T, memory *Memory) {
	key := "string"
	var valPtr *string
	val := "string2"
	valPtr = &val
	if err := memory.Set(context.Background(), key, valPtr, 0); err != nil {
		t.Fatal("set string err: ", err.Error())
	}
	get := ""
	if err := memory.Get(context.Background(), key, &get); err != nil {
		t.Fatal("get string err: ", err.Error())
	}
}

func testMemorySlice(t *testing.T, memory *Memory) {
	key := "slice"
	val := []string{"a", "b"}
	if err := memory.Set(context.Background(), key, val, 0); err != nil {
		t.Fatal("set string err: ", err.Error())
	}
	var get []string
	if err := memory.Get(context.Background(), key, &get); err != nil {
		t.Fatal("get string err: ", err.Error())
	}

	val[0] = "c"
	var get2 []string
	if err := memory.Get(context.Background(), key, &get2); err != nil {
		t.Fatal("get string err: ", err.Error())
	}
}

// go test -v  -bench=BenchmarkReflect
//BenchmarkReflect-8      13439902                87.3 ns/op
func BenchmarkReflect(b *testing.B) {
	type Student struct {
		Name string
		Age  int
	}
	var s1 = &Student{
		"zander",
		18,
	}

	for i := 0; i < b.N; i++ {
		var s2 = new(Student)
		if reflect.ValueOf(s1).Type().Kind() == reflect.Ptr {
			reflect.ValueOf(s2).Elem().Set(reflect.ValueOf(s1).Elem())
		} else {
			reflect.ValueOf(s2).Elem().Set(reflect.ValueOf(s1))
		}
	}
}

//go test -v  -bench=BenchmarkAssert
//BenchmarkAssert-8       1000000000               0.339 ns/op
func BenchmarkAssert(b *testing.B) {
	type Student struct {
		Name string
		Age  int
	}
	var s1 = Student{
		"zandera",
		18,
	}

	for i := 0; i < b.N; i++ {
		var s2 interface{}
		s2 = s1
		if _, ok := s2.(Student); ok {

		}
	}
}

// go test -v  -bench=BenchmarkUnmarshal
//BenchmarkUnmarshal-8     1000000              1015 ns/op
func BenchmarkUnmarshal(b *testing.B) {
	type Student struct {
		Name string
		Age  int
	}
	var s1 = Student{
		"zander",
		18,
	}
	vs1, _ := json.Marshal(s1)
	for i := 0; i < b.N; i++ {
		var s2 Student
		json.Unmarshal(vs1, &s2)
	}
}
