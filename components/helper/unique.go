package helper

import (
	"fmt"
	"sync/atomic"
	"time"
)

var _id int64

func UniqueID(prefix string) string {
	atomic.AddInt64(&_id, 1)
	return fmt.Sprintf("%s%d%s", prefix, _id, time.Now().Format("60102150405"))
}
