package logger

type Logger interface {
	Debug(msg string, field ...Field)
	Info(msg string, field ...Field)
	Error(msg string, field ...Field)
	Panic(msg string, field ...Field)
	Fatal(msg string, field ...Field)
}

type Field struct {
	TraceId   string                 `json:"trace_id"`
	Level     string                 `json:"level"`
	Msg       string                 `json:"msg"`
	From      string                 `json:"from"`
	Ts        string                 `json:"ts"`
	Belong    string                 `json:"belong"` //
	SubFields map[string]interface{} `json:"sub_fields"`
}

type Instance struct {
	name string
	log  Logger
}
