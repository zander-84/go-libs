package zlog

import (
	"github.com/zander-84/go-libs/components/logger"
	"go.uber.org/zap"
)

func (this *ZLog) Debug(msg string, files ...logger.Filed) {
	this.engine.Debug(msg, this.fs(files...)...)
}

func (this *ZLog) Info(msg string, files ...logger.Filed) {
	this.engine.Info(msg, this.fs(files...)...)
}

func (this *ZLog) Error(msg string, files ...logger.Filed) {
	this.engine.Error(msg, this.fs(files...)...)
}

func (this *ZLog) Panic(msg string, files ...logger.Filed) {
	this.engine.Panic(msg, this.fs(files...)...)
}

func (this *ZLog) Fatal(msg string, files ...logger.Filed) {
	this.engine.Fatal(msg, this.fs(files...)...)
}

func (this *ZLog) fs(data ...logger.Filed) (fields []zap.Field) {
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