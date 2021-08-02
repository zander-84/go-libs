package think

import (
	"errors"
)

type Err error

// ErrRecordNotFound 数据未找到
type ErrRecordNotFound Err

var ErrInstanceRecordNotFound = ErrRecordNotFound(errors.New("record not found"))

func IsErrRecordNotFound(err error) bool {
	_, ok := err.(ErrRecordNotFound)
	return ok
}

// ErrRepeat 重复操作
type ErrRepeat Err

// ErrUnDone 未完成
type ErrUnDone Err

var ErrInstanceUnDone = ErrUnDone(errors.New("err undone"))

// ErrType 错误类型
type ErrType Err

// ErrIgnore 忽略错误
type ErrIgnore Err

// ErrUnknown 未知错误
type ErrUnknown Err

// ErrTooManyRequest 请求过多
type ErrTooManyRequest Err

// ErrForbidden 禁止访问
type ErrForbidden Err

// ErrSign 签名错误
type ErrSign Err

// ErrUnauthorized 未认证
type ErrUnauthorized Err

// ErrSystemSpace 系统空间
type ErrSystemSpace Err

// ErrAlter 简单错误 提示
type ErrAlter Err

// ErrParam 参数错误
type ErrParam Err

// ErrTimeOut 参数错误
type ErrTimeOut Err

// ErrBiz 业务错误
type ErrBiz Err

func ToErrBiz(subCode int, err error) error {
	return &errWithBiz{
		bizCode: subCode,
		cause:   err,
	}
}

func NewErrBiz(subCode int, errMsg string) error {
	return &errWithBiz{
		bizCode: subCode,
		cause:   errors.New(errMsg),
	}
}
func IsErrBiz(err error) bool {
	_, ok := err.(ErrBiz)
	return ok
}
func Err2Code(err error) Code {
	c, _ := errExplain(err)
	return c
}

func ErrString(err error) string {
	_, s := errExplain(err)
	return s
}
func errExplain(err error) (Code, string) {
	switch err.(type) {
	case ErrRecordNotFound:
		return CodeRecordNotFound, "ErrRecordNotFound"
	case ErrUnDone:
		return CodeUnDone, "ErrUnDone"
	case ErrIgnore:
		return CodeIgnore, "ErrIgnore"
	case ErrUnknown:
		return CodeUndefined, "ErrUnknown"
	case ErrTooManyRequest:
		return CodeTooManyRequests, "ErrTooManyRequest"
	case ErrForbidden:
		return CodeForbidden, "ErrForbidden"
	case ErrSign:
		return CodeSignError, "ErrSign"
	case ErrUnauthorized:
		return CodeUnauthorized, "ErrUnauthorized"
	case ErrSystemSpace:
		return CodeSystemSpaceError, "ErrSystemSpace"
	case ErrAlter:
		return CodeAlterError, "ErrAlter"
	case ErrType:
		return CodeTypeError, "ErrType"
	case ErrParam:
		return CodeParamError, "ErrParam"
	case ErrTimeOut:
		return CodeTimeOut, "ErrTimeOut"
	case ErrBiz:
		return CodeBizError, "ErrBiz"
	case ErrRepeat:
		return CodeRepeat, "ErrRepeat"
	default:
		return CodeUndefined, "error"
	}
}

type errWithBiz struct {
	bizCode int
	cause   error
}

type bizCode interface {
	BizCode() int
}

func (w *errWithBiz) Error() string { return w.cause.Error() }
func (w *errWithBiz) BizCode() int  { return w.bizCode }

// Unwrap provides compatibility for Go 1.13 error chains.
func (w *errWithBiz) Unwrap() error { return w.cause }

func BizCode(err error) int {
	if err == nil {
		return 0
	}
	wb, ok := err.(bizCode)
	if !ok {
		return 0
	}
	return wb.BizCode()
}
