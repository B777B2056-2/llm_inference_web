package controller

import (
	"llm_online_inference/accessor/resource"
	"net/http"

	"llm_online_inference/accessor/dto"
	"llm_online_inference/accessor/services"

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
	if err := services.NewOnlineInferenceOperator(userId).ChatCompletion(ctx, &params); err != nil {
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

	var params dto.ChatHistoryReq
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if err := resource.Validator.Struct(params); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := services.NewOnlineInferenceOperator(userId).GetChatHistory(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": resp})
}

func CreateBatchInferenceTask(ctx *gin.Context) {
	userIdVal, ok := ctx.Get("user_id")
	if !ok {
		ctx.Status(http.StatusUnauthorized)
		ctx.Abort()
		return
	}
	userId := userIdVal.(int)

	var params dto.CreateBatchInferenceTaskReq
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if err := resource.Validator.Struct(params); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := services.NewBatchInferenceOperator(userId).CreateTask(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}
	ctx.Status(http.StatusOK)
}

func GetBatchInferenceTaskResults(ctx *gin.Context) {
	userIdVal, ok := ctx.Get("user_id")
	if !ok {
		ctx.Status(http.StatusUnauthorized)
		ctx.Abort()
		return
	}
	userId := userIdVal.(int)

	var params dto.GetBatchInferenceTaskResultsReq
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if err := resource.Validator.Struct(params); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := services.NewBatchInferenceOperator(userId).TaskResults(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": resp})
}
