package compound

import (
	"sync"
	"time"
)

var Location *time.Location
var mutex sync.Mutex

type Data struct {
	Value        string `json:"value"`
	Typ          int    `json:"typ"` //0：字符串  10：时间 20:相对图片 21绝对图片 30换行
	Days         int    `json:"days"`
	Months       int    `json:"months"`
	Years        int    `json:"years"`
	Format       string `json:"format"`
	MaxWidth     string `json:"max_width"`
	FontSize     string `json:"font_size"`
	ResizeWidth  uint   `json:"resize_width"`
	ResizeHeight uint   `json:"resize_height"`
	PositionX    int    `json:"position_x"`
	PositionY    int    `json:"position_y"`
	IsRise       bool   `json:"is_rise"`
	Remark       string `json:"remark"`
}

func (this *Data) IsString() bool {
	if this.Typ == 0 {
		return true
	}
	return false
}

func (this *Data) IsTime() bool {
	if this.Typ == 10 {
		return true
	}
	return false
}

func (this *Data) IsRelativeImage() bool {
	if this.Typ == 20 {
		return true
	}
	return false
}

func (this *Data) IsAbsoluteImage() bool {
	if this.Typ == 21 {
		return true
	}
	return false
}

func (this *Data) IsLine() bool {
	if this.Typ == 30 {
		return true
	}
	return false
}

func (this *Data) GetTimeString() string {
	if Location == nil {
		mutex.Lock()
		Location, _ = time.LoadLocation("Asia/Shanghai")
		mutex.Unlock()
	}
	return time.Now().In(Location).AddDate(this.Years, this.Months, this.Days).Format(this.Format)
}

func (this *Data) GetString() string {
	return this.Value
}
