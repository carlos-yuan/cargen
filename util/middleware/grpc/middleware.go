package middleware

import (
	"comm/controller/token"
	e "comm/error"
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/bytedance/gopkg/cloud/metainfo"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/cloudwego/kitex/pkg/remote"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
)

func WhiteIpMiddleware(whiteList map[string][]string) server.Option {
	return server.WithMiddleware(func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request, response interface{}) error {
			info := rpcinfo.GetRPCInfo(ctx)
			if info.From().Address().Network() == "tcp" {
				ips := whiteList[info.Invocation().MethodName()]
				if len(ips) > 0 {
					find := false
					for _, ip := range ips {
						if strings.Index(info.From().Address().String(), ip) == 0 {
							find = true
							break
						}
					}
					if !find {
						return errors.New("invalid ip address")
					}
				}
			}
			err := next(ctx, request, response)
			return err
		}
	})
}

// ErrClientHandlerOption 业务错误再封装
var ErrClientHandlerOption = client.WithErrorHandler(func(ctx context.Context, err error) error {
	transErr, ok := err.(*remote.TransError)
	if ok && transErr.TypeID() == 6 {
		errStr := strings.ReplaceAll(transErr.Error(), "biz error: ", "")
		var errInfo e.Err
		err = json.Unmarshal([]byte(errStr), &errInfo)
		if err == nil {
			return errInfo
		}
		err = e.RPCClientErrorCodeError.SetErr(err, strings.ReplaceAll(transErr.Error(), "biz error: ", ""))
	} else {
		err = e.RPCClientErrorCodeError.SetErr(err, "数据访问失败")
	}
	return err
})

var ErrServerHandlerOption = server.WithErrorHandler(func(ctx context.Context, err error) error {
	transErr, ok := err.(*kerrors.DetailedError)
	if ok {
		var errInfo e.Err
		if transErr.As(&errInfo) {
			var b []byte
			b, err = json.Marshal(errInfo)
			if err != nil {
				err = e.RPCServerErrorCodeError.SetErr(err, "错误转换失败")
			} else {
				err = e.RPCServerErrorCodeError.SetErr(err, string(b))
			}
		}
	} else {
		err = e.RPCServerErrorCodeError.SetErr(err, "数据访问失败")
	}
	return err
})

var ClientMetaHandler = client.WithMetaHandler(&MetaTTHeaderHandler{})

var ServerMetaHandler = server.WithMetaHandler(&MetaTTHeaderHandler{})

type MetaTTHeaderHandler struct {
	WhiteApi   []string //接口白名单
	CheckWhite bool
}

func (ch *MetaTTHeaderHandler) WriteMeta(ctx context.Context, msg remote.Message) (context.Context, error) {
	tk := ctx.Value(token.GrpcTokenStringKey)
	hd := make(map[string]string)
	if metainfo.HasMetaInfo(ctx) {
		metainfo.SaveMetaInfoToMap(ctx, hd)
	}
	if tk != nil {
		hd[token.GrpcTokenStringKey] = tk.(*token.Payload).GrpcToken
	}
	msg.TransInfo().PutTransStrInfo(hd)
	return ctx, nil
}

// ReadMeta of MetaTTHeaderHandler reads headers of TTHeader protocol from transport
func (ch *MetaTTHeaderHandler) ReadMeta(ctx context.Context, msg remote.Message) (context.Context, error) {
	tk := msg.TransInfo().TransStrInfo()[token.GrpcTokenStringKey]
	// 接口白名单验证
	if ch.CheckWhite {
		method := msg.RPCInfo().Invocation().MethodName()
		for _, s := range ch.WhiteApi {
			if s == method {
				return ctx, nil
			}
		}

	}
	return context.WithValue(ctx, token.GrpcTokenKey, tk), nil
}
