package middleware

import (
	"fmt"
	"llm_online_interence/llmgateway/client"
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
func canPassUserLimiter(ctx *gin.Context, svcName string) (bool, error) {
	uri := ctx.Param("proxyPath")
	// 从user服务获取用户信息
	tokenString := ctx.GetHeader("Authorization")
	clt, err := client.NewUserCenterGRPCClient()
	if err != nil {
		return false, err
	}
	userInfo, err := clt.GetUserInfo(ctx, tokenString)
	if err != nil {
		return false, err
	}
	// 检查并消费
	err = limiter.CheckAndConsumeWithUserID(ctx, svcName, uri, fmt.Sprintf("%d", userInfo.Id))
	if err != nil {
		return false, err
	}
	return true, nil
}

func RateLimit(svcName string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !canPassSvcLimiter(ctx, svcName) {
			ctx.Status(http.StatusTooManyRequests)
			ctx.Abort()
			return
		}

		ok, err := canPassUserLimiter(ctx, svcName)
		if err != nil {
			ctx.Status(http.StatusInternalServerError)
			ctx.Abort()
			return
		}
		if !ok {
			ctx.Status(http.StatusTooManyRequests)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
