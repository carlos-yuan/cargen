package e

import (
	"strconv"
	"sync"
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

	NoPermissionErrorCode = 1008
)

// 控制器处理出错 错误码
var (
	ParamsValidatorError = Err{Code: ParamsValidatorErrorCode}

	ParamsDealError = Err{Code: ParamsDealErrorCode}

	NotFindError = Err{Code: NotFindErrorCode}

	AuthorizeError = Err{Code: AuthorizeErrorCode}

	AuthorizeTimeOutError = Err{Code: AuthorizeErrorCode}

	NoPermissionError = Err{Code: NoPermissionErrorCode}

	InternalServerError = Err{Code: InternalServerErrorCode}

	RPCClientErrorCodeError = Err{Code: RPCClientErrorCode}

	RPCServerErrorCodeError = Err{Code: RPCServerErrorCode}
)

var codeNameMap map[int]string = map[int]string{
	InternalServerErrorCode:   "服务内部错误:",
	RPCClientErrorCode:        "服务连接错误:",
	AuthorizeErrorCode:        "授权校验失败:",
	AuthorizeTimeOutErrorCode: "授权失效:",
	NoPermissionErrorCode:     "没有权限",
}

var codeNameMutex sync.Mutex

func SetCodeNameMap(code int, name string) {
	codeNameMutex.Lock()
	codeNameMap[code] = name
	codeNameMutex.Unlock()
}

func (e Err) CodeName() string {
	for k, v := range codeNameMap {
		if e.Code == k {
			return v + strconv.Itoa(e.Code)
		}
	}
	return "Unknown err:" + strconv.Itoa(e.Code)
}

var (
	RecordNotFound = "record not found"
)
