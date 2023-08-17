package e

import (
	"strconv"
)

const (
	ParamsValidatorErrorCode = 1000

	ParamsDealErrorCode = 1001

	NotFindErrorCode = 1002

	AuthorizeErrorCode = 1003

	InternalServerErrorCode = 1004

	RPCClientErrorCode = 1005

	RPCServerErrorCode = 1006

	AuthorizeTimeOutErrorCode = 1007

	EdgeUserGrpcErrorCode = 2001

	EdgePermitGrpcErrorCode = 2002

	EdgeDictGrpcErrorCode = 2003

	EdgeProductGrpcErrorCode = 2004
)

// 控制器处理出错 错误码
var (
	ParamsValidatorError = Err{Code: ParamsValidatorErrorCode}

	ParamsDealError = Err{Code: ParamsDealErrorCode}

	NotFindError = Err{Code: NotFindErrorCode}

	AuthorizeError = Err{Code: AuthorizeErrorCode}

	AuthorizeTimeOutError = Err{Code: AuthorizeErrorCode}

	InternalServerError = Err{Code: InternalServerErrorCode}

	RPCClientErrorCodeError = Err{Code: RPCClientErrorCode}

	RPCServerErrorCodeError = Err{Code: RPCServerErrorCode}
)

// 业务rpc对应错误码
var (
	EdgeUserGrpcError = Err{Code: EdgeUserGrpcErrorCode}

	EdgePermitGrpcError = Err{Code: EdgePermitGrpcErrorCode}

	EdgeDictGrpcError = Err{Code: EdgeDictGrpcErrorCode}

	EdgeProductGrpcError = Err{Code: EdgeProductGrpcErrorCode}
)

func (e Err) CodeName() string {
	switch e.Code {
	case EdgeUserGrpcErrorCode:
		return "用户GRPC服务:" + strconv.Itoa(EdgeUserGrpcErrorCode)
	case EdgePermitGrpcErrorCode:
		return "许可证GRPC服务:" + strconv.Itoa(EdgePermitGrpcErrorCode)
	case InternalServerErrorCode, NotFindErrorCode:
		return "服务内部错误:" + strconv.Itoa(InternalServerErrorCode)
	case RPCClientErrorCode:
		return "服务连接错误:" + strconv.Itoa(RPCClientErrorCode)
	case AuthorizeErrorCode:
		return "授权校验失败:" + strconv.Itoa(AuthorizeErrorCode)
	case AuthorizeTimeOutErrorCode:
		return "授权失效:" + strconv.Itoa(AuthorizeTimeOutErrorCode)
	}
	return "Unknown err:" + strconv.Itoa(e.Code)
}

var (
	RecordNotFound = "record not found"
)
