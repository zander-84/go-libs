package helper

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Lock interface {
	Lock(ctx context.Context, key string, val interface{}, duration time.Duration) error
	UnLock(ctx context.Context, key string, val interface{}) error
	Data(ctx context.Context) (map[string]interface{}, error)
	Release()
	Println()
	Exit()
}

type memoryLock struct {
	sMap sync.Map
	exit chan struct{}
}
type singleLockData struct {
	expiredAt time.Time // 过期时间
	tag       int32
	val       interface{}
}

var errLock = errors.New("lock err")

func NewMemoryLock() Lock {
	l := new(memoryLock)
	l.exit = make(chan struct{}, 1)
	go func() {
		for {
			select {
			case <-l.exit:
				return
			case <-time.After(time.Second * 5):
				l.Release()
			}
		}
	}()
	return l
}

// Lock duration时间0 永不过期
func (l *memoryLock) Lock(ctx context.Context, key string, val interface{}, duration time.Duration) error {
	var expiredAt time.Time
	if duration < 1 {
		expiredAt = time.Time{}
	} else {
		expiredAt = time.Now().Add(duration)
	}

	idata, ok := l.sMap.LoadOrStore(key, &singleLockData{
		expiredAt: expiredAt,
		tag:       0,
		val:       val,
	})

	//不存在 第一次
	if !ok {
		return nil
	}

	// 非第一次
	user := idata.(*singleLockData)
	//并发下 同一个人只有第一次人才可以进去判断时间，第二次以上直接返回错误
	if atomic.CompareAndSwapInt32(&user.tag, 0, 1) {
		// 时间过期下才可以返回true
		if !user.expiredAt.IsZero() && time.Now().After(user.expiredAt) {
			user.expiredAt = expiredAt
			user.val = val
			user.tag = 0
			l.sMap.Store(key, user)
			return nil
		} else {
			user.tag = 0
			return errLock
		}
	} else {
		return errLock
	}
}

func (l *memoryLock) UnLock(ctx context.Context, key string, val interface{}) error {
	valInterface, ok := l.sMap.Load(key)
	if !ok {
		return nil
	}
	lockData, ok := valInterface.(*singleLockData)
	if ok && lockData.val == val {
		l.sMap.Delete(key)
	}
	return nil
}

func (l *memoryLock) Release() {
	l.sMap.Range(func(key, value interface{}) bool {
		user, ok := value.(*singleLockData)
		if !ok {
			return true
		}
		// 时间0  永不过期
		if user.expiredAt.IsZero() {
			return true
		}
		if time.Now().After(user.expiredAt) {
			l.sMap.Delete(key)
		}
		return true
	})
}
func (l *memoryLock) Exit() {
	l.exit <- struct{}{}
}
func (l *memoryLock) Println() {
	var cnt int64 = 0
	l.sMap.Range(func(key, value interface{}) bool {
		fmt.Printf("%v ==> %+v \n", key, value)
		atomic.AddInt64(&cnt, 1)
		return true
	})

	fmt.Println("len: ", cnt)
}
func (l *memoryLock) Data(ctx context.Context) (map[string]interface{}, error) {
	data := make(map[string]interface{}, 0)
	l.sMap.Range(func(key, value interface{}) bool {
		data[key.(string)] = value
		return true
	})
	return data, nil
}
