package memory

import (
	"github.com/zander-84/go-libs/components/cache"
	"github.com/zander-84/go-libs/components/errs"
	"reflect"
	"time"
)

var _ cache.Cache = (*Memory)(nil)

func (this *Memory) Get(key string, ptrValue interface{}) (err error) {
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

func (this *Memory) GetFast(key string, ptr interface{}) (interface{}, error) {
	value, ok := this.engine.Get(key)
	if !ok {
		return nil, errs.RecordNotFoundError
	}
	return value, nil
}

func (this *Memory) Set(key string, value interface{}, ttl time.Duration) (err error) {
	this.engine.Set(key, value, ttl)
	return this.err
}

func (this *Memory) GetOrSet(key string, ptrValue interface{}, ttl time.Duration, f func() (value interface{}, err error)) (err error) {
	defer func() {
		if rerr := recover(); rerr != nil {
			err = errs.TypeError
		}
	}()
	err = this.Get(key, ptrValue)
	if err == nil {
		return err
	}
	if err != errs.RecordNotFoundError {
		return err
	}

	fv, fe := f()
	if fe != nil {
		err = fe
		return err
	}
	err = this.Set(key, fv, ttl)
	if err == nil {
		err = this.Get(key, ptrValue)
	}
	return err
}
func (this *Memory) GetOrSetFast(key string, ptrValue interface{}, ttl time.Duration, f func() (value interface{}, err error)) (value interface{}, err error) {
	defer func() {
		if rerr := recover(); rerr != nil {
			err = errs.TypeError
		}
	}()

	value, err = this.GetFast(key, ptrValue)
	if err == nil {
		return value, err
	}
	if err != errs.RecordNotFoundError {
		return value, err
	}

	fv, fe := f()
	if fe != nil {
		return value, fe
	}
	err = this.Set(key, fv, ttl)
	if err != nil {
		return fv, err
	}

	return fv, nil
}

func (this *Memory) Delete(keys ...string) error {
	for _, key := range keys {
		this.engine.Delete(key)
	}
	return nil
}

func (this *Memory) Exists(keys ...string) (bool, error) {
	for _, key := range keys {
		if _, ok := this.engine.Get(key); !ok {
			return false, nil
		}
	}
	return true, nil
}

func (this *Memory) FlushDB() error {
	this.engine.Flush()
	return nil
}
