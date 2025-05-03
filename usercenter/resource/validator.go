package resource

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"llm_inference_web/usercenter/utils"
)

var Validator *validator.Validate

func CheckToken(fl validator.FieldLevel) bool {
	tokenString := fl.Field().String()
	_, err := utils.ValidateToken(tokenString)
	return err == nil
}

func initValidator() {
	Validator = validator.New()
	err := Validator.RegisterValidation("checkToken", CheckToken)
	if err != nil {
		panic(errors.New("failed to init checkToken against validator"))
	}
}
