package logger

type Logger interface {
	Debug(from string, field ...Field)
	Info(from string, field ...Field)
	Error(from string, field ...Field)
	Panic(from string, field ...Field)
	Fatal(from string, field ...Field)
}

type Field struct {
	TraceId string
	Level   string
	Msg     string
	From    string
	Ts      string
}
