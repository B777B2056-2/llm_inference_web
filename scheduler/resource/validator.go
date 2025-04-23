package resource

import (
	"github.com/go-playground/validator/v10"
)

var Validator *validator.Validate

func initValidator() {
	Validator = validator.New()
	//err := Validator.RegisterValidation("checkToken", CheckToken)
	//if err != nil {
	//	panic(errors.New("failed to init checkToken against validator"))
	//}
}
