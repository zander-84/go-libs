package zlog

type Writer func(p []byte) (n int, err error)

func (f Writer) Write(p []byte) (n int, err error) {
	return f(p)
}
