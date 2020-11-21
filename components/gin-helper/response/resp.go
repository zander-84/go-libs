package response

import (
	"github.com/gin-gonic/gin"
	"github.com/zander-84/go-libs/components/helper"
	"github.com/zander-84/go-libs/components/response"
)

var GinTraceId = "gin_trace_id"

func Success(c *gin.Context, dat interface{}) {
	data := response.NewData()
	data.Data = dat
	data.TraceID = GetTraceId(c)
	SuccessData(c, data)
}

func SuccessData(c *gin.Context, data *response.Data) {
	response.Success(data)
	c.JSON(data.HttpCode, data)
}

func SuccessAction(c *gin.Context) {
	data := response.NewData()
	SuccessActionData(c, data)
	return
}

func SuccessActionData(c *gin.Context, data *response.Data) {
	response.SuccessAction(data)
	c.JSON(data.HttpCode, data)
}

func SystemSpaceError(c *gin.Context, msg string) {
	re := response.NewData()
	re.Msg = "服务器出了点小差，请稍后重试！"
	re.TraceID = GetTraceId(c)
	re.Debug = msg
	SystemSpaceErrorData(c, re)
}

func SystemSpaceErrorData(c *gin.Context, data *response.Data) {
	response.SystemSpaceError(data)
	c.JSON(data.HttpCode, data)
}

func UserSpaceError(c *gin.Context, dat interface{}) {
	data := response.NewData()
	data.Data = dat
	data.TraceID = GetTraceId(c)
	UserSpaceErrorData(c, data)
}

func UserSpaceErrorData(c *gin.Context, data *response.Data) {
	response.UserSpaceError(data)
	c.JSON(data.HttpCode, data)
}

func UserAlterError(c *gin.Context, msg string) {
	data := response.NewData()
	data.Msg = msg
	data.TraceID = GetTraceId(c)
	UserAlterErrorData(c, data)
}

func UserAlterErrorData(c *gin.Context, data *response.Data) {
	response.UserAlterError(data)
	c.JSON(data.HttpCode, data)
}

func UserParamsError(c *gin.Context, dat interface{}) {
	data := response.NewData()
	data.Data = dat
	data.TraceID = GetTraceId(c)
	UserParamsErrorData(c, data)
}

func UserParamsErrorData(c *gin.Context, data *response.Data) {
	response.UserParamsError(data)
	c.JSON(data.HttpCode, data)
}

// 后台Rbac认证
func UserForbiddenRbacError(c *gin.Context, dat interface{}) {
	data := response.NewData()
	data.Msg = "当前用户无权限访问！"
	data.Data = dat
	data.TraceID = GetTraceId(c)
	UserForbiddenErrorData(c, data)
}

// 签名认证
func UserForbiddenSignError(c *gin.Context, dat interface{}) {
	data := response.NewData()
	data.Msg = "当前用户签名验证未通过！"
	data.Data = dat
	data.TraceID = GetTraceId(c)
	UserForbiddenErrorData(c, data)
}

func UserForbiddenErrorData(c *gin.Context, data *response.Data) {
	response.UserForbiddenError(data)
	c.JSON(data.HttpCode, data)
}

func UserSignError(c *gin.Context, dat interface{}) {
	data := response.NewData()
	data.Data = dat
	data.TraceID = GetTraceId(c)
	UserSignErrorData(c, data)
}

func UserSignErrorData(c *gin.Context, data *response.Data) {
	response.UserSignError(data)
	c.JSON(data.HttpCode, data)
}

func UserTooManyRequestError(c *gin.Context, dat interface{}) {
	data := response.NewData()
	data.Data = dat
	data.TraceID = GetTraceId(c)
	UserTooManyRequestErrorData(c, data)
}

func UserTooManyRequestErrorData(c *gin.Context, data *response.Data) {
	// data.ShowDebug = true
	response.UserTooManyRequestsError(data)
	c.JSON(data.HttpCode, data)
}

func UserUnauthorizedError(c *gin.Context, dat interface{}) {
	data := response.NewData()
	data.Data = dat
	data.TraceID = GetTraceId(c)
	UserUnauthorizedErrorData(c, data)
}

func UserUnauthorizedErrorData(c *gin.Context, data *response.Data) {
	// data.ShowDebug = true
	response.UserUnauthorizedError(data)
	c.JSON(data.HttpCode, data)
}

func CustomError(c *gin.Context, dat interface{}) {
	data := response.NewData()
	data.Data = dat
	data.TraceID = GetTraceId(c)
	CustomErrorData(c, data)
}

func CustomErrorData(c *gin.Context, data *response.Data) {
	// data.ShowDebug = true
	response.CustomError(data)
	c.JSON(data.HttpCode, data)
}

func GetTraceId(ctx *gin.Context) string {
	tid, exits := ctx.Get(GinTraceId)
	if exits {
		return tid.(string)
	}
	return ""
}

func SetTraceId(ctx *gin.Context, prefix string) {
	ctx.Set(GinTraceId, helper.UniqueID(prefix))
}
