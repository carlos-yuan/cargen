package gin

import (
	"github.com/alibaba/sentinel-golang/core/system"
	sentinel "github.com/alibaba/sentinel-golang/pkg/adapters/gin"
	"github.com/gin-gonic/gin"
)

// Sentinel 限流
func Sentinel() gin.HandlerFunc {
	if _, err := system.LoadRules([]*system.Rule{
		{
			MetricType:   system.AvgRT,
			TriggerCount: 1000,
			Strategy:     system.NoAdaptive,
		},
	}); err != nil {
		println("Unexpected error: %+v", err)
	}
	return sentinel.SentinelMiddleware(
		sentinel.WithBlockFallback(func(ctx *gin.Context) {
			ctx.AbortWithStatusJSON(200, map[string]interface{}{
				"msg":  "too many request; the quota used up!",
				"code": 500,
			})
		}),
	)
}
