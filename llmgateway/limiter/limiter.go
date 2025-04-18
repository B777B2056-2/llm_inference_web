package limiter

import (
	"context"
	"errors"
	"llm_online_interence/llmgateway/confparser"
	"llm_online_interence/llmgateway/resource"
	"time"
)

// checkAndConsumeImpl 尝试消费一个令牌，如果令牌桶为空则返回错误
func checkAndConsumeImpl(ctx context.Context, redisKey string, tokenPerSecond, bucketSize, requestNum uint) error {
	result, err := resource.RedisClient.EvalSha(
		ctx, luaScriptSHA,
		[]string{redisKey},
		tokenPerSecond, bucketSize, time.Now().UnixMilli(), requestNum).Result()
	if err != nil {
		return err
	}
	if result == int64(0) {
		return errors.New("Rate limit exceeded")
	}
	return nil
}

// buildSvcRateLimiterKey 构建服务限流器redis键
func buildSvcRateLimiterKey(svcName, uri string) string {
	return RedisKeyPrefix + svcName + ":" + uri
}

// CheckAndConsumeWithSvc 根据服务名和uri检查并消费令牌
func CheckAndConsumeWithSvc(ctx context.Context, svcName, uri string) error {
	svcRateLimitConf, ok := confparser.ProxyConfig.SvcURIRateLimit[svcName]
	if !ok {
		return nil
	}
	uriRateLimitConf, ok := svcRateLimitConf[uri]
	if !ok {
		return nil
	}
	tokenPerSecond, bucketSize := uriRateLimitConf.TokenPerSecond, uriRateLimitConf.BucketSize
	return checkAndConsumeImpl(ctx, buildSvcRateLimiterKey(svcName, uri), tokenPerSecond, bucketSize, 1)
}

// buildUserRateLimiterKey 构建服务单用户限流器redis键
func buildUserRateLimiterKey(svcName, userID string) string {
	return RedisKeyPrefix + svcName + ":" + userID
}

// CheckAndConsumeWithUserID 根据服务名和用户id检查并消费令牌
func CheckAndConsumeWithUserID(ctx context.Context, svcName, uri, userID string) error {
	if userID == "" {
		return errors.New("userID is empty")
	}
	svcRateLimitConf, ok := confparser.ProxyConfig.UserRateLimit[svcName]
	if !ok {
		return nil
	}
	userRateLimitConf, ok := svcRateLimitConf[uri]
	if !ok {
		return nil
	}
	tokenPerSecond, bucketSize := userRateLimitConf.TokenPerSecond, userRateLimitConf.BucketSize
	return checkAndConsumeImpl(ctx, buildUserRateLimiterKey(svcName, userID), tokenPerSecond, bucketSize, 1)
}
