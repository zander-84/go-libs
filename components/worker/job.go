package worker

type Job interface {
	Run() error
}

type JobFunc func() error

func (j JobFunc) Run() error {
	return j()
}
