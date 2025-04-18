package middleware

import (
	"llm_online_interence/llmgateway/limiter"
	"net/http"

	"github.com/gin-gonic/gin"
)

// canPassSvcLimiter 判断服务限流（水平限流）
func canPassSvcLimiter(ctx *gin.Context, svcName string) bool {
	uri := ctx.Param("proxyPath")
	if err := limiter.CheckAndConsumeWithSvc(ctx, svcName, uri); err != nil {
		return false
	}
	return true
}

// canPassUserLimiter 判断用户限流（垂直限流/多租户限流），防止某一个用户单独占满整个服务资源
func canPassUserLimiter(ctx *gin.Context, svcName string) bool {
	uri := ctx.Param("proxyPath")
	// TODO 从user服务获取用户信息
	userID := ""
	// 检查并消费
	if err := limiter.CheckAndConsumeWithUserID(ctx, svcName, uri, userID); false && err != nil {
		return false
	}
	return true
}

func RateLimit(svcName string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !canPassUserLimiter(ctx, svcName) || !canPassSvcLimiter(ctx, svcName) {
			ctx.Status(http.StatusTooManyRequests)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
