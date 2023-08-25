package gin

import (
	"fmt"

	"github.com/carlos-yuan/cargen/util/log"

	"github.com/gin-gonic/gin"
)

func Log(logLevel int8) gin.HandlerFunc {
	logger := log.New(logLevel)
	return gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: func(params gin.LogFormatterParams) string {
			return fmt.Sprintf("%s|%s|%d|%s|%s",
				params.Method,
				params.Path,
				params.StatusCode,
				params.ClientIP,
				params.Request.UserAgent(),
			)
		},
		Output: logger,
	})
}
