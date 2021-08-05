package helper

import (
	"github.com/zander-84/go-libs/think"
	"log"
	"time"
)

type Time struct {
	location *time.Location
	timeZone string
}

//"2006-01-02 15:04:05"
func NewTime(timeZone string) *Time {
	this := new(Time)
	if timeZone == "" {
		timeZone = think.DefaultTimeZone
	}
	var err error
	this.location, err = time.LoadLocation(timeZone)
	if err != nil {
		log.Fatal("err timeZone: ", timeZone)
	}
	this.timeZone = timeZone
	return this
}
func (t *Time) TimeZone() string {
	return t.timeZone
}
func (t *Time) Location() *time.Location {
	return t.location
}

func (t *Time) Now() time.Time {
	return time.Now().In(t.location)
}

func (t *Time) LocationNow(src time.Time) time.Time {
	return src.In(t.location)
}
func (t *Time) Format(layout string) string {
	return t.Now().Format(layout) //"2006-01-02 15:04:05"
}

// layout "01/02/2006", sourceTime "02/08/2020"
func (t *Time) Parse(layout string, sourceTime string) (time.Time, error) {
	t1, err := time.ParseInLocation(layout, sourceTime, t.location)
	if err != nil {
		return time.Time{}, err
	}
	return t1.In(t.location), err
}

func (t *Time) DayTime(calcDay int) time.Time {
	day := t.Now().AddDate(0, 0, calcDay).Format("2006-01-02")
	today, _ := time.ParseInLocation("2006-01-02", day, t.location)
	return today
}

//UnixToTime  时间戳转时间 sec秒和 nsec纳秒
func (t *Time) UnixToTime(sec int64, nsec int64) time.Time {
	return time.Unix(sec, nsec).In(t.location)
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
	sourceTime, err := time.ParseInLocation(layout, source, t.location)
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
	return timer.In(t.location).Format("2006/01/02 15:04:05")
}

func (t *Time) FormatHyphenFromTime(timer time.Time) string {
	if timer.IsZero() {
		return ""
	}
	return timer.In(t.location).Format("2006-01-02 15:04:05")
}

func (t *Time) FormatDaySlashFromTime(timer time.Time) string {
	if timer.IsZero() {
		return ""
	}
	return timer.In(t.location).Format("2006/01/02")
}

func (t *Time) FormatDayHyphenFromTime(timer time.Time) string {
	if timer.IsZero() {
		return ""
	}
	return timer.In(t.location).Format("2006-01-02")
}
