package pipeline

import "context"

type Event struct {
	Ctx     context.Context
	Name    string
	Handler func(e *Event) error
	Input   interface{}
	Output  interface{}
	Error   error
}
