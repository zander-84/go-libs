package goredis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/zander-84/go-libs/components/cache"
	"github.com/zander-84/go-libs/think"
	"reflect"
	"runtime"
	"time"
)

var _ cache.Cache = (*Rdb)(nil)

func (this *Rdb) MGet(ctx context.Context, keys []string, ptrSliceData interface{}) (lostKeys []string, err error) {
	if reflect.ValueOf(ptrSliceData).Elem().Type().Kind() != reflect.Slice {
		return nil, errors.New("data  must be slice ptr")
	}

	res, err := this.engine.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}
	lostKeys = make([]string, 0)
	var keyCnt = len(keys)
	reflectValue := reflect.ValueOf(ptrSliceData).Elem()
	for key, v := range res {
		if v != nil {
			tmp := reflect.New(reflectValue.Type().Elem())
			tmp2, ok := v.(string)
			if !ok {
				return nil, errors.New("unmarshal error")
			}
			if err := json.Unmarshal([]byte(tmp2), tmp.Interface()); err != nil {
				return nil, err
			}
			reflectValue.Set(reflect.Append(reflectValue, tmp.Elem()))
		} else {
			if key >= keyCnt {
				return nil, fmt.Errorf("index out of range [%d] with length %d", key, keyCnt)
			}
			lostKeys = append(lostKeys, keys[key])
			reflectValue.Set(reflect.Append(reflectValue, reflect.Zero(reflectValue.Type().Elem())))
		}
	}

	return lostKeys, nil
}

func (this *Rdb) MustMGetOrSet(ctx context.Context, rawKeys []string, redisKeys []string, ptrSliceData interface{}, ttl time.Duration, f func(id string) (interface{}, error)) (lostKey string, err error) {
	defer func() {
		if rerr := recover(); rerr != nil {
			buf := make([]byte, 64<<10)
			n := runtime.Stack(buf, false)
			buf = buf[:n]
			err = errors.New(string(buf))
		}
	}()

	lostKeys, err := this.MGet(ctx, redisKeys, ptrSliceData)
	if err != nil {
		return "", err
	}
	if len(lostKeys) < 1 {
		return "", nil
	}

	reflectValue := reflect.ValueOf(ptrSliceData).Elem()

	for k, v := range redisKeys {
		for _, vv := range lostKeys {
			if v == vv {
				if fv, fe := f(rawKeys[k]); fe != nil {
					return rawKeys[k], fe
				} else {
					if err := this.SetNX(ctx, vv, fv, ttl); err != nil {
						return rawKeys[k], err
					}

					tmp := reflect.New(reflectValue.Type().Elem())
					if err := this.Get(ctx, vv, tmp.Interface()); err != nil {
						return rawKeys[k], err
					}
					reflectValue.Index(k).Set(tmp.Elem())
				}
			}
		}
	}
	return "", nil
}

func (this *Rdb) Get(ctx context.Context, key string, toPtr interface{}) (err error) {
	defer func() {
		if rerr := recover(); rerr != nil {
			buf := make([]byte, 64<<10)
			n := runtime.Stack(buf, false)
			buf = buf[:n]
			err = errors.New(string(buf))
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
