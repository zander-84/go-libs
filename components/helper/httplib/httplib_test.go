package httplib

import (
	"fmt"
	"testing"
	"time"
)

func TestHttpRequest(t *testing.T) {
	//if res,err:=Get("http://test.local/test.php").
	//	Param("name","zander").
	//	Debug(true).
	//	String();err!=nil{
	//	t.Error(err.Error())
	//}else {
	//	fmt.Println(res)
	//}
	//g := Get("http://test.local/test.php")
	//if _, err := g.
	//	Param("name", "zander").
	//	Debug(true).
	//	DumpBody(false).
	//	Response(); err != nil {
	//	t.Error(err.Error())
	//} else {
	//	fmt.Println(string(g.DumpRequest()))
	//	//fmt.Println(res)
	//}

	p := Get("http://test.local/test.php")
	p, _ = p.Param("name", "abc").
		JSONBody(map[string]interface{}{
			"age":  18,
			"from": "china",
		})

	if _, err := p.
		Debug(true).
		DumpBody(true).
		SetTimeout(7*time.Second, 10*time.Second).
		Response(); err != nil {
		t.Error(err.Error())
	} else {
		fmt.Println(string(p.DumpRequest()))
		//fmt.Println(res )
	}
}
