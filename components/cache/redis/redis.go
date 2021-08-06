package redis

import (
	"context"
	"github.com/zander-84/go-libs/components/cache"
	"github.com/zander-84/go-libs/components/redis"
	"time"
)

var _ cache.Cache = (*Redis)(nil)

type Redis struct {
	r *redis.Rdb
}

func New(r *redis.Rdb) *Redis {
	return &Redis{
		r: r,
	}
}

func (this *Redis) Get(ctx context.Context, key string, recPtr interface{}) error {
	return this.r.Get(ctx, key, recPtr)
}

func (this *Redis) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) (err error) {
	return this.r.Set(ctx, key, value, ttl)
}

func (this *Redis) GetOrSet(ctx context.Context, key string, recPtr interface{}, ttl time.Duration, f func() (value interface{}, err error)) (err error) {
	return this.r.GetOrSet(ctx, key, recPtr, ttl, f)
}

func (this *Redis) Delete(ctx context.Context, keys ...string) error {
	return this.r.Delete(ctx, keys...)
}

func (this *Redis) Exists(ctx context.Context, keys ...string) (bool, error) {
	return this.r.Exists(ctx, keys...)
}

func (this *Redis) FlushDB(ctx context.Context) error {
	return this.r.FlushDB(ctx)
}
