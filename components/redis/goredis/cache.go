package goredis

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/zander-84/go-libs/components/cache"
	"github.com/zander-84/go-libs/think"
	"time"
)

var _ cache.Cache = (*Rdb)(nil)
var nilByte, _ = json.Marshal(nil)
var NilStr = string(nilByte)

func (this *Rdb) Get(ctx context.Context, key string, toPtr interface{}) error {
	b, err := this.engine.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return think.ErrInstanceRecordNotFound
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(b, toPtr)
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
func (this *Rdb) GetNoNil(ctx context.Context, key string, toPtr interface{}) error {
	cmd := this.engine.Get(ctx, key)
	err := cmd.Err()
	if err == redis.Nil {
		return think.ErrInstanceRecordNotFound
	}
	if err != nil {
		return err
	}

	if cmd.Val() == NilStr {
		return cache.ErrDataNil
	}

	b, _ := cmd.Bytes()
	return json.Unmarshal(b, toPtr)
}

// GetOrSetConsistent 后台编辑操作必须更val为nil
func (this *Rdb) GetOrSetConsistent(ctx context.Context, key string, ptrValue interface{}, ttl time.Duration, f func() (interface{}, error)) error {
	return this.getOrSetConsistent(ctx, key, ptrValue, ttl, 0, f)
}

func (this *Rdb) getOrSetConsistent(ctx context.Context, key string, ptrValue interface{}, ttl time.Duration, maxtry int, f func() (interface{}, error)) error {
	if maxtry > 10 {
		return errors.New("数据异常")
	}
	maxtry = maxtry + 1
	if err := this.GetNoNil(ctx, key, ptrValue); err != nil {
		if err == cache.ErrDataNil {
			if err := this.Delete(ctx, key); err != nil {
				return err
			}

			return this.getOrSetConsistent(ctx, key, ptrValue, ttl, maxtry, f)
		}

		if !think.IsErrNotFound(err) {
			return err
		}

		fv, fe := this.singleflight.Do(key, f)
		if fe != nil {
			return fe
		}
		if fv == nil {
			return errors.New("data can't nil")
		}

		if err := this.SetNX(ctx, key, fv, ttl); err != nil {
			return err
		}
		err = this.GetNoNil(ctx, key, ptrValue)
		if err == cache.ErrDataNil {
			if err := this.Delete(ctx, key); err != nil {
				return err
			}

			return this.getOrSetConsistent(ctx, key, ptrValue, ttl, maxtry, f)
		}
		return nil
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
