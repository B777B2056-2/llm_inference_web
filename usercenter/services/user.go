package services

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"llm_online_inference/usercenter/confparser"
	"llm_online_inference/usercenter/dao"
	"llm_online_inference/usercenter/dto"
	"llm_online_inference/usercenter/resource"
	"llm_online_inference/usercenter/utils"
	"time"
)

func getUserTokenRedisKey(userToken string) string {
	const userTokenRedisKeyPrefix = "user_token_"
	return userTokenRedisKeyPrefix + userToken
}

// UserLogin 用户登录：检查验证码 -> 检查用户是否存在 -> 检查密码是否匹配 -> 生成访问令牌 -> 令牌存入redis，记录用户状态
func UserLogin(ctx context.Context, params dto.LoginReq) (resp dto.LoginResp, err error) {
	// 检查验证码是否有效
	if !verifyCaptcha(params.CaptchaID, params.Answer) {
		return resp, utils.NewError(1003)
	}

	// 检查用户是否存在，以及密码是否匹配
	userInfo, err := dao.NewUserDao().GetByName(params.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return resp, utils.NewError(1001)
		}
		return resp, err
	}
	if userInfo.Password != params.Password {
		return resp, utils.NewError(1002)
	}

	// 生成jwt令牌
	tokenString, err := utils.GenerateToken(userInfo.ID, params.Username)
	if err != nil {
		return resp, err
	}

	// 存入redis，记录用户登录状态，并设置过期时间
	key := getUserTokenRedisKey(tokenString)
	expirationTime := time.Duration(confparser.ResourceConfig.Server.TokenExpirationInSecond) * time.Second
	err = resource.RedisClient.SetEX(ctx, key, userInfo.ID, expirationTime).Err()
	if err != nil {
		return resp, err
	}

	resp.Token = tokenString
	return resp, nil
}

// GetUserInfo 获取用户信息
func GetUserInfo(ctx context.Context, tokenString string) (dto.UserInfoResp, error) {
	// 检查redis内用户访问令牌是否存在
	key := getUserTokenRedisKey(tokenString)
	userID, err := resource.RedisClient.Get(ctx, key).Uint64()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return dto.UserInfoResp{}, nil
		}
		return dto.UserInfoResp{}, err
	}

	// 查询数据库
	userInfo, err := dao.NewUserDao().GetByID(uint(userID))
	if err != nil {
		// 未找到记录，删除redis内用户访问令牌
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = resource.RedisClient.Del(ctx, key)
			return dto.UserInfoResp{}, utils.NewError(1001)
		}
		return dto.UserInfoResp{}, err
	}

	return dto.UserInfoResp{UserID: userInfo.ID, Username: userInfo.Name}, nil
}

// CheckUserAuth 用户鉴权
func CheckUserAuth(ctx context.Context, tokenString string) (bool, error) {
	key := getUserTokenRedisKey(tokenString)
	_, err := resource.RedisClient.Exists(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// UpdateUserToken 更新用户活跃状态：重置redis TTL
func UpdateUserToken(ctx context.Context, tokenString string) error {
	key := getUserTokenRedisKey(tokenString)
	expirationTime := time.Duration(confparser.ResourceConfig.Server.TokenExpirationInSecond) * time.Second
	return resource.RedisClient.Expire(ctx, key, expirationTime).Err()
}

// UserLogout 用户注销登录状态
func UserLogout(ctx context.Context, tokenString string) error {
	key := getUserTokenRedisKey(tokenString)
	err := resource.RedisClient.Del(ctx, key).Err()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}
		return err
	}
	return nil
}

func UserSignUp(_ context.Context, params dto.UserSignUpReq) error {
	// 检查验证码是否有效
	if !verifyCaptcha(params.CaptchaID, params.Answer) {
		return utils.NewError(1003)
	}

	// 检查用户是否已经存在
	userInfo, err := dao.NewUserDao().GetByName(params.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if userInfo != nil {
		return utils.NewError(1004)
	}

	// 创建用户记录
	return dao.NewUserDao().Create(params.Username, params.Password)
}
