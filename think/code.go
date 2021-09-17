package think

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Code uint32

const (
	CodeSuccess Code = 100200 // 成功

	CodeSuccessAction Code = 100201 // 成功行为

	CodeBizError Code = 100400 // 业务错误 用户空间错误

	CodeAlterError Code = 102400 // 简单错误

	CodeParamError Code = 101400 // 参数错误

	CodeNotFound Code = 100404 // 记录错误

	CodeRepeat Code = 100405 // 重复操作

	CodeUnDone Code = 103400 // 记录错误

	CodeForbidden Code = 100403 // 禁止访问

	CodeSignError Code = 101403 // 签名错误

	CodeUnauthorized Code = 100401 // 未认证

	CodeTooManyRequests Code = 100429 // 请求过于频繁

	CodeSystemSpaceError Code = 100500 // 系统空间错误  不外抛

	CodeIgnore Code = 101500 // 忽略

	CodeUndefined Code = 102500 // 未定义

	CodeTimeOut Code = 102504 // 超时

	CodeException Code = 103500 // 异常

	CodeTypeError Code = 104500 // 类型错误

)

func (c Code) HttpCode() int {
	switch c {
	case CodeSuccess, CodeSuccessAction:
		return http.StatusOK

	case CodeParamError, CodeNotFound, CodeAlterError, CodeBizError, CodeUnDone, CodeRepeat:
		return http.StatusBadRequest

	case CodeForbidden, CodeSignError:
		return http.StatusForbidden

	case CodeUnauthorized:
		return http.StatusUnauthorized

	case CodeTooManyRequests:
		return http.StatusTooManyRequests

	case CodeSystemSpaceError, CodeIgnore, CodeUndefined, CodeException, CodeTypeError:
		return http.StatusInternalServerError

	case CodeTimeOut:
		return http.StatusGatewayTimeout

	default:
		return http.StatusInternalServerError

	}
}

func (c Code) ToString() string {
	switch c {
	case CodeSuccess, CodeSuccessAction:
		return "成功"
	case CodeParamError:
		return "参数错误"
	case CodeNotFound:
		return "404"
	case CodeRepeat:
		return "重复操作"
	case CodeAlterError:
		return "提示错误"
	case CodeBizError:
		return "业务错误"
	case CodeUnDone:
		return "未完成"
	case CodeForbidden:
		return "禁止访问"
	case CodeSignError:
		return "签名错误"
	case CodeUnauthorized:
		return "未认证"
	case CodeTooManyRequests:
		return "反问过于频繁"
	case CodeSystemSpaceError:
		return "系统空间错误"
	case CodeIgnore:
		return "忽略"
	case CodeUndefined:
		return "未定义"
	case CodeTimeOut:
		return "超时"
	case CodeException:
		return "异常"
	case CodeTypeError:
		return "错误类型"
	default:
		return fmt.Sprintf("未定义： %d", c)
	}
}

func HttpWriteJson(w http.ResponseWriter, code int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	var err error

	switch v := data.(type) {
	case string:
		_, err = w.Write([]byte(v))
	case []byte:
		_, err = w.Write(v)
	default:
		err = json.NewEncoder(w).Encode(data)
	}
	return err
}
