package cache

import (
	"context"
	"time"
)

type Cache interface {
	Exists(ctx context.Context, key ...string) (bool, error)

	Get(ctx context.Context, key string, ptrValue interface{}) error

	GetFast(ctx context.Context, key string, ptrValue interface{}) (interface{}, error)

	Set(ctx context.Context, key string, value interface{}, expires time.Duration) error

	Delete(ctx context.Context, key ...string) error

	GetOrSet(ctx context.Context, key string, ptrValue interface{}, expires time.Duration, f func() (value interface{}, err error)) error

	GetOrSetFast(ctx context.Context, key string, ptrValue interface{}, expires time.Duration, f func() (value interface{}, err error)) (interface{}, error)

	FlushDB(ctx context.Context) error
}
