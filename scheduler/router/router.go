package router

import (
	"github.com/gin-gonic/gin"
	"llm_online_inference/scheduler/controller"
)

func Init(r *gin.Engine) {
	chatGroup := r.Group("chat")
	chatGroup.Use(userIDMiddleware)
	chatGroup.POST("completion", controller.ChatCompletion) // 对话流式接口（SSE）
}
