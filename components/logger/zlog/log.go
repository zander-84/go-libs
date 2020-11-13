package zlog

import (
	"github.com/zander-84/go-libs/components/logger"
	"go.uber.org/zap"
	"reflect"
)

func (this *ZLog) Debug(msg string, field ...logger.Field) {
	this.engine.Debug(msg, this.fs(field...)...)
}

func (this *ZLog) Info(msg string, field ...logger.Field) {
	this.engine.Info(msg, this.fs(field...)...)
}

func (this *ZLog) Error(msg string, field ...logger.Field) {
	this.engine.Error(msg, this.fs(field...)...)
}

func (this *ZLog) Panic(msg string, field ...logger.Field) {
	this.engine.Panic(msg, this.fs(field...)...)
}

func (this *ZLog) Fatal(msg string, field ...logger.Field) {
	this.engine.Fatal(msg, this.fs(field...)...)
}

func (this *ZLog) fs(data ...logger.Field) (fields []zap.Field) {

	if len(data) > 0 {
		fields = make([]zap.Field, 0)
		for _, v := range data {
			typ := reflect.TypeOf(v)
			val := reflect.ValueOf(v)
			if val.Kind() != reflect.Struct {
				return fields
			}
			for i := 0; i < val.NumField(); i++ {
				if typ.Field(i).Tag.Get("json") == "ts" || typ.Field(i).Tag.Get("json") == "level" || typ.Field(i).Tag.Get("json") == "msg" {
					continue
				}

				if typ.Field(i).Tag.Get("json") != "" {

					if val.Field(i).Kind() == reflect.String {
						fields = append(fields, zap.String(typ.Field(i).Tag.Get("json"), val.Field(i).String()))
					} else if val.Field(i).Kind() == reflect.Int {
						fields = append(fields, zap.Int(typ.Field(i).Tag.Get("json"), int(val.Field(i).Int())))
					} else if val.Field(i).Kind() == reflect.Int64 {
						fields = append(fields, zap.Int64(typ.Field(i).Tag.Get("json"), val.Field(i).Int()))
					} else if val.Field(i).Kind() == reflect.Float64 {
						fields = append(fields, zap.Float64(typ.Field(i).Tag.Get("json"), val.Field(i).Float()))
					}
				}

			}

			//if v.From != "" {
			//	fields = append(fields, zap.String("from", v.From))
			//}
			//if v.TraceId != "" {
			//	fields = append(fields, zap.String("trace_id", v.TraceId))
			//}
		}
	}
	return fields
}
