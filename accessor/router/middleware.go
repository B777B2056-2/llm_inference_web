package router

import (
	"github.com/gin-gonic/gin"
	"llm_online_inference/accessor/resource"
	"net/http"
	"strconv"
)

// sseMiddleware 设置 SSE 头
func sseMiddleware(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Next()
}

// userIDMiddleware 获取用户id
func userIDMiddleware(ctx *gin.Context) {
	userIdStr := ctx.GetHeader("user_id")
	if userIdStr == "" {
		ctx.Status(http.StatusUnauthorized)
		ctx.Abort()
		return
	}
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		ctx.Abort()
		return
	}
	ctx.Set("user_id", userId)
	ctx.Next()
}

func traceIdMiddleware(ctx *gin.Context) {
	traceId := ctx.GetHeader("X-Trace-Id")
	if traceId == "" {
		ctx.Status(http.StatusUnauthorized)
		ctx.Abort()
		return
	}
	ctx.Set("trace_id", traceId)
	resource.Logger.AddHook(resource.NewTraceIdHook(traceId))
	ctx.Next()
}
