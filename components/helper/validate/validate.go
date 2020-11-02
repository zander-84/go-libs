package validate

import (
	"errors"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
	"strings"
)

// 兼容gin
type Validate struct {
	engine *validator.Validate
	trans  ut.Translator
	uni    *ut.UniversalTranslator
	conf   Conf
}

//
func NewValidate(conf Conf) *Validate {
	this := &Validate{
		engine: validator.New(),
		conf:   conf.SetDefault(),
	}

	if conf.Locale == "zh" {
		lzh := zh.New()
		uni := ut.New(lzh, lzh)
		this.trans, _ = uni.GetTranslator(conf.Locale)
		_ = zh_translations.RegisterDefaultTranslations(this.engine, this.trans)

	}
	this.engine.SetTagName(conf.ValidateTag) //validate binding
	this.engine.RegisterTagNameFunc(func(fld reflect.StructField) string {
		label := fld.Tag.Get(conf.CommentTag)
		if label == "" {
			return fld.Name
		}
		label = strings.Replace(label, ":", "\":\"", 1)
		label = strings.Replace(label, "：", "\":\"", 1)

		return label
	})

	return this
}

// ValidateStruct receives any kind of type, but only performed struct or pointer to struct type.
func (v *Validate) ValidateStruct(obj interface{}) error {
	if err := v.engine.Struct(obj); err != nil {
		errStr := "{"
		hasErr := false
		for _, err := range err.(validator.ValidationErrors) {
			errStr += "\"" + err.Translate(v.trans) + "\","
			hasErr = true
		}
		if hasErr {
			errStr = strings.Trim(errStr, ",")
			errStr += "}"
			return errors.New(errStr)
		} else {
			return err
		}
	}
	return nil
}

// Engine returns the underlying validator engine which powers the default
// Validator instance. This is useful if you want to register custom validations
// or struct level validations. See validator GoDoc for more info -
// https://godoc.org/gopkg.in/go-playground/validator.v10
func (v *Validate) Engine() *validator.Validate {
	return v.engine
}

func (v *Validate) lazyinit() {

}

func kindOfData(data interface{}) reflect.Kind {

	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}
