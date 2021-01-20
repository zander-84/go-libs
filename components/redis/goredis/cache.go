package goredis

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/zander-84/go-libs/components/cache"
	"github.com/zander-84/go-libs/components/errs"
	"reflect"
	"time"
)

var _ cache.Cache = (*Rdb)(nil)

func (this *Rdb) Get(ctx context.Context, key string, value interface{}) error {
	b, err := this.engine.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return errs.RecordNotFoundError
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(b, value)
}

func (this *Rdb) GetFast(ctx context.Context, key string, value interface{}) (interface{}, error) {
	err := this.Get(ctx, key, value)
	if err == redis.Nil {
		err = errs.RecordNotFoundError
	}
	return value, err
}

func (this *Rdb) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return this.engine.Set(ctx, key, b, ttl).Err()
}

func (this *Rdb) GetOrSet(ctx context.Context, key string, ptrValue interface{}, ttl time.Duration, f func() (interface{}, error)) error {
	if err := this.Get(ctx, key, ptrValue); err != nil {
		if err != redis.Nil {
			return errs.RecordNotFoundError
		}

		// todo once
		if fv, fe := this.singleflight.Do(key, f); fe != nil {
			return fe
		} else {
			if err := this.Set(ctx, key, fv, ttl); err != nil {
				return err
			}
			return this.Get(ctx, key, ptrValue)
		}
	}

	return nil
}
func (this *Rdb) GetOrSetFast(ctx context.Context, key string, ptrValue interface{}, ttl time.Duration, f func() (interface{}, error)) (interface{}, error) {
	err := this.Get(ctx, key, ptrValue)
	if err == nil {
		v := reflect.ValueOf(ptrValue)
		if v.Kind() == reflect.Ptr {
			return v.Elem().Interface(), err
		} else {
			return nil, err
		}
	}

	if err != redis.Nil {
		return nil, errs.RecordNotFoundError
	}

	if fv, fe := this.singleflight.Do(key, f); fe != nil {
		return nil, fe
	} else {
		if err := this.Set(ctx, key, fv, ttl); err != nil {
			return nil, err
		}
		return fv, fe
	}
}

func (this *Rdb) Delete(ctx context.Context, key ...string) error {
	return this.engine.Del(ctx, key...).Err()
}

func (this *Rdb) Exists(ctx context.Context, key ...string) (bool, error) {
	cmd := this.engine.Exists(ctx, key...)
	if cmd.Err() != nil {
		return false, cmd.Err()
	}

	if cmd.Val() == int64(len(key)) {
		return true, nil
	} else {
		return false, nil
	}
}

func (this *Rdb) FlushDB(ctx context.Context) error {
	return this.engine.FlushDB(ctx).Err()
}
