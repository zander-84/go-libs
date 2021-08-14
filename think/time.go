package think

import (
	"sync"
	"time"
)

var setDong8TimeOnce sync.Once

// SetDong8Time 东八区时间
func SetDong8Time() {
	setDong8TimeOnce.Do(func() {
		var cstZone = time.FixedZone("CST", 8*3600) // 东八
		time.Local = cstZone
	})
}
