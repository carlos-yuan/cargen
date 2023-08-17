package timeUtil

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const (
	FormatNormal     = "2006-01-02 15:04:05"
	FormatNormalDate = "2006-01-02"
)

func NewTs(year int, month time.Month, day, hour, min, sec, nsec int) time.Time {
	return time.Date(year, month, day, hour, min, sec, nsec, time.FixedZone("CST", 8*60*60))
}

func StrToInt() int64 {
	format := time.Now().Format("20060102150405.000Z07")
	now, _ := strconv.ParseInt(format[:14]+format[15:18]+format[19:21], 10, 64)
	return now
}

func NowHourTime() time.Time {
	now := time.Now()
	zn, offset := now.Zone()
	return time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, time.FixedZone(zn, offset))
}

func NowDayTime() time.Time {
	now := time.Now()
	zn, offset := now.Zone()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.FixedZone(zn, offset))
}

func ToHourTime(t time.Time) time.Time {
	zn, offset := t.Zone()
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, time.FixedZone(zn, offset))
}

func ToDayTime(t time.Time) time.Time {
	zn, offset := t.Zone()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.FixedZone(zn, offset))
}

func ParseTime(val string, format ...string) (time.Time, error) {
	if len(format) == 0 {
		if len(val) == 19 {
			format = append(format, FormatNormal)
		} else if len(val) == 10 {
			format = append(format, FormatNormalDate)
		} else {
			return time.Time{}, errors.New("format not nil")
		}
	}
	return time.ParseInLocation(format[0], val, time.Local)
}

func FormatTime(t time.Time, format string) string {
	return FormatDate(t.UnixMilli(), format)
}

func ParseTimeNormal(t time.Time) string {
	return FormatDate(t.UnixMilli(), "yyyy-MM-dd HH:mm:ss")
}

func Milli() int64 {
	return time.Now().UnixMilli()
}

func Nano() int64 {
	return time.Now().UnixNano()
}

func NowTo(dateStyle string) string {
	return FormatDate(time.Now().UnixNano()/int64(time.Millisecond), dateStyle)
}

func NowNormal() string {
	return FormatDate(time.Now().UnixNano()/int64(time.Millisecond), "yyyy-MM-dd HH:mm:ss")
}

func DayStart(t int64) int64 {
	date := FormatDate(t, "yyyyMMdd") + "000000"
	return ToTimeStamp(date, "yyyyMMddHHmmss")
}

func MonthStart(t int64) int64 {
	date := FormatDate(t, "yyyyMM") + "01000000"
	return ToTimeStamp(date, "yyyyMMddHHmmss")
}

func NextMonth(t int64) int64 {
	month, _ := strconv.ParseInt(FormatDate(t, "MM"), 10, 32)
	year, _ := strconv.ParseInt(FormatDate(t, "yyyy"), 10, 32)
	if month+1 > 12 {
		year++
		month = 1
	} else {
		month++
	}
	var date string
	if month > 9 {
		date = strconv.Itoa(int(year)) + strconv.Itoa(int(month)) + "01000000"
	} else {
		date = strconv.Itoa(int(year)) + "0" + strconv.Itoa(int(month)) + "01000000"
	}
	return ToTimeStamp(date, "yyyyMMddHHmmss")
}

func ToyyyyMM() string {
	return FormatDate(Milli(), "yyyyMM")
}

func ToyyyyMMInt() int32 {
	i, _ := strconv.ParseInt(FormatDate(Milli(), "yyyyMM"), 10, 64)
	return int32(i)
}

func ToyyyyMMdd() string {
	return FormatDate(Milli(), "yyyyMMdd")
}

func ToyyyyMMddHH() string {
	return FormatDate(Milli(), "yyyyMMddHH")
}

func ToyyyyMMddHHmmss() string {
	return FormatDate(Milli(), "yyyyMMddHHmmss")
}

func ToyyyyMMddInt() int32 {
	i, _ := strconv.ParseInt(FormatDate(Milli(), "yyyyMMdd"), 10, 64)
	return int32(i)
}

func ToyyyyMMddHHInt() int32 {
	i, _ := strconv.ParseInt(FormatDate(Milli(), "yyyyMMddHH"), 10, 32)
	return int32(i)
}

func ToyyyyMMddHHIntByStamp(time int64) int32 {
	i, _ := strconv.ParseInt(FormatDate(time, "yyyyMMddHH"), 10, 32)
	return int32(i)
}

func ToyyyyMMddIntByStamp(time int64) int32 {
	i, _ := strconv.ParseInt(FormatDate(time, "yyyyMMdd"), 10, 32)
	return int32(i)
}

func ToyyyyMMIntByStamp(time int64) int32 {
	i, _ := strconv.ParseInt(FormatDate(time, "yyyyMM"), 10, 32)
	return int32(i)
}

func FormatDate(t int64, dateStyle string) string {
	date := time.Unix(t/1000, 0)
	dateStyle = strings.Replace(dateStyle, "yyyy", "2006", 1)
	dateStyle = strings.Replace(dateStyle, "yy", "06", 1)
	dateStyle = strings.Replace(dateStyle, "MM", "01", 1)
	dateStyle = strings.Replace(dateStyle, "dd", "02", 1)
	dateStyle = strings.Replace(dateStyle, "HH", "15", 1)
	dateStyle = strings.Replace(dateStyle, "mm", "04", 1)
	dateStyle = strings.Replace(dateStyle, "ss", "05", 1)
	dateStyle = strings.Replace(dateStyle, "SSS", "000", -1)
	return date.Format(dateStyle)
}

func FormatStyle(style string) string {
	style = strings.Replace(style, "yyyy", "2006", 1)
	style = strings.Replace(style, "yy", "06", 1)
	style = strings.Replace(style, "MM", "01", 1)
	style = strings.Replace(style, "dd", "02", 1)
	style = strings.Replace(style, "HH", "15", 1)
	style = strings.Replace(style, "mm", "04", 1)
	style = strings.Replace(style, "ss", "05", 1)
	style = strings.Replace(style, "SSS", "000", -1)
	return style
}

func ToTimeStamp(t string, dateStyle string) int64 {
	dateStyle = strings.Replace(dateStyle, "yyyy", "2006", 1)
	dateStyle = strings.Replace(dateStyle, "yy", "06", 1)
	dateStyle = strings.Replace(dateStyle, "MM", "01", 1)
	dateStyle = strings.Replace(dateStyle, "dd", "02", 1)
	dateStyle = strings.Replace(dateStyle, "HH", "15", 1)
	dateStyle = strings.Replace(dateStyle, "mm", "04", 1)
	dateStyle = strings.Replace(dateStyle, "ss", "05", 1)
	dateStyle = strings.Replace(dateStyle, "SSS", "000", -1)
	date, err := time.ParseInLocation(dateStyle, t, time.Local)
	if err != nil {
		println(err.Error())
	}
	return date.UnixNano() / int64(time.Millisecond)
}

const (
	Millisecond int64 = 1
	Second      int64 = 1000 * Millisecond
	Minute      int64 = 60 * Second
	Minute2     int64 = 2 * Minute
	Minute3     int64 = 3 * Minute
	Minute5     int64 = 5 * Minute
	Minute10    int64 = 10 * Minute
	Minute15    int64 = 15 * Minute
	Minute20    int64 = 20 * Minute
	Minute30    int64 = 30 * Minute
	Hour        int64 = 60 * Minute
	Day         int64 = 24 * Hour
	Week        int64 = 7 * Day
	Month       int64 = 30 * Day
)
