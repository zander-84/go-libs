package errs

type SimpleError string

func (this SimpleError) Error() string { return string(this) }
func (this SimpleError) SetError(err string) {
}

const (
	UninitializedError  = SimpleError("组件未初始化")
	RecordNotFoundError = SimpleError("错误记录")
	TypeError           = SimpleError("错误类型")
	IDError             = SimpleError("错误ID")
	IDExistError        = SimpleError("ID不存在")

	SystemError          = SimpleError("系统空间错误")
	UserError            = SimpleError("用户空间错误")
	UnDefinitionError    = SimpleError("未定义错误")
	AlterError           = SimpleError("提示错误")
	ParamsError          = SimpleError("参数错误")
	TooManyRequestsError = SimpleError("访问过于频繁")
	ForbiddenError       = SimpleError("禁止访问")
	UnauthorizedError    = SimpleError("未授权")
	Success              = SimpleError("成功")
	SuccessAction        = SimpleError("成功行为")
)

type Level int8

const (
	SuccessLevel Level = iota + 1

	DebugLevel
	InfoLevel

	// 用户空间
	WarnLevel
	ErrorLevel

	//系统空间
	PanicLevel
	FatalLevel

	DefaultLevel = SuccessLevel
)

type CustomError struct {
	Err   string
	Level Level
}

func IsCustomError(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(CustomError)
	return ok
}

func ToCustomError(err error) CustomError {
	e, _ := err.(CustomError)
	return e
}

func NewDefaultError(err string) CustomError {
	return CustomError{
		Err:   err,
		Level: SuccessLevel,
	}
}

func NewCustomError(err string, level Level) CustomError {
	return CustomError{
		Err:   err,
		Level: level,
	}
}

func NewSuccess(err string) CustomError {
	return CustomError{
		Err:   err,
		Level: ErrorLevel,
	}
}

func NewUserSpaceError(err string) CustomError {
	return CustomError{
		Err:   err,
		Level: SuccessLevel,
	}
}

func NewSystemSpaceError(err string) CustomError {
	return CustomError{
		Err:   err,
		Level: FatalLevel,
	}
}

func (this CustomError) Error() string { return this.Err }

func (this CustomError) SetError(err string) CustomError {
	this.Err = err
	return this
}
func (this CustomError) SetLevel(l Level) CustomError {
	this.Level = l
	return this
}

func (this CustomError) IsSystemSpaceError() bool {
	if this.Level >= PanicLevel {
		return true
	}
	return false
}

func (this CustomError) IsUserSpaceError() bool {
	if this.Level <= ErrorLevel && this.Level >= WarnLevel {
		return true
	}
	return false
}

func (this CustomError) IsSuccess() bool {
	if this.Level <= InfoLevel {
		return true
	}
	return false
}
