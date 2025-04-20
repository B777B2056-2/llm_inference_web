package proxy

import (
	"llm_online_interence/llmgateway/confparser"
	"llm_online_interence/llmgateway/middleware"

	"github.com/gin-gonic/gin"
)

// initRouter 初始化路由
func initRouter(r *gin.Engine) {
	for _, backend := range confparser.ProxyConfig.Backends {
		group := r.Group(backend.GroupName)
		group.Use(middleware.UserAuth(backend))
		group.Use(middleware.RefreshUserToken(backend))
		group.Use(middleware.RateLimit(backend.SvcName))
		group.Use(middleware.Breaker(backend.SvcName)) // 断路器需保证为最后一个中间件，否则无法生效
		group.Any("/*proxyPath", func(ctx *gin.Context) {
			commonProxyHandler(ctx, backend)
		})
	}
}
