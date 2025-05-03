package controller

import (
	"github.com/gin-gonic/gin"
	"llm_inference_web/usercenter/dto"
	"llm_inference_web/usercenter/resource"
	"llm_inference_web/usercenter/services"
	"llm_inference_web/usercenter/utils"
	"net/http"
)

func GenerateCaptcha(ctx *gin.Context) {
	var err error
	var resp dto.GenerateCaptchaResp
	resp.CaptchaID, err = services.GenerateDigitCaptcha()
	utils.NewResponse(ctx, resp, err)
}

func Login(ctx *gin.Context) {
	var params dto.LoginReq
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := services.UserLogin(ctx, params)
	utils.NewResponse(ctx, resp, err)
}

func SignUp(ctx *gin.Context) {
	var params dto.UserSignUpReq
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := resource.Validator.Struct(params); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	utils.NewResponse(ctx, nil, services.UserSignUp(ctx, params))
}

func Logout(ctx *gin.Context) {
	var params dto.UserLogoutReq
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := resource.Validator.Struct(params); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	utils.NewResponse(ctx, nil, services.UserLogout(ctx, params.Token))
}
