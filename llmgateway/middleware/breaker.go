package middleware

import (
	"llm_online_interence/llmgateway/breaker"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Breaker(svcName string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		b := breaker.GetBreakerBySvcName(svcName)
		if b == nil {
			ctx.Next()
			return
		}
		if !b.Execute(ctx) {
			ctx.Status(http.StatusTooManyRequests)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
