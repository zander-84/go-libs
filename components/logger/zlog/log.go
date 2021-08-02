package zlog

import (
	"github.com/zander-84/go-libs/components/logger"
	"go.uber.org/zap"
)

func (this *ZLog) Debug(msg string, fds ...logger.Field) {
	this.engine.Debug(msg, this.fs(fds)...)
}

func (this *ZLog) Info(msg string, fds ...logger.Field) {
	this.engine.Info(msg, this.fs(fds)...)
}

func (this *ZLog) Error(msg string, fds ...logger.Field) {
	this.engine.Error(msg, this.fs(fds)...)
}

func (this *ZLog) Panic(msg string, fds ...logger.Field) {
	this.engine.Panic(msg, this.fs(fds)...)
}

func (this *ZLog) Fatal(msg string, fds ...logger.Field) {
	this.engine.Fatal(msg, this.fs(fds)...)
}

func (this *ZLog) fs(fds []logger.Field) []zap.Field {
	if len(fds) < 1 {
		return nil
	}

	var fields = make([]zap.Field, 0)
	for _, v := range fds {
		for kk, vv := range v {
			fields = append(fields, zap.String(kk, vv))
		}
	}
	return fields
}
