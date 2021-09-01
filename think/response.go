package think

import (
	"encoding/json"
)

var ResponseDebug = false

type Response struct {
	Code Code
	Msg  string
	Data interface{}
}

func ResponseSuccessData(data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"Data": data,
	}
}

func NewResponseFromErr(err error, debug bool) *Response {
	var r = Response{
		Code: Err2Code(err),
		Msg:  Err2Code(err).ToString(),
		Data: nil,
	}

	if r.Code != CodeSystemSpaceError {
		r.Data = err.Error()
	}
	if IsErrBiz(err) {
		r.Data = []string{
			BizCode(err),
			err.Error(),
		}
	}
	if debug {
		r.Data = err.Error()
	}
	return &r
}

func NewResponseBytesFromErr(err error, debug bool) []byte {
	res, _ := json.Marshal(NewResponseFromErr(err, debug))
	return res
}
