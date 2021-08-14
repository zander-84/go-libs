package helper

import (
	"errors"
	"time"
)

type Time struct {
}

//NewTime "2006-01-02 15:04:05"
func NewTime() *Time {
	this := new(Time)
	return this
}

func (t *Time) Now() time.Time {
	return time.Now()
}

func (t *Time) Location() *time.Location {
	return time.Local
}

func (t *Time) Format(layout string) string {
	return t.Now().Format(layout) //"2006-01-02 15:04:05"
}

//Parse layout "01/02/2006", sourceTime "02/08/2020"
func (t *Time) Parse(layout string, sourceTime string) (time.Time, error) {
	return time.Parse(layout, sourceTime)
}

func (t *Time) DayTime(calcDay int) time.Time {
	day := t.Now().AddDate(0, 0, calcDay).Format("2006-01-02")
	today, _ := time.Parse("2006-01-02", day)
	return today
}

//UnixToTime  时间戳转时间 sec秒和 nsec纳秒
func (t *Time) UnixToTime(sec int64, nsec int64) time.Time {
	return time.Unix(sec, nsec)
}

//Unix 时间戳秒 例：1628127247
func (t *Time) Unix() int64 {
	return t.Now().Unix()
}

//UnixNano 时间戳纳秒 例：1628127212214105000
func (t *Time) UnixNano() int64 {
	return t.Now().UnixNano()
}

//UnixMsec 时间戳毫秒级 例：1628127443176
func (t *Time) UnixMsec() int64 {
	return t.Now().UnixNano() / 1e6
}

//TimeToUnixNano 时间转时间戳纳秒级
//例：timeIns.TimeToUnixNano("2006-01-02T15:04:05.999999999", "2006-01-02T15:04:05.999999999")
//返回结果 1136185445999999999
func (t *Time) TimeToUnixNano(source string, layout string) (int64, error) {
	sourceTime, err := time.Parse(layout, source)
	if err != nil {
		return 0, err
	}
	return sourceTime.UnixNano(), err
}

//TimeToUnix 时间转时间戳秒级
//例：TimeToUnix("2010-03-04", "2006-01-02") = 1265126400
func (t *Time) TimeToUnix(source string, layout string) (int64, error) {
	unixNano, err := t.TimeToUnixNano(source, layout)
	if err != nil {
		return 0, err
	}
	return unixNano / 1e9, nil
}

//TimeToUnixMsec 时间转时间戳毫秒级
//例：TimeToUnixMsec("2010-03-04", "2006-01-02") = 1265126400000
func (t *Time) TimeToUnixMsec(source string, layout string) (int64, error) {
	unixNano, err := t.TimeToUnixNano(source, layout)
	if err != nil {
		return 0, err
	}
	return unixNano / 1e6, nil
}

func (t *Time) Year() string {
	return t.Now().Format("2006")
}

func (t *Time) Month() string {
	return t.Now().Format("01")
}

func (t *Time) Day() string {
	return t.Now().Format("02")
}

func (t *Time) FormatSlash() string {
	return t.Now().Format("2006/01/02 15:04:05")
}

func (t *Time) FormatHyphen() string {
	return t.Now().Format("2006-01-02 15:04:05")
}

func (t *Time) FormatDaySlash() string {
	return t.Now().Format("2006/01/02")
}

func (t *Time) FormatDayHyphen() string {
	return t.Now().Format("2006-01-02")
}

func (t *Time) FormatSlashFromTime(timer time.Time) string {
	if timer.IsZero() {
		return ""
	}
	return timer.Format("2006/01/02 15:04:05")
}

func (t *Time) FormatHyphenFromTime(timer time.Time) string {
	if timer.IsZero() {
		return ""
	}
	return timer.Format("2006-01-02 15:04:05")
}

func (t *Time) FormatDaySlashFromTime(timer time.Time) string {
	if timer.IsZero() {
		return ""
	}
	return timer.Format("2006/01/02")
}

func (t *Time) FormatDayHyphenFromTime(timer time.Time) string {
	if timer.IsZero() {
		return ""
	}
	return timer.Format("2006-01-02")
}
func (t *Time) IntervalSliceTimes(startAt time.Time, endAt time.Time, interval time.Duration) ([]time.Time, error) {
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

func (t *Time) IntervalArrayTimes(startAt time.Time, endAt time.Time, interval time.Duration) ([][2]time.Time, error) {
	ts, err := t.IntervalSliceTimes(startAt, endAt, interval)
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
