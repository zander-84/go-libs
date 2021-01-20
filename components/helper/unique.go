package helper

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	mathrand "math/rand"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var _id uint64
var _start = time.Now().Format("060102150405")
var _uniqueLock sync.Mutex
var Prefix string

func UniqueID() string {
	return fmt.Sprintf("%s%d%s", Prefix, atomic.AddUint64(&_id, 1), _start)
}

func Rand(count int, raw bool) string {
	maxTry := 5
	var err error
	var res int64

	for i := 0; i < maxTry; i++ {
		if raw {
			res, err = _rand1(int64(math.Pow10(count)))
			if err == nil {
				return fmt.Sprintf("%0"+fmt.Sprintf("%d", count)+"d", res)
			}
		} else {
			res, err = _rand1(int64(math.Pow10(count)) - int64(math.Pow10(count-1)))
			if err == nil {
				res = res + int64(math.Pow10(count-1))
				return fmt.Sprintf("%0"+fmt.Sprintf("%d", count)+"d", res)
			}
		}
	}

	return _rand2(count)
}

func _rand1(max int64) (int64, error) {
	result, err := rand.Int(rand.Reader, big.NewInt(max))
	return result.Int64(), err
}

func _rand2(count int) string {
	nums := make([]string, 0)
	for i := 0; i < count; i++ {
		mathrand.Seed(time.Now().UnixNano())
		nums = append(nums, fmt.Sprintf("%d", mathrand.Intn(10)))
	}
	return strings.Join(nums, "")
}
