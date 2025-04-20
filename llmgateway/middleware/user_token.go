package middleware

import (
	"github.com/gin-gonic/gin"
	"llm_online_interence/llmgateway/client"
	"llm_online_interence/llmgateway/confparser"
	"net/http"
)

func RefreshUserToken(backendConf confparser.BackendConfigItem) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !backendConf.NeedRefreshToken {
			ctx.Next()
			return
		}

		tokenString, ok := ctx.Get("token")
		if !ok {
			ctx.Status(http.StatusUnauthorized)
			ctx.Abort()
			return
		}

		clt, err := client.NewUserCenterGRPCClient()
		if err != nil {
			ctx.Status(http.StatusInternalServerError)
			ctx.Abort()
			return
		}
		if err := clt.UpdateUserToken(ctx, tokenString.(string)); err != nil {
			ctx.Status(http.StatusInternalServerError)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
