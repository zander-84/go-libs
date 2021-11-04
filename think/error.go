package think

import (
	"errors"
	"fmt"
	"github.com/zander-84/go-libs/components/helper"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
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

// ErrNotFound 数据未找到
func ErrNotFound(err error) error {
	return &Error{code: CodeNotFound, err: errors.New(err.Error())}
}

var ErrInstanceRecordNotFound = ErrNotFound(errors.New("record not found"))

func IsThinkErr(err error) bool {
	e, ok := err.(*Error)
	if ok == false {
		return false
	}
	if e.code < MinCode {
		return false
	}
	return true
}

func IsErrNotFound(err error) bool {
	e, ok := err.(*Error)
	if ok == false {
		return false
	}
	if e.code == CodeNotFound {
		return true
	}
	return false
}

// ErrRepeat 重复操作
func ErrRepeat(err error) error {
	return &Error{code: CodeRepeat, err: errors.New(err.Error())}
}

// ErrUnDone 未完成
func ErrUnDone(err error) error {
	return &Error{code: CodeUnDone, err: errors.New(err.Error())}
}

var ErrInstanceUnDone = ErrUnDone(errors.New("err undone"))

// ErrType 错误类型
func ErrType(err error) error {
	return &Error{code: CodeTypeError, err: errors.New(err.Error())}
}

// ErrIgnore 忽略错误
func ErrIgnore(err error) error {
	return &Error{code: CodeIgnore, err: errors.New(err.Error())}
}

// ErrUnknown 未知错误
func ErrUnknown(err error) error {
	return &Error{code: CodeUndefined, err: errors.New(err.Error())}
}

// ErrTooManyRequest 请求过多
func ErrTooManyRequest(err error) error {
	return &Error{code: CodeTooManyRequests, err: errors.New(err.Error())}
}

// ErrForbidden 禁止访问
func ErrForbidden(err error) error {
	return &Error{code: CodeForbidden, err: errors.New(err.Error())}
}

// ErrSign 签名错误
func ErrSign(err error) error {
	return &Error{code: CodeSignError, err: errors.New(err.Error())}
}

// ErrUnauthorized 未认证
func ErrUnauthorized(err error) error {
	return &Error{code: CodeUnauthorized, err: errors.New(err.Error())}
}

// ErrSystemSpace 系统空间
func ErrSystemSpace(err error) error {
	return &Error{code: CodeSystemSpaceError, err: errors.New(err.Error())}
}

func IsErrSystemSpace(err error) bool {
	e, ok := err.(*Error)
	if ok == false {
		return false
	}
	if e.code == CodeSystemSpaceError {
		return true
	}
	return false
}

// ErrAlter 简单错误 提示
func ErrAlter(err error) error {
	return &Error{code: CodeAlterError, err: errors.New(err.Error())}
}

// ErrParam 参数错误
func ErrParam(err error) error {
	return &Error{code: CodeParamError, err: errors.New(err.Error())}
}

// ErrTimeOut 参数错误
func ErrTimeOut(err error) error {
	return &Error{code: CodeTimeOut, err: errors.New(err.Error())}
}

// ErrBiz 业务错误
func ErrBiz(subCode int, err error) error {
	return &Error{code: CodeBizError, bizCode: subCode, err: errors.New(err.Error())}
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
		return CodeSystemSpaceError
	}
	return e.code
}

func BizErrToString(bizErr error) string {
	return fmt.Sprintf("%d-%s", BizCode(bizErr), bizErr.Error())
}

// BizErrStringParse 解析业务错误
func BizErrStringParse(bizErrString string) error {
	arr := strings.Split(bizErrString, "-")
	subCode := helper.GetConv().ShouldStoI(arr[0])
	var errString string
	arr1 := arr[1:]
	errString = strings.Join(arr1, "-")

	return ErrBiz(subCode, errors.New(errString))
}

func ErrFromGrpc(err error) error {
	s, ok := status.FromError(err)
	if !ok {
		return ErrSystemSpace(err)
	}

	if s.Code() == codes.Unavailable {
		return &Error{code: CodeUnavailable, err: errors.New(CodeUnavailable.ToString())}
	}

	if uint32(s.Code()) < uint32(MinCode) {
		return ErrSystemSpace(err)
	} else {
		return &Error{code: Code(s.Code()), err: errors.New(s.Message())}
	}
}
