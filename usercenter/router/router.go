package router

import (
	"github.com/gin-gonic/gin"
	"llm_inference_web/usercenter/controller"
)

func Init(r *gin.Engine) {
	userGroup := r.Group("user")
	userGroup.POST("/logout", controller.Logout)                   // 用户注销登录状态
	userGroup.POST("/login", controller.Login)                     // 用户登录
	userGroup.POST("/signUp", controller.SignUp)                   // 用户注册
	userGroup.GET("/captcha/generate", controller.GenerateCaptcha) // 生成验证码
}
