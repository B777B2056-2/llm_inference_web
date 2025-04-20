package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"llm_online_inference/usercenter/confparser"
	"net/http"
)

type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func newSuccessResponse(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

func newErrorResponse(ctx *gin.Context, err error) {
	var realErr *confparser.ErrorConfig
	if errors.As(err, &realErr) {
		ctx.JSON(http.StatusOK, APIResponse{
			Code:    realErr.Code,
			Message: realErr.Message,
		})
	} else {
		ctx.JSON(http.StatusOK, APIResponse{
			Code:    500,
			Message: err.Error(),
		})
	}
}

func NewResponse(ctx *gin.Context, data any, err error) {
	if err != nil {
		newErrorResponse(ctx, err)
		return
	}
	newSuccessResponse(ctx, data)
}

func NewError(code int) error {
	for _, e := range confparser.Errors {
		if e.Code == code {
			return &e
		}
	}
	return errors.New("unknown error")
}
