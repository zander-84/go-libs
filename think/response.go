package think

var ResponseDebug = false

type Response struct {
	Code    Code
	SubCode int `json:",omitempty"`
	Msg     string
	Data    interface{}
	Debug   interface{} `json:",omitempty"`
}

func (r *Response) Reset() {
	r.Code = 0
	r.Msg = ""
	r.Data = nil
	r.Debug = nil
}

func NewResponseFromErr(err error, debug bool) *Response {
	var code = Err2Code(err)
	var r = Response{
		Code: code,
		Msg:  code.ToString(),
		Data: nil,
	}

	if IsErrBiz(err) {
		r.SubCode = BizCode(err)
		r.Data = err.Error()
	} else if code == CodeSystemSpaceError {
		r.Data = nil
	} else {
		r.Data = err.Error()
	}

	if debug {
		r.Debug = err.Error()
	}
	return &r
}
