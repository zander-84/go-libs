package cache

import (
	"context"
	"time"
)

type Cache interface {
	Exists(ctx context.Context, key ...string) (bool, error)

	Get(ctx context.Context, key string, recPtr interface{}) error

	Set(ctx context.Context, key string, value interface{}, expires time.Duration) error
	SetNX(ctx context.Context, key string, value interface{}, expires time.Duration) error

	Delete(ctx context.Context, key ...string) error

	GetOrSet(ctx context.Context, key string, recPtr interface{}, expires time.Duration, f func() (value interface{}, err error)) error
	GetOrSetNX(ctx context.Context, key string, recPtr interface{}, expires time.Duration, f func() (value interface{}, err error)) error

	FlushDB(ctx context.Context) error
}
