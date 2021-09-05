package worker

type worker struct {
	jobChannel chan Job
	workerPool chan chan Job
	quit       chan bool
}

func newWorker(workerPool chan chan Job) *worker {
	return &worker{
		jobChannel: make(chan Job),
		workerPool: workerPool,
		quit:       make(chan bool),
	}
}

func (this *worker) start(i int) {
	go func() {
		for {
			this.workerPool <- this.jobChannel
			select {
			case job := <-this.jobChannel:
				if job != nil {
					func() {
						defer func() {
							if err := recover(); err != nil {
							}
						}()
						_ = job.Run()
					}()
				}
			case <-this.quit:
				return
			}
		}
	}()

}

func (this *worker) stop() {
	this.quit <- true
}
