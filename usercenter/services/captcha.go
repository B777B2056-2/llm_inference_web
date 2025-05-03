package services

import (
	"context"
	"github.com/mojocn/base64Captcha"
	"llm_inference_web/usercenter/confparser"
	"llm_inference_web/usercenter/resource"
	"time"
)

// RedisStore 实现 base64Captcha.Store 接口
type RedisStore struct {
	keyPrefix string // Redis 键前缀（如 "captcha:"）
	expire    time.Duration
}

// NewRedisStore 创建 Redis 存储实例
func NewRedisStore(keyPrefix string, expire time.Duration) *RedisStore {
	return &RedisStore{
		keyPrefix: keyPrefix,
		expire:    expire,
	}
}

// Set 保存验证码答案到 Redis
func (s *RedisStore) Set(id string, value string) error {
	key := s.keyPrefix + id
	return resource.RedisClient.Set(context.Background(), key, value, s.expire).Err()
}

// Get 从 Redis 获取验证码答案
func (s *RedisStore) Get(id string, clear bool) string {
	key := s.keyPrefix + id
	val, err := resource.RedisClient.Get(context.Background(), key).Result()
	if err != nil {
		return ""
	}

	if clear {
		// 验证后删除键
		go resource.RedisClient.Del(context.Background(), key)
	}
	return val
}

// Verify 验证答案并删除键
func (s *RedisStore) Verify(id, answer string, clear bool) bool {
	storedAnswer := s.Get(id, clear)
	return storedAnswer == answer
}

// 验证码存储
var store = NewRedisStore(
	"user_center_captcha_",
	time.Duration(confparser.ResourceConfig.Server.CaptchaExpirationInSecond)*time.Second,
)

// GenerateDigitCaptcha 生成数字验证码（返回验证码ID、Base64图片、答案）
func GenerateDigitCaptcha() (string, error) {
	// 配置验证码参数
	driver := base64Captcha.DriverDigit{
		Height:   80,  // 图片高度
		Width:    240, // 图片宽度
		Length:   6,   // 验证码长度
		MaxSkew:  0.7, // 数字最大倾斜角度
		DotCount: 100, // 干扰点数量
	}

	// 生成验证码
	captcha := base64Captcha.NewCaptcha(&driver, store)
	id, _, _, err := captcha.Generate()
	if err != nil {
		return "", err
	}
	return id, nil
}

// verifyCaptcha 验证用户输入
func verifyCaptcha(id, answer string) bool {
	return store.Verify(id, answer, true) // true: 验证后删除缓存
}
