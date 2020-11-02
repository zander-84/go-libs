package request

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
)

// go test -v  -run TestHttpRequest
func TestHttpRequest(t *testing.T) {

	resp := HttpRequest(HttpRequestParams{
		Client: &http.Client{
			Transport:     nil,
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       0,
		},
		Method: "post",
		Url:    "http://test.local/test.php?name=zander&age=9",
		RequestHandle: func(r *http.Request) {
			r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		},
		Body: BodyValue(url.Values{
			"name": []string{"zander"},
			"age":  []string{"18"},
		}),
		Debug:   true,
		Attempt: 2,
	})
	t.Log(string(resp.Dump))
	if resp.Err != nil {
		t.Error(resp.Err.Error())
	} else {
		fmt.Println(string(resp.Body))
	}
}
