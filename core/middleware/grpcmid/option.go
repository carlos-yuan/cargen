package grpcmid

import (
	"context"
	"strings"

	e "github.com/carlos-yuan/cargen/core/error"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/remote"
	"github.com/cloudwego/kitex/pkg/remote/trans/nphttp2/status"
)

// ErrHandlerOption 业务错误再封装
var ErrHandlerOption = client.WithErrorHandler(func(ctx context.Context, err error) error {
	transErr, ok := err.(*remote.TransError)
	if ok && transErr.TypeID() == 6 {
		err = e.RPCClientErrorCodeError.SetErr(err, strings.ReplaceAll(transErr.Error(), "biz error: ", ""))
	} else {
		statusErr, ok := err.(*status.Error)
		if ok && uint32(statusErr.GRPCStatus().Code()) == 13 {
			err = e.RPCClientErrorCodeError.SetErr(err, strings.TrimSpace(strings.ReplaceAll(statusErr.GRPCStatus().Message(), "[biz error]", "")))
		} else {
			err = e.RPCClientErrorCodeError.SetErr(err, "数据访问失败")
		}
	}
	return err
})
