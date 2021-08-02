package logger

import "strconv"

type Field map[string]string

type Logger interface {
	Debug(msg string, fds ...Field)
	Info(msg string, fds ...Field)
	Error(msg string, fds ...Field)
	Panic(msg string, fds ...Field)
	Fatal(msg string, fds ...Field)
}

//type Field struct {
//	TraceId   string                 `json:"trace_id"`
//	Level     string                 `json:"level"`
//	Msg       string                 `json:"msg"`
//	From      string                 `json:"from"` // 错误来源
//	Ts        string                 `json:"ts"`
//	Belong    string                 `json:"belong"` // 接口分类（项目来源）
//	SubFields map[string]interface{} `json:"sub_fields"`
//}

type Instance struct {
	name string
	log  Logger
}

func (f Field) GetStr(key string) string {
	if v, ok := f[key]; ok {
		return v
	} else {
		return ""
	}
}

func (f Field) GetInt64(key string) int64 {
	d := f.GetStr(key)
	int10, err := strconv.ParseInt(d, 10, 32)
	if err != nil {
		return 0
	}
	return int10
}

func (f Field) GetInt32(key string) int32 {
	d := f.GetStr(key)
	int10, err := strconv.ParseInt(d, 10, 32)
	if err != nil {
		return 0
	}
	return int32(int10)
}
