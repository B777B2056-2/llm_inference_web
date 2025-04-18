package middleware

import (
	"llm_online_interence/llmgateway/confparser"

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
			ctx.JSON(403, gin.H{
				"message": "Forbidden",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
