package think

import (
	"math/rand"
	"sync"
	"time"
)

var bsOnce sync.Once

func Bootstrap() {
	bsOnce.Do(func() {
		// SetDong8Time 东八区时间
		var cstZone = time.FixedZone("CST", 8*3600) // 东八
		time.Local = cstZone

		// 随机数
		rand.Seed(time.Now().UnixNano())
	})
}
