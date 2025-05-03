package utils

import (
	"fmt"
	"time"

	"llm_inference_web/usercenter/confparser"

	"github.com/golang-jwt/jwt/v5"
)

// UserClaims 签发用户token请求
type UserClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken 生成Token
func GenerateToken(userID uint, username string) (string, error) {
	// 定义 Token 过期时间（例如 24 小时）
	expirationTime := time.Now().Add(24 * time.Hour)

	// 创建 Claims
	claims := UserClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "your-app-name", // 可自定义签发者
			Subject:   "user-auth",     // 可自定义主题
		},
	}

	// 创建 Token 并签名
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(confparser.ResourceConfig.Server.TokenSecretKey)
	if err != nil {
		return "", fmt.Errorf("生成 Token 失败: %w", err)
	}

	return tokenString, nil
}

// ValidateToken 验证用户token
func ValidateToken(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("意外的签名方法: %v", token.Header["alg"])
		}
		return confparser.ResourceConfig.Server.TokenSecretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
