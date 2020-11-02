package robfig

import (
	"github.com/robfig/cron/v3"
	cron2 "github.com/zander-84/go-libs/components/cron"
	"github.com/zander-84/go-libs/components/errs"
)

var _ cron2.Crontab = (*Robfig)(nil)

func (this *Robfig) AddJob(job *cron2.Job) error {
	this.lock.Lock()
	defer this.lock.Unlock()
	if _, ok := this.jobs[job.ID]; ok {
		return errs.IDExistError
	}
	// 添加job
	id, err := this.engine.AddJob(job.Spec, job.Cmd)
	if err == nil {
		job.Obj = this.engine.Entry(id)
		this.jobs[job.ID] = job
	}
	return err
}

// 移除
func (this *Robfig) RemoveJob(id string) error {
	this.lock.Lock()
	defer this.lock.Unlock()
	if job, ok := this.jobs[id]; ok {
		entry := job.Obj.(cron.Entry)
		this.engine.Remove(entry.ID)
		delete(this.jobs, id)
	}

	return this.err
}

func (this *Robfig) StatusJobs() map[string]*cron2.Job {
	this.lock.Lock()
	defer this.lock.Unlock()

	for _, job := range this.jobs {
		entry := job.Obj.(cron.Entry)
		job.Obj = this.engine.Entry(entry.ID)
	}
	return this.jobs
}

func (this *Robfig) StartJobs() error {
	this.engine.Start()
	return this.err
}

func (this *Robfig) StopJobs() error {
	this.engine.Stop()
	return this.err

}

func (this *Robfig) RestartJob(id string) error {
	if job, ok := this.jobs["id"]; ok {
		if err := this.RemoveJob(id); err != nil {
			return err
		}
		return this.AddJob(job)
	}
	return errs.IDError
}
