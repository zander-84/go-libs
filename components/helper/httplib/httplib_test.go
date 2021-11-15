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

	p := Post("http://ks-tmp.52haoka.com/dxjx/PoolService/api")
	p, _ = p.
		JSONBody(map[string]interface{}{
			"spcode": "18",
			"from":   "china",
			"method": "tyk-address",
		})

	if _, err := p.
		Debug(true).
		DumpBody(true).
		SetTimeout(7*time.Second, 10*time.Second).
		Response(); err != nil {
		t.Error(err.Error())
	} else {
		fmt.Println(p.String())
		fmt.Println(string(p.DumpRequest()))
		fmt.Println(p.String())

		//fmt.Println(res )
	}
}
