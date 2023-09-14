package cartime

import (
	"strconv"
	"strings"
	"time"

	"github.com/carlos-yuan/cargen/util/convert"
)

const timeFormat = "20060102150405.000Z07"

const DefaultFormat = "2006-01-02 15:04:05"
const DefaultFormatDate = "2006-01-02"

func NowToInt() int64 {
	return toInt(time.Now().Format(timeFormat))
}

func toInt(format string) int64 {
	var timebytes = make([]byte, 0, 19)
	timebytes = append(timebytes, format[:14]...)
	timebytes = append(timebytes, format[15:18]...)
	if len(format) > 20 {
		timebytes = append(timebytes, format[19:21]...)
	}
	intTime, _ := strconv.ParseInt(convert.Bytes2str(timebytes), 10, 64)
	return intTime
}

func StrToInt(str, layout string) int64 {
	t, err := time.Parse(layout, str)
	if err != nil {
		return 0
	}
	if t.Location() == nil || !strings.Contains(layout, "Z07") { //时间如果没有时区信息，默认加载本地时区
		t = t.In(time.Local)
	}
	return toInt(t.Format(timeFormat))
}

func IntToStr(t int64, layout string) string {
	str := strconv.Itoa(int(t))
	var timebytes = make([]byte, 0, 24)
	timebytes = append(timebytes, str[:14]...)
	timebytes = append(timebytes, '.', str[14], str[15], str[16], '+', str[17], str[18])
	ct, err := time.Parse(timeFormat, convert.Bytes2str(timebytes))
	if err != nil {
		return ""
	}
	return ct.Format(layout)
}

// TimeStr 特殊时间格式转换使用 yyyyMMddHHmmssSSSZ 2023072518181866608 年-月-日-时-分-秒-毫秒-时区(0-23 超过11为-值，12为-1)
type TimeStr string
