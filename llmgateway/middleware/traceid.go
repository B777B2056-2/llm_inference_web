package middleware

import (
	"llm_online_interence/llmgateway/resource"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

// TraceID 生成traceID的中间件
func TraceID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		traceID := ctx.Request.FormValue("reqID")
		if traceID == "" {
			uid, _ := uuid.NewV4()
			ctx.Request.Form.Set("reqID", uid.String())
			traceID = uid.String()
		}
		resource.Logger.AddHook(resource.NewTraceIdHook(traceID))
		ctx.Writer.Header().Set("X-Trace-Id", traceID)
		ctx.Next()
	}
}
