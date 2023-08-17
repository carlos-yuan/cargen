package e

func (e Err) GetMsg() string {
	switch e.Code {
	case ParamsValidatorErrorCode:
		return "参数验证有误"
	case ParamsDealErrorCode:
		return "参数处理失败"
	}
	return "未知错误"
}
