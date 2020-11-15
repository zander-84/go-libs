package zlog

import (
	"github.com/zander-84/go-libs/components/logger"
	"go.uber.org/zap"
	"reflect"
)

func (this *ZLog) Debug(msg string, field ...logger.Field) {
	if len(field) > 0 {
		for _, v := range field {
			this.engine.Debug(msg, this.fs(v)...)
		}
	} else {
		this.engine.Debug(msg)
	}

}

func (this *ZLog) Info(msg string, field ...logger.Field) {
	if len(field) > 0 {
		for _, v := range field {
			this.engine.Info(msg, this.fs(v)...)
		}
	} else {
		this.engine.Info(msg)
	}
}

func (this *ZLog) Error(msg string, field ...logger.Field) {
	if len(field) > 0 {
		for _, v := range field {
			this.engine.Error(msg, this.fs(v)...)
		}
	} else {
		this.engine.Error(msg)
	}
}

func (this *ZLog) Panic(msg string, field ...logger.Field) {
	if len(field) > 0 {
		for _, v := range field {
			this.engine.Panic(msg, this.fs(v)...)
		}
	} else {
		this.engine.Panic(msg)
	}
}

func (this *ZLog) Fatal(msg string, field ...logger.Field) {
	if len(field) > 0 {
		for _, v := range field {
			this.engine.Fatal(msg, this.fs(v)...)
		}
	} else {
		this.engine.Fatal(msg)
	}
}

func (this *ZLog) fs(data logger.Field) (fields []zap.Field) {
	fields = make([]zap.Field, 0)
	typ := reflect.TypeOf(data)
	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Struct {
		return fields
	}
	for i := 0; i < val.NumField(); i++ {
		if typ.Field(i).Tag.Get("json") == "ts" || typ.Field(i).Tag.Get("json") == "level" || typ.Field(i).Tag.Get("json") == "msg" {
			continue
		}

		if typ.Field(i).Tag.Get("json") != "" {
			if val.Field(i).Kind() == reflect.Map {
				for _, mk := range val.Field(i).MapKeys() {
					if mk.String() == "ts" || mk.String() == "level" || mk.String() == "msg" {
						continue
					}
					fields = append(fields, zap.Any(mk.String(), val.Field(i).MapIndex(mk).Elem().Interface()))
				}
			} else {
				fields = append(fields, zap.Any(typ.Field(i).Tag.Get("json"), val.Field(i).Interface()))
			}
		}
	}

	return fields
}
