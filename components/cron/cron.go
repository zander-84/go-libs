package cron

/* Spec
Field name   | Mandatory? | Allowed values  | Allowed special characters
----------   | ---------- | --------------  | --------------------------
Seconds      | Yes        | 0-59            | * / , -
Minutes      | Yes        | 0-59            | * / , -
Hours        | Yes        | 0-23            | * / , -
Day of month | Yes        | 1-31            | * / , - ?
Month        | Yes        | 1-12 or JAN-DEC | * / , -
Day of week  | Yes        | 0-6 or SUN-SAT  | * / , - ?
*/
type Job struct {
	ID   string
	Desc string
	Spec string
	Cmd  Cmd
	Obj  interface{}
}

type Crontab interface {
	AddJob(jobInfo *Job) error

	RemoveJob(id string) error

	StatusJobs() map[string]*Job

	StartJobs() error

	StopJobs() error

	RestartJob(id string) error
}

type Cmd interface {
	Run()
}

type CmdFunc func()

func (this CmdFunc) Run() {
	this()
}
