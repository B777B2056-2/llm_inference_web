package dto

type GenerateCaptchaResp struct {
	CaptchaID string `json:"captchaID"`
}

type LoginResp struct {
	Token string `json:"token"`
}

type UserInfoResp struct {
	UserID   uint   `json:"userId"`
	Username string `json:"username"`
}
