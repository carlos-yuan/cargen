package gin

import (
	"github.com/gin-gonic/gin"
)

// 404页面
func NoRoute() gin.HandlerFunc {

	return func(c *gin.Context) {
		c.JSON(404, map[string]interface{}{
			"code": 404,
			"msg":  "内容不存在",
			"path": c.Request.RequestURI,
		})
	}
}
