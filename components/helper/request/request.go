package request

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

//r.Header.Set("Content-Type", "application/xml")
//r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
//r.Header.Set("Content-Type", "application/json")
type HttpRequestParams struct {
	Client        *http.Client
	Method        string
	Url           string
	RequestHandle func(r *http.Request)
	Body          io.Reader
	Debug         bool
	Attempt       int // 推荐get使用
}

func UrlValue(obj url.Values) string {
	if obj == nil {
		return ""
	}

	return obj.Encode()
}

func BodyValue(obj url.Values) io.Reader {
	if obj == nil {
		return nil
	}
	return strings.NewReader(obj.Encode())
}

func BodyMap(obj map[string]interface{}) (io.Reader, error) {
	return BodyJSON(obj)
}

func BodyJSON(obj interface{}) (io.Reader, error) {
	if obj == nil {
		return nil, nil
	}

	byts, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(byts), nil
}

func BodyXML(obj interface{}) (io.Reader, error) {
	if obj == nil {
		return nil, nil
	}

	byts, err := xml.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(byts), nil
}

type HttpResponse struct {
	StatusCode int
	Header     http.Header
	Body       []byte
	Dump       []byte
	Err        error
}

func HttpRequest(requestParams HttpRequestParams) *HttpResponse {
	resp := httpRequest(requestParams)
	if resp.Err == nil {
		return resp
	}
	if requestParams.Attempt > 0 {
		for i := 0; i < requestParams.Attempt; i++ {
			resp = httpRequest(requestParams)
			if resp.Err == nil {
				break
			}
		}
	}
	return resp
}

func httpRequest(requestParams HttpRequestParams) *HttpResponse {
	if requestParams.Client == nil {
		requestParams.Client = new(http.Client)
		requestParams.Client.Timeout = 2 * 60 * time.Second
	}

	httpResponse := new(HttpResponse)
	var request *http.Request
	request, httpResponse.Err = http.NewRequest(strings.ToUpper(requestParams.Method), requestParams.Url, requestParams.Body)
	if httpResponse.Err != nil {
		return httpResponse
	}

	if requestParams.RequestHandle != nil {
		requestParams.RequestHandle(request)
	}
	request.Close = true

	var response *http.Response
	response, httpResponse.Err = requestParams.Client.Do(request)
	if response != nil {
		httpResponse.Header = response.Header
		httpResponse.StatusCode = response.StatusCode
		defer response.Body.Close()
	}

	if requestParams.Debug {
		httpResponse.Dump, _ = httputil.DumpRequest(request, true)
	}

	if httpResponse.Err != nil {
		return httpResponse
	} else {
		if response != nil {
			httpResponse.Body, httpResponse.Err = ioutil.ReadAll(response.Body)
			return httpResponse
		} else {
			httpResponse.Err = errors.New("no response")
			return httpResponse
		}
	}
}

func (b *HttpResponse) ToJSON(v interface{}) error {
	if b.Err != nil {
		return b.Err
	}

	return json.Unmarshal(b.Body, v)
}

// ToXML returns the map that marshals from the body bytes as xml in response .
// it calls Response inner.
func (b *HttpResponse) ToXML(v interface{}) error {
	if b.Err != nil {
		return b.Err
	}

	return xml.Unmarshal(b.Body, v)
}

// ToYAML returns the map that marshals from the body bytes as yaml in response .
// it calls Response inner.
func (b *HttpResponse) ToYAML(v interface{}) error {
	if b.Err != nil {
		return b.Err
	}

	return yaml.Unmarshal(b.Body, v)
}
