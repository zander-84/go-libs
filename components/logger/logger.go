package logger

type Logger interface {
	Debug(from string, files ...Filed)
	Info(from string, files ...Filed)
	Error(from string, files ...Filed)
	Panic(from string, files ...Filed)
	Fatal(from string, files ...Filed)
}

type Filed struct {
	TraceId string
	Level   string
	Msg     string
	From    string
	Ts      string
}
