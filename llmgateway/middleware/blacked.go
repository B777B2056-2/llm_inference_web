package middleware

import (
	"llm_online_interence/llmgateway/confparser"
	"net/http"

	"github.com/gin-gonic/gin"
)

func BlackedIPs() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		isBlackedIP := false
		for _, ip := range confparser.ProxyConfig.BlackedIPs {
			if ctx.ClientIP() == ip {
				isBlackedIP = true
				break
			}
		}
		if isBlackedIP {
			ctx.Status(http.StatusForbidden)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
