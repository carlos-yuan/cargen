package ginmid

import (
	"runtime/debug"

	e "github.com/carlos-yuan/cargen/core/error"

	"github.com/gin-gonic/gin"
)

// Panic 业务抛出PANIC时处理
func Panic() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				err, ok := rec.(e.Err)
				if ok {
					switch err.Code {
					case e.AuthorizeErrorCode, e.AuthorizeTimeOutErrorCode:
						err.Msg = "尚未授权"
						c.JSON(401, err)
						return
					case e.ParamsDealErrorCode, e.ParamsValidatorErrorCode:
						c.JSON(400, err)
						return
					case e.NoPermissionErrorCode:
						c.JSON(403, err)
						return
					}
				}
				println(string(debug.Stack()))
				c.JSON(500, e.InternalServerError.SetRecover(rec))
			}
		}()
		c.Next()
	}
}
