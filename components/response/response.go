package response

import (
	"github.com/zander-84/go-libs/components/errs"
	"net/http"
)

type Data struct {
	HttpCode  int         `json:"-"`
	Code      string      `json:"code"`
	Msg       interface{} `json:"msg"`
	Data      interface{} `json:"data"`
	TraceID   interface{} `json:"trace_id,omitempty"`
	Debug     interface{} `json:"debug,omitempty"`
	ShowDebug bool        `json:"-"`
}

const (
	//系统code码预留
	SuccessCode                  = "10200" // 成功
	SuccessActionCode            = "10201" // 成功行为
	UserSpaceErrorCode           = "10400" // 用户空间错误
	UserAlterErrorCode           = "10406" //  简单错误
	UserParamsErrorCode          = "10402" // 参数错误
	UserForbiddenErrorCode       = "10403" // 禁止访问
	UserUnauthorizedErrorCode    = "10401" // 未认证
	UserTooManyRequestsErrorCode = "10429" //请求过于频繁
	SystemSpaceErrorCode         = "10500" //系统空间错误
)

func Success(data Data) Data {
	data.HttpCode = http.StatusOK
	data.Code = SuccessCode
	if data.Msg == "" {
		data.Msg = string(errs.Success)
	}
	HandleData(&data)
	return data
}

func SuccessAction(data Data) Data {
	data.HttpCode = http.StatusNoContent
	data.Code = SuccessActionCode
	if data.Msg == "" {
		data.Msg = string(errs.SuccessAction)
	}
	HandleData(&data)
	return data
}

func SystemSpaceError(data Data) Data {
	data.HttpCode = http.StatusInternalServerError
	data.Code = SystemSpaceErrorCode
	if data.Msg == "" {
		data.Msg = string(errs.SystemError)
	}
	HandleData(&data)
	return data
}

func UserSpaceError(data Data) Data {
	data.HttpCode = http.StatusBadRequest
	data.Code = UserSpaceErrorCode

	if data.Msg == "" {
		data.Msg = string(errs.UserError)
	}
	HandleData(&data)
	return data
}

func UserAlterError(data Data) Data {
	data.HttpCode = http.StatusBadRequest
	data.Code = UserAlterErrorCode
	if data.Msg == "" {
		data.Msg = string(errs.AlterError)
	}
	HandleData(&data)
	return data
}

func UserParamsError(data Data) Data {
	data.HttpCode = http.StatusBadRequest
	data.Code = UserParamsErrorCode
	if data.Msg == "" {
		data.Msg = string(errs.ParamsError)
	}
	HandleData(&data)
	return data
}
func UserForbiddenError(data Data) Data {
	data.HttpCode = http.StatusForbidden
	data.Code = UserForbiddenErrorCode
	if data.Msg == "" {
		data.Msg = string(errs.ForbiddenError)
	}
	HandleData(&data)
	return data
}

func UserTooManyRequestsError(data Data) Data {
	data.HttpCode = http.StatusTooManyRequests
	data.Code = UserTooManyRequestsErrorCode
	if data.Msg == "" {
		data.Msg = string(errs.TooManyRequestsError)
	}
	HandleData(&data)
	return data
}
func UserUnauthorizedError(data Data) Data {
	data.HttpCode = http.StatusUnauthorized
	data.Code = UserUnauthorizedErrorCode
	if data.Msg == "" {
		data.Msg = string(errs.UnauthorizedError)
	}
	HandleData(&data)
	return data
}

func CustomError(data Data) Data {
	if data.HttpCode == 0 {
		data.HttpCode = http.StatusBadRequest
	}
	if data.Msg == "" {
		data.Msg = string(errs.UnDefinitionError)
	}
	HandleData(&data)
	return data
}

var Debug = false

type DebugRecordInterface interface {
	Record(data *Data)
}
type DebugRecord func(data *Data)

func (this DebugRecord) Record(data *Data) {
	this(data)
}

var DebugRecordEngine DebugRecordInterface

func HandleData(data *Data) {
	if DebugRecordEngine != nil {
		DebugRecordEngine.Record(data)
	}

	if data.ShowDebug {
		return
	}
	if Debug {
		return
	}

	data.Debug = nil
	return
}
