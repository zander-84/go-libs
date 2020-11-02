package helper

import "sync"

type GoGroup struct {
	wg  sync.WaitGroup
	err error
}

func NewGoGroup() *GoGroup {
	return new(GoGroup)
}

func (this *GoGroup) AsyncFunc(actions ...func()) {
	if len(actions) > 0 {
		for _, action := range actions {
			this.wg.Add(1)
			go func(a func()) {
				defer this.wg.Done()
				a()
			}(action)
		}
		this.wg.Wait()
	}
}

func (this *GoGroup) SyncFunc(actions ...func()) {
	if len(actions) > 0 {
		for _, action := range actions {
			action()
		}
	}
}
