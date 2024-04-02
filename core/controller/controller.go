package ctl

import (
	"context"
	"errors"

	e "github.com/carlos-yuan/cargen/core/error"

	"github.com/jinzhu/copier"
)

// ControllerContext 请求体接口
type ControllerContext interface {
	// CheckToken 检查token是否有效 该方法出错会panic 并被panic handler拦截
	CheckToken(Token)
	// GetToken 获取token包含的内容
	GetToken() Payload
	// GetHeader 获取header
	GetHeader(key string) string
	// SetHeader 设置header
	SetHeader(key, val string)
	// GetRequestInfo 获取请求方法，和链接地址
	GetRequestInfo() (method, url string)
	// Bind 绑定参数 该方法出错会panic 并被panic handler拦截
	Bind(params any, bindings ...BindOption)
	// SetContext 设置web框架的上下文对象
	SetContext(context.Context) ControllerContext
	// GetContext 获得上下文
	GetContext() context.Context
	// 设置加密方法
	SetEncryption(func(ctx ControllerContext, data any) (any, error))
	// 成功返回统一接口
	Success(data any) *Result
}

// Result 统一返回体
type Result struct {
	Code      int    `json:"code"`      //状态码 200正常 500错误
	Msg       string `json:"msg"`       //错误消息
	Data      any    `json:"data"`      //返回数据
	RequestId string `json:"requestId"` //请求编号链路中来
	err       e.Err
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
func (r *Result) Err(args ...any) *Result {
	var err error
	var msg string
	if len(args) == 2 {
		err = args[0].(error)
		msg = args[1].(string)
	} else if len(args) == 1 {
		var ok bool
		msg, ok = args[0].(string)
		if !ok {
			err = args[0].(error)
		}
	}
	if err != nil {
		if errors.As(err, &r.err) {
			r.err = err.(e.Err)
			if msg != "" {
				r.Msg = msg
				r.Data = err.Error()
			} else {
				r.Msg = r.err.ErrMsg
			}
		} else {
			r.err = r.err.SetErr(err, msg)
			r.Data = err.Error()
			r.Msg = r.err.ErrMsg
		}
	} else {
		if msg != "" {
			r.Msg = msg
		}
	}
	r.Code = 500
	return r
}

func Copy[T any](to T, from any, opts ...copier.Option) T {
	if len(opts) > 0 {
		err := copier.CopyWithOption(to, from, opts[0])
		if err != nil {
			panic(e.ParamsDealError.SetErr(err, "参数转换失败"))
		}
	} else if len(opts) == 0 {
		err := copier.Copy(to, from)
		if err != nil {
			panic(e.ParamsDealError.SetErr(err, "参数转换失败"))
		}
	} else {
		panic(e.ParamsDealError.SetErr(nil, "拷贝选项有误"))
	}
	return to
}
