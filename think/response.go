package think

import "fmt"

var ResponseDebug = false

var ResponseKeyError = "Error"
var ResponseKeyMsg = "Msg"
var ResponseKeyData = "Data"
var ResponseKeySubCode = "SubCode"

func ResponseSuccessData(data interface{}) map[string]interface{} {
	return map[string]interface{}{
		ResponseKeyData: data,
	}
}

func ResponseSuccessSource(data interface{}) interface{} {
	return data
}

func ResponseSuccessDataWithSubCode(subCode int32, data interface{}) map[string]interface{} {
	return map[string]interface{}{
		ResponseKeySubCode: subCode,
		ResponseKeyData:    data,
	}
}

type ResponseError struct {
	Err     Code
	SubCode int `json:",omitempty"`
	Msg     string
	Data    interface{}
}

func NewResponseFromErr(err error, debug bool) *ResponseError {

	var r = ResponseError{
		Err:     Err2Code(err),
		SubCode: 0,
		Msg:     Err2Code(err).ToString(),
		Data:    nil,
	}

	if r.Err != CodeSystemSpaceError {
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

var errNilString = fmt.Sprintf(`{"%s":%d,"%s":"%s","%s":"%s"}`, ResponseKeyError, CodeException, ResponseKeyMsg, CodeException.ToString(), ResponseKeyData, "")

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
		return fmt.Sprintf(`{"%s":%d,"%s":"%s","%s":%d,"%s":"%s"}`, ResponseKeyError, code, ResponseKeyMsg, code.ToString(), ResponseKeySubCode, BizCode(err), ResponseKeyData, data)
	} else {
		return fmt.Sprintf(`{"%s":%d,"%s":"%s","%s":"%s"}`, ResponseKeyError, code, ResponseKeyMsg, code.ToString(), ResponseKeyData, data)
	}
}
