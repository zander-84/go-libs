package memory

import (
	"context"
	"github.com/zander-84/go-libs/components/cache"
	"github.com/zander-84/go-libs/components/errs"
	"reflect"
	"time"
)

var _ cache.Cache = (*Memory)(nil)

func (this *Memory) Get(ctx context.Context, key string, ptrValue interface{}) (err error) {
	defer func() {
		if rerr := recover(); rerr != nil {
			err = errs.TypeError
		}
	}()

	if found, ok := this.engine.Get(key); ok {
		if reflect.ValueOf(found).Type().Kind() == reflect.Ptr {
			reflect.ValueOf(ptrValue).Elem().Set(reflect.ValueOf(found).Elem())
		} else {
			reflect.ValueOf(ptrValue).Elem().Set(reflect.ValueOf(found))
		}
		return this.err
	}
	return errs.RecordNotFoundError
}

func (this *Memory) GetFast(ctx context.Context, key string, ptr interface{}) (interface{}, error) {
	value, ok := this.engine.Get(key)
	if !ok {
		return nil, errs.RecordNotFoundError
	}
	return value, nil
}

func (this *Memory) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) (err error) {
	if ttl == 0 {
		ttl = -1
	}
	this.engine.Set(key, value, ttl)
	return this.err
}

func (this *Memory) GetOrSet(ctx context.Context, key string, ptrValue interface{}, ttl time.Duration, f func() (value interface{}, err error)) (err error) {
	defer func() {
		if rerr := recover(); rerr != nil {
			err = errs.TypeError
		}
	}()
	err = this.Get(ctx, key, ptrValue)
	if err == nil {
		return err
	}
	if err != errs.RecordNotFoundError {
		return err
	}
	fv, fe := this.singleflight.Do(key, f)
	//fv, fe := f()
	if fe != nil {
		err = fe
		return err
	}
	err = this.Set(ctx, key, fv, ttl)
	if err == nil {
		err = this.Get(ctx, key, ptrValue)
	}
	return err
}
func (this *Memory) GetOrSetFast(ctx context.Context, key string, ptrValue interface{}, ttl time.Duration, f func() (value interface{}, err error)) (value interface{}, err error) {
	defer func() {
		if rerr := recover(); rerr != nil {
			err = errs.TypeError
		}
	}()

	value, err = this.GetFast(ctx, key, ptrValue)
	if err == nil {
		return value, err
	}
	if err != errs.RecordNotFoundError {
		return value, err
	}

	//fv, fe := f()
	fv, fe := this.singleflight.Do(key, f)
	if fe != nil {
		return value, fe
	}
	err = this.Set(ctx, key, fv, ttl)
	if err != nil {
		return fv, err
	}

	return fv, nil
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
