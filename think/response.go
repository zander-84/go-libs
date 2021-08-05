package think

import "fmt"

var ResponseDebug = false

type ResponseError struct {
	Code    Code
	SubCode int `json:",omitempty"`
	Msg     string
	Data    interface{}
}

func NewResponseFromErr(err error, debug bool) *ResponseError {

	var r = ResponseError{
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
	return []byte(NewResponseStringFromErr(err, debug))
}

var errNilString = fmt.Sprintf(`{"Code":%d,"Msg":"%s","Data":"%s"}`, CodeException, CodeException.ToString(), "")

func NewResponseStringFromErr(err error, debug bool) string {
	if err == nil {
		return errNilString
	}

	var data string
	var code = Err2Code(err)
	if code != CodeSystemSpaceError {
		data = err.Error()
	}
	if debug && data != "" {
		data = err.Error()
	}

	if IsErrBiz(err) {
		return fmt.Sprintf(`{"Code":%d,"Msg":"%s","SubCode":%d,"Data":"%s"}`, code, code.ToString(), BizCode(err), data)
	} else {
		return fmt.Sprintf(`{"Code":%d,"Msg":"%s","Data":"%s"}`, code, code.ToString(), data)
	}
}
