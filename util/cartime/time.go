package cartime

import (
	"strconv"
	"time"

	"github.com/carlos-yuan/cargen/util/convert"
)

const timeFormat = "20060102150405.000Z07:00"

const DefaultFormat = "2006-01-02 15:04:05"
const DefaultFormatDate = "2006-01-02"

func NowToInt() int64 {
	format := time.Now().Format(timeFormat)
	var timebytes = make([]byte, 0, 19)
	timebytes = append(timebytes, format[:14]...)
	timebytes = append(timebytes, format[15:16]...)
	timebytes = append(timebytes, format[19:21]...)
	timebytes = append(timebytes, format[22:24]...)
	intTime, _ := strconv.ParseInt(convert.Bytes2str(timebytes), 10, 64)
	return intTime
}

func StrToInt(str, layout string) int64 {
	t, err := time.Parse(layout, str)
	if err != nil {
		return 0
	}
	format := t.Format(timeFormat)
	var timebytes = make([]byte, 0, 19)
	timebytes = append(timebytes, format[:14]...)
	timebytes = append(timebytes, format[15:16]...)
	if len(format) == 19 {
		t = time.Now()
		_, o := t.Zone()
		hour := o / 3600
		if hour >= 10 {
			hourStr := strconv.Itoa(hour)
			timebytes = append(timebytes, hourStr[0])
			timebytes = append(timebytes, hourStr[1])
		} else {
			timebytes = append(timebytes, '0', strconv.Itoa(hour)[0])
		}
		min := o % 3600 / 60
		if min >= 10 {
			minStr := strconv.Itoa(min)
			timebytes = append(timebytes, minStr[0])
			timebytes = append(timebytes, minStr[1])
		} else {
			timebytes = append(timebytes, '0', strconv.Itoa(min)[0])
		}
	} else {
		timebytes = append(timebytes, format[19:21]...)
		timebytes = append(timebytes, format[22:24]...)
	}
	intTime, _ := strconv.ParseInt(convert.Bytes2str(timebytes), 10, 64)
	return intTime
}

func IntToStr(t int64, layout string) string {
	str := strconv.Itoa(int(t))
	var timebytes = make([]byte, 0, 24)
	timebytes = append(timebytes, str[:14]...)
	timebytes = append(timebytes, '.', str[14], '0', '0', '+', str[15], str[16], ':', str[17], str[18])
	ct, err := time.Parse(timeFormat, convert.Bytes2str(timebytes))
	if err != nil {
		return ""
	}
	return ct.Format(layout)
}

// TimeStr 特殊时间格式转换使用 yyyyMMddHHmmssSSSZ 2023072518181866608 年-月-日-时-分-秒-毫秒-时区(0-23 超过11为-值，12为-1)
type TimeStr string
