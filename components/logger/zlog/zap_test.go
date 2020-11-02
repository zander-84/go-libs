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
		fmt.Println(string(p))

		var data logger.Filed
		if err := json.Unmarshal(p, &data); err != nil {
			return len(p), err
		}
		return 0, nil
	})})
	if err := zlog.Start(); err != nil {
		t.Fatal("start log err: ", err.Error())
	}

	zlog.Debug("debug", logger.Filed{From: "aaa"})
	zlog.Info("Info", logger.Filed{From: "bbb", TraceId: "ssssss"})
	zlog.Error("Error", logger.Filed{From: "aaa"}, logger.Filed{From: "bbb", TraceId: "ssssss"})
}
