package ctl

import (
	"context"
	"errors"

	"github.com/carlos-yuan/cargen/core/controller/token"
	e "github.com/carlos-yuan/cargen/core/error"
	"github.com/jinzhu/copier"
)

// ControllerContext 请求体接口
type ControllerContext interface {
	// CheckToken 检查token是否有效 该方法出错会panic 并被panic handler拦截
	CheckToken()
	// GetToken 获取token包含的内容
	GetToken() *token.Payload
	// GetHeader 获取header
	GetHeader(key string) string
	// SetHeader 设置header
	SetHeader(key, val string)
	// Bind 绑定参数 该方法出错会panic 并被panic handler拦截
	Bind(params any, bindings ...BindOption)
	// SetContext 设置web框架的上下文对象
	SetContext(context.Context) ControllerContext
	// GetContext 获得上下文
	GetContext() context.Context
}

// Result 统一返回体
type Result struct {
	Code      int    `json:"code"`      //状态码 200正常 500错误
	Msg       string `json:"msg"`       //错误消息
	Data      any    `json:"data"`      //返回数据
	RequestId string `json:"requestId"` //请求编号链路中来
	err       e.Err
}

// Success 成功返回
func (r *Result) Success(data any) *Result {
	r.Data = data
	r.Code = 200
	return r
}

// String 成功返回字符串
func (r *Result) String(data string) *Result {
	r.Data = data
	r.Code = 200
	return r
}

// Bytes 成功返回字节数组
func (r *Result) Bytes(data []byte) *Result {
	r.Data = data
	r.Code = 200
	return r
}

// SuccessIn 成功拷贝回复
func (r *Result) SuccessIn(to any, from any) *Result {
	err := copier.Copy(to, from)
	if err != nil {
		return r.Err(err, "数据转换失败")
	}
	r.Data = to
	r.Code = 200
	return r
}

// Err 错误返回
func (r *Result) Err(err error, msg ...string) *Result {
	if errors.As(err, &r.err) {
		r.err = err.(e.Err)
		if len(msg) == 1 {
			r.Msg = msg[0]
		} else {
			r.Msg = r.err.Error()
		}
	} else {
		r.err = r.err.SetErr(err, msg...)
		r.Msg = r.err.Error()
	}
	r.Code = 500
	return r
}

func Copy[T any](to T, from any) T {
	err := copier.Copy(to, from)
	if err != nil {
		panic(e.ParamsDealError.SetErr(err, "参数转换失败"))
	}
	return to
}
