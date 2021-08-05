package goredis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

func (this *Rdb) TryLockWithTimeout(ctx context.Context, key string, identify string, duration time.Duration) (bool, error) {
	return this.engine.SetNX(ctx, key, identify, duration).Result()
}

func (this *Rdb) TryLockWithWaiting(ctx context.Context, key string, identify string, duration time.Duration, waitTime int) (bool, error) {
	for i := 0; i < waitTime; i++ {
		ok, err := this.engine.SetNX(ctx, key, identify, duration).Result()
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

func (this *Rdb) ReleaseLock(ctx context.Context, key string, identify string) error {
	data, err := this.GetString(ctx, key)
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		return err
	}

	if data == identify {
		if err := this.Dels(ctx, key); err != nil {
			if err == redis.Nil {
				return nil
			}
			return err
		}
	}
	return err
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
