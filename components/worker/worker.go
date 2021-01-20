package worker

type worker struct {
	jobChannel chan Job
	workerPool chan chan Job
	quit       chan bool
	running    bool // 0:等待状态  1：运行状态
}

func newWorker(workerPool chan chan Job) *worker {
	return &worker{
		jobChannel: make(chan Job),
		workerPool: workerPool,
		quit:       make(chan bool),
		running:    false,
	}
}

func (this *worker) start(i int) {
	go func() {
		for {
			this.workerPool <- this.jobChannel
			this.running = false
			select {
			case job := <-this.jobChannel:
				this.running = true
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
	this.running = false
	this.quit <- true
}
