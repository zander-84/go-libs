package memory

import (
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
	testRdbGetOrSetFast(t, memory)
	testRdbFlushDB(t, memory)

}

func testRdbFlushDB(t *testing.T, memory *Memory) {
	if err := memory.FlushDB(); err != nil {
		t.Fatal("FlushDB  err: ", err.Error())
	}
}

func testRdbGetOrSetFast(t *testing.T, memory *Memory) {
	key := "GetOrSetFast"
	val := "aaaa"
	get := ""

	get2, err := memory.GetOrSetFast(key, &get, 100*time.Second, func() (interface{}, error) {
		return val, nil
	})
	if err != nil {
		t.Fatal("Exists  err: ", err.Error())
	}

	if get3, ok := get2.(string); ok {
		if val != get3 {
			t.Fatal("testRdbGetOrSetFast err: ", get2)
		}
	} else {
		t.Fatal("testRdbGetOrSetFast err: ", get2)
	}

	get2, err = memory.GetOrSetFast(key, &get, 100*time.Second, func() (interface{}, error) {
		return "err", nil
	})
	if err != nil {
		t.Fatal("Exists  err: ", err.Error())
	}

	if get3, ok := get2.(string); ok {
		if val != get3 {
			t.Fatal("testRdbGetOrSetFast err: ", get2)
		}
	} else {
		t.Fatal("testRdbGetOrSetFast err: ", get2)
	}
}

func testRdbGetOrSet(t *testing.T, memory *Memory) {
	key := "str"
	val := "aaaa"
	get := ""

	if err := memory.GetOrSet(key, &get, 100*time.Second, func() (interface{}, error) {
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
	if ok, err := memory.Exists(key, key); err != nil {
		t.Fatal("Exists  err: ", err.Error())
	} else if !ok {

	}
	testMemoryString(t, memory)
	if ok, err := memory.Exists("string", "string"); err != nil {
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
	if err := memory.Set(key, val, 100*time.Second); err != nil {
		t.Fatal("set struct err: ", err.Error())
	}
	var get = new(Student)
	if err := memory.Get(key, get); err != nil {
		t.Fatal("get struct err: ", err.Error())
	}
	if get.Name != val.Name || get.Age != val.Age {
		t.Fatal("set get struct err: ", get)
	}
	//fmt.Println(get)
	//get.Age=18
	//get.Class.term=99
	var get2 = new(Student)
	if err := memory.Get(key, get2); err != nil {
		t.Fatal("get struct err: ", err.Error())
	}
	//fmt.Println(get2.Class.term)

}

func testMemoryString(t *testing.T, memory *Memory) {
	key := "string"
	var valPtr *string
	val := "string2"
	valPtr = &val
	if err := memory.Set(key, valPtr, 0); err != nil {
		t.Fatal("set string err: ", err.Error())
	}
	get := ""
	if err := memory.Get(key, &get); err != nil {
		t.Fatal("get string err: ", err.Error())
	}
}

func testMemorySlice(t *testing.T, memory *Memory) {
	key := "slice"
	val := []string{"a", "b"}
	if err := memory.Set(key, val, 0); err != nil {
		t.Fatal("set string err: ", err.Error())
	}
	var get []string
	if err := memory.Get(key, &get); err != nil {
		t.Fatal("get string err: ", err.Error())
	}

	val[0] = "c"
	var get2 []string
	if err := memory.Get(key, &get2); err != nil {
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
		"zander",
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
		"zander ",
		18,
	}
	vs1, _ := json.Marshal(s1)
	for i := 0; i < b.N; i++ {
		var s2 Student
		json.Unmarshal(vs1, &s2)
	}
}
