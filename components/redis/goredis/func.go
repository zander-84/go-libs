package goredis

import "time"

func (this *Rdb) TryLockWithTimeout(identify string, duration time.Duration) (bool, error) {
	return this.engine.SetNX(this.context, identify, true, duration).Result()
}

func (this *Rdb) TryLockWithWaiting(identify string, duration time.Duration, waitTime int) (bool, error) {
	for i := 0; i < waitTime; i++ {
		ok, err := this.engine.SetNX(this.context, identify, true, duration).Result()
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

func (this *Rdb) SetString(key string, str string, expires time.Duration) (err error) {
	return this.engine.Set(this.context, key, str, expires).Err()
}

func (this *Rdb) GetString(key string) (string, error) {
	return this.engine.Get(this.context, key).Result()
}

func (this *Rdb) GetBytes(key string) ([]byte, error) {
	return this.engine.Get(this.context, key).Bytes()
}

func (this *Rdb) Dels(keys ...string) (err error) {
	return this.engine.Del(this.context, keys...).Err()
}
