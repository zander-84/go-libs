package zlog

import (
	"github.com/zander-84/go-libs/components/logger"
	"go.uber.org/zap"
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
			if v.From != "" {
				fields = append(fields, zap.String("from", v.From))
			}
			if v.TraceId != "" {
				fields = append(fields, zap.String("trace_id", v.TraceId))
			}
		}
	}
	return fields
}
