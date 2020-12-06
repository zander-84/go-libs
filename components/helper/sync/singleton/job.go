package singleton

import "context"

type Job struct {
	Ctx     context.Context
	Name    string
	Handler func(j *Job) error
	Input   interface{}
	Output  interface{}
	Error   error
}
