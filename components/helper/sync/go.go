package sync

import (
	"context"
	"time"
)

func Go(fn func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				return
			}
		}()
		fn()
	}()
}

func RunWithTimeout(fn func() error, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer func() {
		cancel()
	}()
	fin := make(chan error, 1)
	go func() {
		fin <- fn()
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case e := <-fin:
		return e
	}
}
