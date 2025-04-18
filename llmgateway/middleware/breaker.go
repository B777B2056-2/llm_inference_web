package middleware

import (
	"llm_online_interence/llmgateway/breaker"

	"github.com/gin-gonic/gin"
)

func Breaker(svcName string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		breaker := breaker.GetBreakerBySvcName(svcName)
		if breaker == nil {
			return
		}
		if !breaker.Execute(ctx) {
			ctx.Abort()
		}
	}
}
