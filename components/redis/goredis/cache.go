package goredis

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/zander-84/go-libs/think"
	"reflect"
	"runtime"
	"time"
)

func (this *Rdb) Get(ctx context.Context, key string, toPtr interface{}) (err error) {
	defer func() {
		if rerr := recover(); rerr != nil {
			buf := make([]byte, 64<<10)
			n := runtime.Stack(buf, false)
			buf = buf[:n]
			err = think.ErrSystemSpace(errors.New(string(buf)))
		}
	}()

	data, err := this.getSingleFlight.Do(key, func() (interface{}, error) {
		b, err := this.engine.Get(ctx, key).Bytes()
		if err == redis.Nil {
			return nil, think.ErrInstanceRecordNotFound
		}
		if err != nil {
			b, err = this.engine.Get(ctx, key).Bytes()
			if err == redis.Nil {
				return nil, think.ErrInstanceRecordNotFound
			}
			if err != nil {
				return nil, err
			}
		}
		err = json.Unmarshal(b, toPtr)
		if err != nil {
			return nil, err
		}
		return toPtr, nil
	})

	if err != nil {
		return err
	}

	if data != nil && toPtr != data {
		if reflect.ValueOf(data).Type().Kind() == reflect.Ptr {
			reflect.ValueOf(toPtr).Elem().Set(reflect.ValueOf(data).Elem())
		} else {
			reflect.ValueOf(toPtr).Elem().Set(reflect.ValueOf(data))
		}
	}

	return err
}

func (this *Rdb) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return this.engine.Set(ctx, key, b, ttl).Err()
}

func (this *Rdb) SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return this.engine.SetNX(ctx, key, b, ttl).Err()
}

func (this *Rdb) GetOrSet(ctx context.Context, key string, ptrValue interface{}, ttl time.Duration, f func() (interface{}, error)) error {
	if err := this.Get(ctx, key, ptrValue); err != nil {
		if !think.IsErrNotFound(err) {
			return err
		}

		if fv, fe := this.funSingleFlight.Do(key, f); fe != nil {
			return fe
		} else {
			// 不允许覆盖
			if err := this.SetNX(ctx, key, fv, ttl); err != nil {
				return err
			}
			return this.Get(ctx, key, ptrValue)
		}
	}

	return nil
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
