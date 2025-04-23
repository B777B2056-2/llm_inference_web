package controller

import (
	"llm_online_inference/scheduler/resource"
	"net/http"

	"llm_online_inference/scheduler/dto"
	"llm_online_inference/scheduler/services"

	"github.com/gin-gonic/gin"
)

func ChatCompletion(ctx *gin.Context) {
	userIdVal, ok := ctx.Get("user_id")
	if !ok {
		ctx.Status(http.StatusUnauthorized)
		ctx.Abort()
		return
	}
	userId := userIdVal.(int)

	// 解析 POST 数据
	var params dto.ChatCompletionReq
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if err := resource.Validator.Struct(params); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// sse返回数据
	if err := services.ChatCompletion(ctx, userId, &params); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}
}

func ChatHistory(ctx *gin.Context) {
	userIdVal, ok := ctx.Get("user_id")
	if !ok {
		ctx.Status(http.StatusUnauthorized)
		ctx.Abort()
		return
	}
	userId := userIdVal.(int)

	var params struct {
		ChatSessionID string `json:"chat_session_id" validate:"required"`
	}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if err := resource.Validator.Struct(params); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := services.GetChatHistory(ctx, userId, params.ChatSessionID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": resp})
}
