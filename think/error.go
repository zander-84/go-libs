package think

import (
	"errors"
)

type Error struct {
	code    Code
	bizCode int
	err     error
}
type bizCode interface {
	BizCode() int
}

func (e *Error) BizCode() int { return e.bizCode }

// Unwrap provides compatibility for Go 1.13 error chains.
func (e *Error) Unwrap() error { return e.err }

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

func (e *Error) Error() string {
	return e.err.Error()
}

// ErrRecordNotFound 数据未找到
func ErrRecordNotFound(err error) error {
	return &Error{code: CodeRecordNotFound, err: err}
}

var ErrInstanceRecordNotFound = ErrRecordNotFound(errors.New("record not found"))

func IsErrRecordNotFound(err error) bool {
	e, ok := err.(*Error)
	if ok == false {
		return false
	}
	if e.code == CodeRecordNotFound {
		return true
	}
	return false
}

// ErrRepeat 重复操作
func ErrRepeat(err error) error {
	return &Error{code: CodeRepeat, err: err}
}

// ErrUnDone 未完成
func ErrUnDone(err error) error {
	return &Error{code: CodeUnDone, err: err}
}

var ErrInstanceUnDone = ErrUnDone(errors.New("err undone"))

// ErrType 错误类型
func ErrType(err error) error {
	return &Error{code: CodeTypeError, err: err}
}

// ErrIgnore 忽略错误
func ErrIgnore(err error) error {
	return &Error{code: CodeIgnore, err: err}
}

// ErrUnknown 未知错误
func ErrUnknown(err error) error {
	return &Error{code: CodeUndefined, err: err}
}

// ErrTooManyRequest 请求过多
func ErrTooManyRequest(err error) error {
	return &Error{code: CodeTooManyRequests, err: err}
}

// ErrForbidden 禁止访问
func ErrForbidden(err error) error {
	return &Error{code: CodeForbidden, err: err}
}

// ErrSign 签名错误
func ErrSign(err error) error {
	return &Error{code: CodeSignError, err: err}
}

// ErrUnauthorized 未认证
func ErrUnauthorized(err error) error {
	return &Error{code: CodeUnauthorized, err: err}
}

// ErrSystemSpace 系统空间
func ErrSystemSpace(err error) error {
	return &Error{code: CodeSystemSpaceError, err: err}
}

// ErrAlter 简单错误 提示
func ErrAlter(err error) error {
	return &Error{code: CodeAlterError, err: err}
}

// ErrParam 参数错误
func ErrParam(err error) error {
	return &Error{code: CodeParamError, err: err}
}

// ErrTimeOut 参数错误
func ErrTimeOut(err error) error {
	return &Error{code: CodeTimeOut, err: err}
}

// ErrBiz 业务错误
func ErrBiz(subCode int, err error) error {
	return &Error{code: CodeBizError, bizCode: subCode, err: err}
}

func IsErrBiz(err error) bool {
	e, ok := err.(*Error)
	if ok == false {
		return false
	}
	if e.code == CodeBizError {
		return true
	}
	return false
}
func Err2Code(err error) Code {
	e, ok := err.(*Error)
	if ok == false {
		return CodeUndefined
	}
	return e.code
}
