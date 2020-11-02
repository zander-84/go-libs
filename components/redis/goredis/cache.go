package goredis

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/zander-84/go-libs/components/cache"
	"reflect"
	"time"
)

var _ cache.Cache = (*Rdb)(nil)

func (this *Rdb) Get(key string, value interface{}) error {
	b, err := this.engine.Get(this.context, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(b, value)
}

func (this *Rdb) GetFast(key string, value interface{}) (interface{}, error) {
	err := this.Get(key, value)
	return value, err
}

func (this *Rdb) Set(key string, value interface{}, ttl time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return this.engine.Set(this.context, key, b, ttl).Err()
}

func (this *Rdb) GetOrSet(key string, ptrValue interface{}, ttl time.Duration, f func() (interface{}, error)) error {
	if err := this.Get(key, ptrValue); err != nil {
		if err != redis.Nil {
			return err
		}

		// todo once
		if fv, fe := f(); fe != nil {
			return fe
		} else {
			if err := this.Set(key, fv, ttl); err != nil {
				return err
			}
			return this.Get(key, ptrValue)
		}
	}

	return nil
}
func (this *Rdb) GetOrSetFast(key string, ptrValue interface{}, ttl time.Duration, f func() (interface{}, error)) (interface{}, error) {
	err := this.Get(key, ptrValue)
	if err == nil {
		v := reflect.ValueOf(ptrValue)
		if v.Kind() == reflect.Ptr {
			return v.Elem().Interface(), err
		} else {
			return nil, err
		}
	}

	if err != redis.Nil {
		return nil, err
	}

	if fv, fe := f(); fe != nil {
		return nil, fe
	} else {
		if err := this.Set(key, fv, ttl); err != nil {
			return nil, err
		}
		return fv, fe
	}
}

func (this *Rdb) Delete(key ...string) error {
	return this.engine.Del(this.context, key...).Err()
}

func (this *Rdb) Exists(key ...string) (bool, error) {
	cmd := this.engine.Exists(this.context, key...)
	if cmd.Err() != nil {
		return false, cmd.Err()
	}

	if cmd.Val() == int64(len(key)) {
		return true, nil
	} else {
		return false, nil
	}
}

func (this *Rdb) FlushDB() error {
	return this.engine.FlushDB(this.context).Err()
}
