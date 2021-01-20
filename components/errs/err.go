package errs

import (
	"errors"
	"github.com/zander-84/go-libs/components/response/code"
)

type Err struct {
	level   Level
	path    string
	code    int
	from    string
	traceId string
	isLog   bool
	debug   bool
	err     error
}

func As(err error) *Err {
	if tmp, ok := err.(*Err); ok {
		return tmp
	}
	return nil
}

func (this *Err) Error() string             { return this.err.Error() }
func (this *Err) SetError(err error) *Err   { this.err = err; return this }
func (this *Err) SetPath(path string) *Err  { this.path = path; return this }
func (this *Err) SetFrom(from string) *Err  { this.from = from; return this }
func (this *Err) SetLog() *Err              { this.isLog = true; return this }
func (this *Err) SetNoLog() *Err            { this.isLog = false; return this }
func (this *Err) SetCode(code int) *Err     { this.code = code; return this }
func (this *Err) SetLevel(level Level) *Err { this.level = level; return this }
func (this *Err) SetTraceId(traceId string) *Err {
	if traceId != "" {
		this.traceId = traceId
	}
	return this
}

func (this *Err) SetDebug(debug bool) *Err { this.debug = debug; return this }

func (this *Err) SetLevelDebug() *Err { this.level = DebugLevel; return this }
func (this *Err) SetLevelInfo() *Err  { this.level = InfoLevel; return this }
func (this *Err) SetLevelWarn() *Err  { this.level = WarnLevel; return this }
func (this *Err) SetLevelError() *Err { this.level = ErrorLevel; return this }
func (this *Err) SetLevelPanic() *Err { this.level = PanicLevel; return this }
func (this *Err) SetLevelFatal() *Err { this.level = FatalLevel; return this }

func (this *Err) SetSuccessCode() *Err        { this.code = code.SuccessCode; return this }
func (this *Err) SetCodeSuccessAction() *Err  { this.code = code.SuccessActionCode; return this }
func (this *Err) SetCodeUserSpaceError() *Err { this.code = code.UserSpaceErrorCode; return this }
func (this *Err) SetCodeParamsError() *Err    { this.code = code.ParamsErrorCode; return this }
func (this *Err) SetCodeAlterError() *Err     { this.code = code.AlterErrorCode; return this }
func (this *Err) SetCodeForbiddenError() *Err {
	this.code = code.ForbiddenErrorCode
	return this
}
func (this *Err) SetCodeSignError() *Err { this.code = code.SignErrorCode; return this }
func (this *Err) SetCodeUserUnauthorizedError() *Err {
	this.code = code.UnauthorizedErrorCode
	return this
}
func (this *Err) SetCodeTooManyRequestsError() *Err {
	this.code = code.TooManyRequestsErrorCode
	return this
}
func (this *Err) SetCodeSystemSpaceError() *Err { this.code = code.SystemSpaceErrorCode; return this }
func (this *Err) SetCodeIgnoreError() *Err      { this.code = code.IgnoreErrorCode; return this }

func (this *Err) GetError() string   { return this.err.Error() }
func (this *Err) GetPath() string    { return this.path }
func (this *Err) GetFrom() string    { return this.from }
func (this *Err) GetCode() int       { return this.code }
func (this *Err) GetLevel() Level    { return this.level }
func (this *Err) GetTraceId() string { return this.traceId }
func (this *Err) IsLog() bool        { return this.isLog }
func (this *Err) GetDebug() bool     { return this.debug }

func New(errMsg string) *Err {
	err := errors.New(errMsg)
	return NewErr(err)
}

func NewErr(err error) *Err {
	return &Err{
		level: DefaultLevel,
		code:  code.UserSpaceErrorCode,
		err:   err,
		//traceId: helper.UniqueID(),
	}
}

// 只用于读
var (
	UninitializedError   = errors.New("组件未初始化")
	RecordNotFoundError  = errors.New("数据未找到")
	TypeError            = errors.New("错误类型")
	IDError              = errors.New("错误ID")
	IDExistError         = errors.New("ID不存在")
	SystemError          = errors.New("系统空间错误")
	SystemIgnoreError    = errors.New("忽略错误")
	UserError            = errors.New("用户空间错误")
	UnDefinitionError    = errors.New("未定义错误")
	CutomError           = errors.New("自定义错误")
	UnkonwError          = errors.New("未知错误")
	AlterError           = errors.New("提示错误")
	ParamsError          = errors.New("参数错误")
	TooManyRequestsError = errors.New("访问过于频繁")
	ForbiddenError       = errors.New("禁止访问")
	SignError            = errors.New("签名错误")
	UnauthorizedError    = errors.New("未授权")
	Success              = errors.New("成功")
	SuccessAction        = errors.New("成功行为")
)

type Level int8

const (
	DebugLevel Level = iota + 1
	InfoLevel
	WarnLevel
	ErrorLevel
	PanicLevel
	FatalLevel
	DefaultLevel = DebugLevel
)
