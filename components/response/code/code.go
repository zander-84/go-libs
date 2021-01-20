package code

import "net/http"

const (
	//系统code码预留
	SuccessCode = 100200 // 成功

	SuccessActionCode = 100201 // 成功行为

	UserSpaceErrorCode = 100400 // 用户空间错误

	ParamsErrorCode = 101400 // 参数错误

	AlterErrorCode = 102400 // 简单错误

	ForbiddenErrorCode = 100403 // 禁止访问

	SignErrorCode = 101403 // 签名错误

	UnauthorizedErrorCode = 100401 // 未认证

	TooManyRequestsErrorCode = 100429 // 请求过于频繁

	SystemSpaceErrorCode = 100500 // 系统空间错误

	IgnoreErrorCode = 101500 // 系统空间错误不错处理

	UnknownErrorCode = 102500 // 未知错误

)

var errCode2HttpCode = map[int]int{
	SuccessCode:              http.StatusOK,
	SuccessActionCode:        http.StatusOK,
	UserSpaceErrorCode:       http.StatusBadRequest,
	ParamsErrorCode:          http.StatusBadRequest,
	AlterErrorCode:           http.StatusBadRequest,
	ForbiddenErrorCode:       http.StatusForbidden,
	SignErrorCode:            http.StatusForbidden,
	UnauthorizedErrorCode:    http.StatusUnauthorized,
	TooManyRequestsErrorCode: http.StatusTooManyRequests,
	SystemSpaceErrorCode:     http.StatusInternalServerError,
	IgnoreErrorCode:          http.StatusInternalServerError,
	UnknownErrorCode:         http.StatusInternalServerError,
}

func ErrCode2HttpCode(errCode int) int {
	if httpCode, ok := errCode2HttpCode[errCode]; ok {
		return httpCode
	}

	return http.StatusInternalServerError
}
