package memory

import (
	"context"
	"errors"
	"github.com/zander-84/go-libs/components/cache"
	"github.com/zander-84/go-libs/think"
	"reflect"
	"runtime"
	"time"
)

var _ cache.Cache = (*Memory)(nil)

func (this *Memory) MustMGetOrSet(ctx context.Context, rawKeys []string, redisKeys []string, recPtr interface{}, expires time.Duration, f func(key string) (value interface{}, err error)) (string, error) {
	return "", errors.New("todo")
}

func (this *Memory) MGet(ctx context.Context, keys []string, ptrSliceData interface{}) (lostKeys []string, err error) {
	return nil, errors.New("todo")
}

func (this *Memory) unmarshal(form interface{}, toPtr interface{}) (err error) {
	defer func() {
		if rerr := recover(); rerr != nil {
			buf := make([]byte, 64<<10)
			n := runtime.Stack(buf, false)
			buf = buf[:n]
			err = think.ErrSystemSpace(errors.New(string(buf)))
		}
	}()

	if reflect.ValueOf(form).Type().Kind() == reflect.Ptr {
		reflect.ValueOf(toPtr).Elem().Set(reflect.ValueOf(form).Elem())
	} else {
		reflect.ValueOf(toPtr).Elem().Set(reflect.ValueOf(form))
	}

	return err
}

func (this *Memory) Get(ctx context.Context, key string, recPtr interface{}) error {
	value, ok := this.engine.Get(key)
	if !ok {
		return think.ErrInstanceRecordNotFound
	}

	return this.unmarshal(value, recPtr)
}

func (this *Memory) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) (err error) {
	if ttl == 0 {
		ttl = -1
	}
	this.engine.Set(key, value, ttl)
	return this.err
}

func (this *Memory) SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (err error) {
	return errors.New("todo")
}
func (this *Memory) GetOrSet(ctx context.Context, key string, recPtr interface{}, ttl time.Duration, f func() (value interface{}, err error)) (err error) {
	defer func() {
		if rerr := recover(); rerr != nil {
			err = think.ErrType(errors.New("类型错误"))
		}
	}()
	err = this.Get(ctx, key, recPtr)
	if err == nil {
		return err
	}
	if !think.IsErrNotFound(err) {
		return err
	}

	fv, fe := this.singleflight.Do(key, f)
	if fe != nil {
		err = fe
		return err
	}
	err = this.Set(ctx, key, fv, ttl)
	if err == nil {
		err = this.Get(ctx, key, recPtr)
	}
	return err
}

func (this *Memory) Delete(ctx context.Context, keys ...string) error {
	for _, key := range keys {
		this.engine.Delete(key)
	}
	return nil
}

func (this *Memory) Exists(ctx context.Context, keys ...string) (bool, error) {
	for _, key := range keys {
		if _, ok := this.engine.Get(key); !ok {
			return false, nil
		}
	}
	return true, nil
}

func (this *Memory) FlushDB(ctx context.Context) error {
	this.engine.Flush()
	return nil
}
