package helper

import (
	"errors"
	"time"
)

var defaultTime = NewTime()

type Time struct {
}

//NewTime "2006-01-02 15:04:05"
func NewTime() *Time {
	this := new(Time)
	return this
}

func GetTime() *Time {
	return defaultTime
}

func (t *Time) Now() time.Time {
	return time.Now()
}

func (t *Time) Location() *time.Location {
	return time.Local
}

func (t *Time) NowFormat(layout string) string {
	return t.Now().Format(layout) //"2006-01-02 15:04:05"
}

//Parse layout "01/02/2006", sourceTime "02/08/2020"
func (t *Time) Parse(layout string, sourceTime string) (time.Time, error) {
	return time.ParseInLocation(layout, sourceTime, time.Local)
}

func (t *Time) GetDayTime(calcDay int) time.Time {
	day := t.Now().AddDate(0, 0, calcDay).Format("2006-01-02")
	today, _ := time.ParseInLocation("2006-01-02", day, time.Local)
	return today
}

//UnixToTime  秒时间戳秒转时间
func (t *Time) UnixToTime(sec int64) time.Time {
	return time.Unix(sec, 0)
}

//UnixMilliTime  毫秒时间戳转时间
func (t *Time) UnixMilliTime(msec int64) time.Time {
	return time.Unix(0, msec*int64(time.Millisecond))
}

func (t *Time) AddMonthToUnixMilli(month int) int64 {
	return time.Now().AddDate(0, month, 0).UnixNano() / 1e6
}

//ShouldTimeToUnix 时间戳秒 例：1628127247
func (t *Time) ShouldTimeToUnix(timer time.Time) int64 {
	if timer.IsZero() {
		return 0
	}
	return timer.Unix()
}

func (t *Time) Unix() int64 {
	return time.Now().Unix()
}

//ShouldTimeToUnixMilli 时间戳毫秒级 例：1628127443176
func (t *Time) ShouldTimeToUnixMilli(timer time.Time) int64 {
	if timer.IsZero() {
		return 0
	}
	return timer.UnixNano() / 1e6
}

func (t *Time) UnixMilli() int64 {
	return time.Now().UnixNano() / 1e6
}

//ShouldTimeToUnixNano 时间戳纳秒 例：1628127212214105000
func (t *Time) ShouldTimeToUnixNano(timer time.Time) int64 {
	if timer.IsZero() {
		return 0
	}
	return timer.UnixNano()
}

func (t *Time) UnixNano() int64 {
	return time.Now().UnixNano()
}

//ParseTimeToUnixNano 时间转时间戳纳秒级
//例：timeIns.TimeToUnixNano("2006-01-02T15:04:05.999999999", "2006-01-02T15:04:05.999999999")
//返回结果 1136185445999999999
func (t *Time) ParseTimeToUnixNano(layout string, source string) (int64, error) {
	sourceTime, err := time.ParseInLocation(layout, source, time.Local)
	if err != nil {
		return 0, err
	}
	return sourceTime.UnixNano(), err
}

//ParseTimeToUnix 时间转时间戳秒级
//例：TimeToUnix("2006-01-02","2010-03-04") = 1265126400
func (t *Time) ParseTimeToUnix(layout string, source string) (int64, error) {
	unixNano, err := t.ParseTimeToUnixNano(layout, source)
	if err != nil {
		return 0, err
	}
	return unixNano / 1e9, nil
}

//ParseTimeToUnixMsec 时间转时间戳毫秒级
//例：TimeToUnixMsec("2006-01-02","2010-03-04") = 1265126400000
func (t *Time) ParseTimeToUnixMsec(layout string, source string) (int64, error) {
	unixNano, err := t.ParseTimeToUnixNano(layout, source)
	if err != nil {
		return 0, err
	}
	return unixNano / 1e6, nil
}
func (t *Time) ParseHyphenTimeToUnixMsec(hyphenTime string) (int64, error) {
	return t.ParseTimeToUnixMsec("2006-01-02 15:04:05", hyphenTime)
}

func (t *Time) ShouldYear(timer time.Time) string {
	if timer.IsZero() {
		return ""
	}

	return timer.Format("2006")
}

func (t *Time) ShouldMonth(timer time.Time) string {
	if timer.IsZero() {
		return ""
	}

	return timer.Format("01")
}

func (t *Time) ShouldDay(timer time.Time) string {
	if timer.IsZero() {
		return ""
	}

	return timer.Format("02")
}

func (t *Time) ShouldFormatSlashFromTime(timer time.Time) string {
	if timer.IsZero() {
		return ""
	}

	return timer.Format("2006/01/02 15:04:05")
}

func (t *Time) ShouldFormatHyphenFromTime(timer time.Time) string {
	if timer.IsZero() {
		return ""
	}

	return timer.Format("2006-01-02 15:04:05")
}

//ShouldFormatMilliTimeFromInt 毫秒级时间戳转字符串
func (t *Time) ShouldFormatMilliTimeFromInt(in int64) string {
	if in == 0 {
		return ""
	}
	timer := t.UnixMilliTime(in)
	return timer.Format("2006-01-02 15:04:05")
}

//ShouldFormatTimeFromInt 秒级时间戳转字符串
func (t *Time) ShouldFormatTimeFromInt(in int64) string {
	if in == 0 {
		return ""
	}
	timer := t.UnixToTime(in)
	return timer.Format("2006-01-02 15:04:05")
}

func (t *Time) FormatHyphenFromNow() string {
	return t.Now().Format("2006-01-02 15:04:05")
}

func (t *Time) ShouldFormatDaySlashFromTime(timer time.Time) string {
	if timer.IsZero() {
		return ""
	}
	return timer.Format("2006/01/02")
}
func (t *Time) FormatDaySlashFromNow() string {
	return t.Now().Format("2006/01/02")
}

func (t *Time) ShouldFormatDayHyphenFromTime(timer time.Time) string {
	if timer.IsZero() {
		return ""
	}
	return timer.Format("2006-01-02")
}

func (t *Time) FormatDayHyphenFromNow() string {
	return t.Now().Format("2006-01-02")
}
func (t *Time) SliceTime(startAt time.Time, endAt time.Time, interval time.Duration) ([]time.Time, error) {
	if startAt.IsZero() {
		return nil, errors.New("起始时间不能为空")
	}
	if endAt.IsZero() {
		return nil, errors.New("截止时间不能为空")
	}

	duration := endAt.Sub(startAt)
	if duration <= 0 {
		return nil, errors.New("截止时间必须大于起始时间")
	}

	timeSlice := make([]time.Time, 0)
	timeSlice = append(timeSlice, startAt)
	for {
		startAt = startAt.Add(interval)
		if startAt.Before(endAt) {
			timeSlice = append(timeSlice, startAt)
		} else if startAt.Equal(endAt) {
			timeSlice = append(timeSlice, startAt)
			break
		} else {
			timeSlice = append(timeSlice, endAt)
			break
		}
	}
	return timeSlice, nil
}

func (t *Time) SliceArrayTime(startAt time.Time, endAt time.Time, interval time.Duration) ([][2]time.Time, error) {
	ts, err := t.SliceTime(startAt, endAt, interval)
	if err != nil {
		return nil, err
	}

	if len(ts) < 2 {
		return nil, errors.New("时间切片数量过小")
	}
	var res = make([][2]time.Time, 0)

	tsLen := len(ts)
	for k, _ := range ts {
		if k == tsLen-1 {
			break
		}
		res = append(res, [2]time.Time{ts[k], ts[k+1]})
	}
	return res, nil
}

func (t *Time) FormatUnsigned() string {
	return t.Now().Format("20060102150405")
}
