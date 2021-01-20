package goredis

import (
	"context"
	"time"
)

func (this *Rdb) TryLockWithTimeout(ctx context.Context, identify string, duration time.Duration) (bool, error) {
	return this.engine.SetNX(ctx, identify, true, duration).Result()
}

func (this *Rdb) TryLockWithWaiting(ctx context.Context, identify string, duration time.Duration, waitTime int) (bool, error) {
	for i := 0; i < waitTime; i++ {
		ok, err := this.engine.SetNX(ctx, identify, true, duration).Result()
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
		time.Sleep(time.Second)
	}
	return false, nil
}

func (this *Rdb) SetString(ctx context.Context, key string, str string, expires time.Duration) (err error) {
	return this.engine.Set(ctx, key, str, expires).Err()
}

func (this *Rdb) GetString(ctx context.Context, key string) (string, error) {
	return this.engine.Get(ctx, key).Result()
}

func (this *Rdb) GetBytes(ctx context.Context, key string) ([]byte, error) {
	return this.engine.Get(ctx, key).Bytes()
}

func (this *Rdb) Dels(ctx context.Context, keys ...string) (err error) {
	return this.engine.Del(ctx, keys...).Err()
}
