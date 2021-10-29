package goredis

import (
	"context"
	"flag"
	"fmt"
	"strconv"
	"testing"
	"time"
)

// go test -v  -run TestRdb  -args 172.16.86.160:6379  123456 1
func TestRdb(t *testing.T) {
	if !flag.Parsed() {
		flag.Parse()
	}
	argList := flag.Args()
	db, _ := strconv.Atoi(argList[2])
	rdb := NewRdb(Conf{
		Addr:     argList[0],
		Password: argList[1],
		Db:       db,
		PoolSize: 10,
		MinIdle:  5,
	})
	if err := rdb.Start(); err != nil {
		t.Fatal("start redis err: ", err.Error())
	}
	testRdbMGet(t, rdb)
	//go func() { testRdbString(t, rdb) }()
	//go func() { testRdbString(t, rdb) }()
	//go func() { testRdbString2(t, rdb) }()
	//testRdbInt(t, rdb)
	//testRdbStruct(t, rdb)
	//testRdbMap(t, rdb)
	//testRdbDel(t, rdb)
	//testRdbExists(t, rdb)
	//go func() {testRdbGetOrSet(t, rdb)}()
	//go func() {testRdbGetOrSet2(t, rdb)}()
	//testRdbGetOrSetFast(t, rdb)
	//testRdbFlushDB(t, rdb)
	//time.Sleep(2 * time.Second)

	t.Log("success")
}
func testRdbMGet(t *testing.T, rdb *Rdb) {
	key1 := "string1"
	val1 := 2.1
	if err := rdb.Set(context.Background(), key1, val1, 100*time.Second); err != nil {
		t.Fatal("set string err: ", err.Error())
	}
	key2 := "string2"
	val2 := 3.1
	if err := rdb.Set(context.Background(), key2, val2, 100*time.Second); err != nil {
		t.Fatal("set string err: ", err.Error())
	}
	res := make([]float64, 0)
	if err := rdb.MustMGetOrSet(context.Background(), []string{key1, "fsaara", key2, "fffgggs"}, &res, 10*time.Second, func() (interface{}, error) {
		return 1.1, nil
	}); err != nil {
		t.Fatal("MGet  err: ", err.Error())
	} else {
		fmt.Println(res)
	}
	//if _,err := rdb.MGet(context.Background(), []string{key1, "fs", key2, "ff"}, &res); err != nil {
	//	t.Fatal("MGet  err: ", err.Error())
	//} else {
	//	fmt.Println(len(res))
	//	fmt.Println(res)
	//}
}
func testRdbString(t *testing.T, rdb *Rdb) {
	key := "string"
	val := "string"
	if err := rdb.Set(context.Background(), key, val, 100*time.Second); err != nil {
		t.Fatal("set string err: ", err.Error())
	}
	get := ""
	if err := rdb.Get(context.Background(), key, &get); err != nil {
		t.Fatal("get string err: ", err.Error())
	}
	if val != get {
		t.Fatal("set get string err: ", get)
	} else {
		fmt.Println("testRdbString success")
	}
}
func testRdbString2(t *testing.T, rdb *Rdb) {
	key := "string"
	//val := "string"
	//if err := rdb.Set(context.Background(),key, val, 100*time.Second); err != nil {
	//	t.Fatal("set string err: ", err.Error())
	//}
	get := 0
	if err := rdb.Get(context.Background(), key, &get); err != nil {
		t.Fatal("get string err: ", err.Error())
	}
	//if val != get {
	//	t.Fatal("set get string err: ", get)
	//}else {
	//	fmt.Println("testRdbString success")
	//}
}

func testRdbInt(t *testing.T, rdb *Rdb) {
	key := "int"
	val := 10
	if err := rdb.Set(context.Background(), key, val, 100*time.Second); err != nil {
		t.Fatal("set int err: ", err.Error())
	}
	var get int
	if err := rdb.Get(context.Background(), key, &get); err != nil {
		t.Fatal("get int err: ", err.Error())
	}
	if val != get {
		t.Fatal("set get int err: ", get)
	}
}

func testRdbStruct(t *testing.T, rdb *Rdb) {
	key := "struct"
	type Student struct {
		Name string
		Age  int
	}
	val := Student{
		"zander",
		10,
	}
	if err := rdb.Set(context.Background(), key, val, 100*time.Second); err != nil {
		t.Fatal("set struct err: ", err.Error())
	}
	var get Student
	if err := rdb.Get(context.Background(), key, &get); err != nil {
		t.Fatal("get struct err: ", err.Error())
	}
	if get.Name != val.Name || get.Age != val.Age {
		t.Fatal("set get struct err: ", get)
	}
}

func testRdbMap(t *testing.T, rdb *Rdb) {
	key := "map"
	val := map[string]interface{}{
		"name": "zander",
		"age":  10,
	}
	if err := rdb.Set(context.Background(), key, val, 100*time.Second); err != nil {
		t.Fatal("set map err: ", err.Error())
	}
	var get map[string]interface{}
	if err := rdb.Get(context.Background(), key, &get); err != nil {
		t.Fatal("get map err: ", err.Error())
	}
}

func testRdbDel(t *testing.T, rdb *Rdb) {
	key := "maps"
	val := map[string]interface{}{
		"name": "zander",
		"age":  10,
	}
	if err := rdb.Set(context.Background(), key, val, 100*time.Second); err != nil {
		t.Fatal("set map err: ", err.Error())
	}

	if err := rdb.Delete(context.Background(), key); err != nil {
		t.Fatal("del  err: ", err.Error())
	}
}

func testRdbGetOrSet(t *testing.T, rdb *Rdb) {
	key := "str"
	val := struct {
		Name string
	}{
		Name: "zander",
	}
	type val2 struct {
		Name string
	}
	get := val2{}
	if err := rdb.GetOrSet(context.Background(), key, &get, 100*time.Second, func() (interface{}, error) {
		return val, nil
	}); err != nil {
		t.Fatal("Exists  err: ", err.Error())
	}

	if val.Name != get.Name {
		t.Fatal("set get string err: ", get)
	} else {
		fmt.Println("testRdbGetOrSet success")
	}
}
func testRdbGetOrSet2(t *testing.T, rdb *Rdb) {
	key := "str"
	val := struct {
		Name string
	}{
		Name: "zander",
	}
	type val2 struct {
		Name string
	}
	get := val2{}
	get2 := "&val2{}"
	if err := rdb.GetOrSet(context.Background(), key, &get2, 100*time.Second, func() (interface{}, error) {
		return val, nil
	}); err != nil {
		t.Fatal("Exists  err: ", err.Error())
	}

	if val.Name != get.Name {
		t.Fatal("set get string err: ", get)
	} else {
		fmt.Println("testRdbGetOrSet success")
	}
}

func testRdbExists(t *testing.T, rdb *Rdb) {
	key := "maps"
	if ok, err := rdb.Exists(context.Background(), key, key); err != nil {
		t.Fatal("Exists  err: ", err.Error())
	} else if !ok {

	}
}

func testRdbFlushDB(t *testing.T, rdb *Rdb) {
	if err := rdb.FlushDB(context.Background()); err != nil {
		t.Fatal("FlushDB  err: ", err.Error())
	}
}
