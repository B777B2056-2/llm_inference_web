package dto

type LoginReq struct {
	Username  string `json:"username" validate:"required"`
	Password  string `json:"password" validate:"required"` // sha256后的用户密码
	CaptchaID string `json:"captchaID" validate:"required"`
	Answer    string `json:"answer" validate:"required"`
}

type UpdateUserTokenReq struct {
	Token string `json:"token" validate:"required,checkToken"`
}

type UserLogoutReq struct {
	Token string `json:"token" validate:"required,checkToken"`
}

type UserSignUpReq struct {
	Username  string `json:"username" validate:"required"`
	Password  string `json:"password" validate:"required"` // sha256后的用户密码
	Password2 string `json:"password2" validate:"required"`
	CaptchaID string `json:"captchaID" validate:"required"`
	Answer    string `json:"answer" validate:"required"`
}
