package middleware

import (
	"github.com/gin-gonic/gin"
	"llm_online_interence/llmgateway/client"
	"llm_online_interence/llmgateway/confparser"
	"net/http"
)

func UserAuth(backendConf confparser.BackendConfigItem) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uri := ctx.Param("proxyPath")
		needAuth := false
		for _, u := range backendConf.NeedAuthURLs {
			if u == uri {
				needAuth = true
			}
		}
		if !needAuth {
			ctx.Next()
			return
		}
		// 调用鉴权服务
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
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
		hasAuth, err := clt.CheckUserAuth(ctx, tokenString)
		if err != nil {
			ctx.Status(http.StatusInternalServerError)
			ctx.Abort()
			return
		}
		if !hasAuth {
			ctx.Status(http.StatusUnauthorized)
			ctx.Abort()
			return
		}

		ctx.Set("token", tokenString)
		ctx.Next()
	}
}
