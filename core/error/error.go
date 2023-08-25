package e

import (
	"bytes"
	"encoding/json"
	"errors"
	"runtime/debug"
	"strings"

	"github.com/bytedance/gopkg/util/logger"
	"gorm.io/gorm"
)

type Err struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Site   string `json:"site"`   //发生位置
	Err    error  `json:"-"`      //包裹错误
	ErrMsg string `json:"errMsg"` //包裹错误信息
}

const (
	DuplicateEntry = "Duplicate entry"
)

var UpdateFailure = Err{Msg: "更新失败"}

func (e Err) SetErr(err error, msg ...string) Err {
	if len(msg) == 1 {
		e.Msg = msg[0]
		if err == gorm.ErrRecordNotFound {
			e.Code = NotFindErrorCode
		}
	} else {
		if strings.Contains(err.Error(), DuplicateEntry) {
			e.Code = InternalServerErrorCode
			e.Msg = "请勿重复操作！"
		} else if err == gorm.ErrRecordNotFound {
			e.Code = NotFindErrorCode
			e.Msg = "未找到内容！"
		} else if e.Code == 0 {
			e.Code = InternalServerErrorCode
			e.Msg = "访问失败！"
		}
	}
	sites := strings.Split(string(debug.Stack()), "\n")
	e.Site = strings.ReplaceAll(sites[6], "\t", "")
	e.Site = e.Site[strings.LastIndex(e.Site[:strings.LastIndex(e.Site, "/")], "/")+1:]
	e.Err = err
	if err != nil {
		e.ErrMsg = err.Error()
	}
	return e
}

func (e Err) SetRecover(r interface{}) Err {
	rStr, ok := r.(string)
	if !ok {
		b, _ := json.Marshal(r)
		rStr = string(b)
	}
	e.Err = errors.New(rStr)
	e.Msg = rStr
	e.Site = GetSite(2)
	return e
}

func (e Err) Error() string {
	return e.Msg
}

func (e Err) Is(err error) bool {
	errInfo, ok := err.(Err)
	if ok {
		return errInfo.Code == NotFindErrorCode
	}
	return false
}

func PrintError(err error) {
	var buffer bytes.Buffer
	printError(&buffer, err)
}

func printError(siteBuffer *bytes.Buffer, err error) {
	me, ok := err.(Err)
	if ok {
		if siteBuffer.Len() == 0 {
			siteBuffer.WriteString(me.CodeName())
			siteBuffer.WriteString(":")
			siteBuffer.WriteString(err.Error())
			siteBuffer.WriteString(":")
		}
		siteBuffer.WriteString(me.Site)
		siteBuffer.WriteString(" -> ")
		if me.Err != nil {
			printError(siteBuffer, me.Err)
		} else {
			logger.Error(siteBuffer.String() + ":  " + me.Msg)
		}
	} else {
		logger.Error(siteBuffer.String() + ":  " + err.Error())
	}
}

// GetSite 获取行所在stack 1为当前位置
func GetSite(line int) string {
	line = line*2 + 4
	sites := strings.Split(string(debug.Stack()), "\n")
	site := strings.ReplaceAll(sites[line], "\t", "")
	site = site[strings.LastIndex(site[:strings.LastIndex(site, "/")], "/")+1:]
	return site
}

type WarpErr struct {
	Error Err
}

func (e WarpErr) Err(err error, msg ...string) Err {
	return e.Error.SetErr(err, msg...)
}
