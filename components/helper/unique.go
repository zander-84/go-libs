package helper

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type Unique struct {
	incrementID     uint64 // 自增ID
	incrementTimeID uint64 // 自增ID
	joinSymbol      string // 连接符
	lock            sync.Mutex
	prefixVal       string //前缀加机器码
	lastTime        int64  // 上次更新时间
}

// NewUnique 适合吞吐在 100w之下 qps服务；  重复不可避免，还得通过其它持久化方式优化
func NewUnique(prefix string, machine string, joinSymbol string) *Unique {
	u := new(Unique)
	u.joinSymbol = joinSymbol
	if prefix != "" {
		u.prefixVal = u.prefixVal + prefix + u.joinSymbol
	}
	if machine != "" {
		u.prefixVal = u.prefixVal + machine + u.joinSymbol
	}

	u.incrementTimeID = uint64(rand.Intn(100))
	rand.Seed(time.Now().UnixNano())
	return u
}

func (u *Unique) IDWithTag(tag string) string {
	u.lock.Lock()
	currentTime := u.now()
	d := atomic.AddUint64(&u.incrementID, 1)
	dt := atomic.LoadUint64(&u.incrementTimeID)
	if atomic.LoadInt64(&u.lastTime) == 0 {
		atomic.StoreInt64(&u.lastTime, currentTime.Unix())
	} else {
		if d > 999999 {
			atomic.StoreUint64(&u.incrementID, 0)
			d = atomic.AddUint64(&u.incrementID, 1)
			if u.lastTime >= currentTime.Unix() {
				time.Sleep(1001 * time.Millisecond) //防止一秒破百万造成全局不唯一
				currentTime = u.now()
			}
			if u.lastTime > currentTime.Unix() { // 时间回滚下
				if atomic.AddUint64(&u.incrementTimeID, 1) > 99 {
					atomic.StoreUint64(&u.incrementTimeID, 0)
				}
				dt = atomic.LoadUint64(&u.incrementTimeID)
			}
			atomic.StoreInt64(&u.lastTime, currentTime.Unix())
		}
	}
	u.lock.Unlock()

	if tag != "" {
		return u.prefixVal + tag + u.joinSymbol + fmt.Sprintf("%d%04d", dt, rand.Intn(10000)) + fmt.Sprintf("%06d", d) + u.joinSymbol + currentTime.Format("060102150405")
	} else {
		return u.prefixVal + fmt.Sprintf("%d%04d", dt, rand.Intn(10000)) + fmt.Sprintf("%06d", d) + u.joinSymbol + currentTime.Format("060102150405")
	}
}

// ID 固定长度 LF-A-814707000013-210803222335
func (u *Unique) ID() string {
	return u.IDWithTag("")
}

func (u *Unique) FreeIDWithTag(tag string) string {
	d := atomic.AddUint64(&u.incrementID, 1)
	if tag != "" {
		return u.prefixVal + tag + u.joinSymbol + fmt.Sprintf("%d%04d", atomic.LoadUint64(&u.incrementTimeID), rand.Intn(10000)) + fmt.Sprintf("%04d", d) + u.joinSymbol + time.Now().Format("060102150405")
	} else {
		return u.prefixVal + fmt.Sprintf("%d%04d", atomic.LoadUint64(&u.incrementTimeID), rand.Intn(10000)) + fmt.Sprintf("%04d", d) + u.joinSymbol + time.Now().Format("060102150405")
	}
}

// FreeID 不限长度
func (u *Unique) FreeID() string {
	return u.FreeIDWithTag("")
}

func (u *Unique) now() time.Time {
	return time.Now()
}

//Md5ID 不限长度
func (u *Unique) Md5ID() string {
	data := sha256.Sum256([]byte(u.FreeIDWithTag("")))
	return hex.EncodeToString(data[:])
}

//Md5TagID 不限长度
func (u *Unique) Md5TagID(tag string) string {
	data := sha256.Sum256([]byte(u.FreeIDWithTag(tag)))
	return hex.EncodeToString(data[:])
}

func (u *Unique) CharUpper(cnt int) string {
	var res string
	for i := 0; i < cnt; i++ {
		res += fmt.Sprintf("%c", 65+rand.Intn(26))
	}
	return res
}

func (u *Unique) CharLower(cnt int) string {
	var res string
	for i := 0; i < cnt; i++ {
		res += fmt.Sprintf("%c", 97+rand.Intn(26))
	}
	return res
}
