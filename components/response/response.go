package response

import (
	"github.com/zander-84/go-libs/components/errs"
	"github.com/zander-84/go-libs/components/response/code"
	"net/http"
)

var Debug = false

type Data struct {
	HttpCode  int `json:"-"`
	Code      int
	Msg       interface{}
	Data      interface{}
	Debug     interface{} `json:",omitempty"`
	TraceID   string      `json:",omitempty"`
	ShowDebug bool        `json:"-"`
}

func NewData() *Data {
	return new(Data)
}

func (this *Data) SetData(data interface{}) *Data  { this.Data = data; return this }
func (this *Data) SetDebug(debug bool) *Data       { this.ShowDebug = debug; return this }
func (this *Data) SetTraceID(traceID string) *Data { this.TraceID = traceID; return this }

func Success(data *Data) *Data {
	data.HttpCode = http.StatusOK
	data.Code = code.SuccessCode
	if data.Msg == nil {
		data.Msg = errs.Success.Error()
	}
	return data
}

func SuccessAction(data *Data) *Data {
	data.HttpCode = http.StatusOK
	data.Code = code.SuccessActionCode
	if data.Msg == "" {
		data.Msg = errs.SuccessAction.Error()
	}
	return data
}

func SystemSpaceError(data *Data) *Data {
	data.HttpCode = http.StatusInternalServerError
	data.Code = code.SystemSpaceErrorCode
	if data.Msg == nil {
		data.Msg = errs.SystemError.Error()
	}
	if Debug {
		data.Debug = data.Data
	}
	if data.ShowDebug {
		data.Debug = data.Data
	}
	data.Data = "服务异常"
	return data
}

func UnknowError(data *Data) *Data {
	data.HttpCode = http.StatusInternalServerError
	data.Code = code.UnknownErrorCode
	if data.Msg == nil {
		data.Msg = errs.UnkonwError.Error()
	}
	if data.ShowDebug {
		data.Debug = data.Data
	}
	return data
}

func IgnoreError(data *Data) *Data {
	data.HttpCode = http.StatusInternalServerError
	data.Code = code.IgnoreErrorCode
	if data.Msg == nil {
		data.Msg = errs.SystemIgnoreError.Error()
	}
	if Debug {
		data.Debug = data.Data
	}
	if data.ShowDebug {
		data.Debug = data.Data
	}
	data.Data = nil
	return data
}

func UserSpaceError(data *Data) *Data {
	data.HttpCode = http.StatusBadRequest
	data.Code = code.UserSpaceErrorCode

	if data.Msg == nil {
		data.Msg = errs.UnkonwError.Error()
	}
	if data.ShowDebug {
		data.Debug = data.Data
	}
	return data
}

func AlterError(data *Data) *Data {
	data.HttpCode = http.StatusBadRequest
	data.Code = code.AlterErrorCode
	if data.Msg == nil {
		data.Msg = errs.AlterError.Error()
	}
	if data.ShowDebug {
		data.Debug = data.Data
	}
	return data
}

func ParamsError(data *Data) *Data {
	data.HttpCode = http.StatusBadRequest
	data.Code = code.ParamsErrorCode
	if data.Msg == nil {
		data.Msg = errs.ParamsError.Error()
	}
	if data.ShowDebug {
		data.Debug = data.Data
	}
	return data
}
func ForbiddenError(data *Data) *Data {
	data.HttpCode = http.StatusForbidden
	data.Code = code.ForbiddenErrorCode
	if data.Msg == nil {
		data.Msg = errs.ForbiddenError.Error()
	}
	if data.ShowDebug {
		data.Debug = data.Data
	}
	return data
}
func SignError(data *Data) *Data {
	data.HttpCode = http.StatusForbidden
	data.Code = code.SignErrorCode
	if data.Msg == nil {
		data.Msg = errs.SignError.Error()
	}
	if data.ShowDebug {
		data.Debug = data.Data
	}
	return data
}
func TooManyRequestsError(data *Data) *Data {
	data.HttpCode = http.StatusTooManyRequests
	data.Code = code.TooManyRequestsErrorCode
	if data.Msg == nil {
		data.Msg = errs.TooManyRequestsError.Error()
	}
	if data.ShowDebug {
		data.Debug = data.Data
	}
	return data
}
func UnauthorizedError(data *Data) *Data {
	data.HttpCode = http.StatusUnauthorized
	data.Code = code.UnauthorizedErrorCode
	if data.Msg == nil {
		data.Msg = errs.UnauthorizedError.Error()
	}
	if data.ShowDebug {
		data.Debug = data.Data
	}
	return data
}

func Custom(data *Data) *Data {
	if data.HttpCode == 0 {
		data.HttpCode = http.StatusBadRequest
	}
	if data.Msg == nil {
		data.Msg = errs.CutomError.Error()
	}
	if data.ShowDebug {
		data.Debug = data.Data
	}
	return data
}
