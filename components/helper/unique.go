package helper

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var _id int64
var _start string
var _uniqueLock sync.Mutex

func UniqueIDStart() string {
	return time.Now().Format("60102150405")
}
func UniqueID(prefix string) string {
	if _start == "" {
		_uniqueLock.Lock()
		_start = UniqueIDStart()
		_uniqueLock.Unlock()
	}
	atomic.AddInt64(&_id, 1)
	return fmt.Sprintf("%s%d%s", prefix, _id, _start)
}
