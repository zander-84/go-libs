package helper

import (
	"github.com/zander-84/go-libs/common"
	"log"
	"time"
)

type Time struct {
	location *time.Location
}

//"2006-01-02 15:04:05"
func NewTime(timeZone string) *Time {
	this := new(Time)
	if timeZone == "" {
		timeZone = common.DefaultTimeZone
	}
	var err error
	this.location, err = time.LoadLocation(timeZone)
	if err != nil {
		log.Fatal("err timeZone: ", timeZone)
	}
	return this
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

func (t *Time) UnixTime(unix int64) time.Time {
	return time.Unix(unix, 0).In(t.location)
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
