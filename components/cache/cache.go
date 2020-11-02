package cache

import (
	"time"
)

type Cache interface {
	Exists(key ...string) (bool, error)

	Get(key string, ptrValue interface{}) error

	GetFast(key string, ptrValue interface{}) (interface{}, error)

	Set(key string, value interface{}, expires time.Duration) error

	Delete(key ...string) error

	GetOrSet(key string, ptrValue interface{}, expires time.Duration, f func() (value interface{}, err error)) error

	GetOrSetFast(key string, ptrValue interface{}, expires time.Duration, f func() (value interface{}, err error)) (interface{}, error)

	FlushDB() error
}
