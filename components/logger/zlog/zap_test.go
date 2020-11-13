package zlog

import (
	"encoding/json"
	"fmt"
	"github.com/zander-84/go-libs/components/logger"
	"io"
	"testing"
)

// go test -v  -run TestZap
func TestZap(t *testing.T) {
	zlog := NewZapLog(Conf{
		Level:     "info",
		Name:      "test",
		TimeZone:  "",
		AddCaller: false,
		ConsoleHook: struct {
			Enable bool
		}{
			Enable: true,
		},
		FileHook: struct {
			Enable     bool
			Path       string
			MaxAge     int
			MaxBackups int
			MaxSize    int
		}{
			Enable:     true,
			Path:       "./logs",
			MaxAge:     10,
			MaxBackups: 10,
			MaxSize:    50,
		},
	}, []io.Writer{Writer(func(p []byte) (n int, err error) {
		fmt.Println("aaa: ", string(p))

		var data logger.Field
		if err := json.Unmarshal(p, &data); err != nil {
			return len(p), err
		}
		return 0, nil
	})})
	if err := zlog.Start(); err != nil {
		t.Fatal("start log err: ", err.Error())
	}

	zlog.Debug("debug", logger.Field{From: "aaa"})
	zlog.Info("Info", logger.Field{From: "bbb", TraceId: "ssssss"})
	zlog.Error("Error", logger.Field{From: "aaa"}, logger.Field{From: "bbb", TraceId: "ssssss"})
}
