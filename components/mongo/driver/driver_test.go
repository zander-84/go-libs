package driver

import (
	"context"
	"flag"
	"testing"
)

// go test -v  -run TestMongo  -args 172.16.86.150  27017 dbtest
func TestMongo(t *testing.T) {
	if !flag.Parsed() {
		flag.Parse()
	}
	argList := flag.Args()
	user := ""
	if len(argList) > 3 {
		user = argList[3]
	}
	pwd := ""
	if len(argList) > 4 {
		pwd = argList[4]
	}

	mdb := NewMongo(Conf{
		Host:            argList[0],
		Port:            argList[1],
		MaxPoolSize:     100,
		MinPoolSize:     10,
		MaxConnIdleTime: 5,
		Database:        argList[2],
		User:            user,
		Pwd:             pwd,
	})

	if err := mdb.Start(); err != nil {
		t.Fatal(err.Error())
	}

	defer mdb.Stop()

	type Product struct {
		Code  string
		Price uint
	}

	if _, err := mdb.Collection("m_test").InsertOne(context.Background(), Product{"123", 32}); err != nil {
		t.Fatal(err.Error())
	}
}
