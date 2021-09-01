package think

import (
	"encoding/json"
)

var ResponseDebug = false

type Response struct {
	Code    Code
	SubCode int `json:",omitempty"`
	Msg     string
	Data    interface{}
}

func NewResponseFromErr(err error, debug bool) *Response {
	var r = Response{
		Code:    Err2Code(err),
		SubCode: 0,
		Msg:     Err2Code(err).ToString(),
		Data:    nil,
	}

	if r.Code != CodeSystemSpaceError {
		r.Data = err.Error()
	}
	if IsErrBiz(err) {
		r.SubCode = BizCode(err)
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
