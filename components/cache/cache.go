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
	MGet(ctx context.Context, keys []string, ptrSliceData interface{}) (lostKeys []string, err error)
	MustMGetOrSet(ctx context.Context, rawKeys []string, redisKeys []string, recPtr interface{}, expires time.Duration, f func(id string) (value interface{}, err error)) error
	FlushDB(ctx context.Context) error
}

// MsCache master slaves cache
type MsCache interface {
	Cache
	GetFromMaster(ctx context.Context, key string, recPtr interface{}) error
	GetFromSlave(ctx context.Context, key string, recPtr interface{}) error
}
