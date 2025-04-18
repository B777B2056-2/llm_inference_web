package proxy

import "github.com/gin-gonic/gin"

func Init(r *gin.Engine) {
	initRouter(r)
}
