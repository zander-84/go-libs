package think

import (
	"encoding/json"
	"errors"
	"testing"
)

//BenchmarkNewResponseBytesFromErr-16    	 2676440	       454.8 ns/op
func BenchmarkNewResponseBytesFromErr(b *testing.B) {
	var t = errors.New("alert")
	for i := 0; i < b.N; i++ {
		json.Marshal(NewResponseFromErr(t, false))
	}
}

func responseSuccessData(data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"Data": data,
	}
}

//BenchmarkResponseSuccessData-16    	 1499968	       748.2 ns/op
func BenchmarkResponseSuccessData(b *testing.B) {
	for i := 0; i < b.N; i++ {
		json.Marshal(responseSuccessData(1234))
	}
}

type successResp struct {
	Data interface{}
}

func responseSuccessData2(data interface{}) *successResp {
	s := new(successResp)
	s.Data = data
	return s
}

//BenchmarkResponseSuccessData2-16    	 5015295	       236.0 ns/op
func BenchmarkResponseSuccessData2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		json.Marshal(responseSuccessData2(1234))
	}
}

func responseSuccessData3(data interface{}) []interface{} {
	return []interface{}{data}
}

//BenchmarkResponseSuccessData3-16    	 4575492	       260.5 ns/op
func BenchmarkResponseSuccessData3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		json.Marshal(responseSuccessData3("1234"))
	}
}
