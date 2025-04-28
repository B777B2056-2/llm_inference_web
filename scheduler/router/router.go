package router

import (
	"github.com/gin-gonic/gin"
	"llm_online_inference/scheduler/controller"
)

func Init(r *gin.Engine) {
	chatGroup := r.Group("chat")
	chatGroup.Use(userIDMiddleware)
	chatGroup.POST("completion", controller.ChatCompletion) // 对话流式接口（SSE）
	chatGroup.POST("history", controller.ChatHistory)       // 对话历史记录

	batchInferenceGroup := r.Group("batchInference")
	batchInferenceGroup.Use(userIDMiddleware)
	batchInferenceGroup.POST("create", controller.CreateBatchInferenceTask)      // 创建批量推理任务
	batchInferenceGroup.POST("results", controller.GetBatchInferenceTaskResults) // 获取批量推理任务结果
}
