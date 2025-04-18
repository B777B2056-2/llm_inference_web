package limiter

import (
	"context"
	"llm_online_interence/llmgateway/resource"
	"time"
)

const luaScript = `
-- redis_token_bucket.lua
local key = KEYS[1] -- 令牌桶的键名
local rate = tonumber(ARGV[1]) 	    -- 令牌生成速率 (每秒生成的令牌数)
local capacity = tonumber(ARGV[2])  -- 令牌桶的容量
local now = tonumber(ARGV[3])       -- 当前时间戳（毫秒）
local requested = tonumber(ARGV[4]) -- 请求消耗的令牌数

-- 获取桶的信息
local last_state = redis.call('HMGET', key, 'tokens', 'timestamp')
local last_tokens = tonumber(last_state[1]) or capacity
local last_refreshed = tonumber(last_state[2]) or now

-- 计算生成的令牌数
local delta = math.max(0, now - last_refreshed) * rate / 1000
local filled_tokens = math.min(capacity, last_tokens + delta)

-- 检查是否有足够的令牌
if filled_tokens < requested then
    return 0 	-- 令牌不足，返回失败
else
    -- 更新令牌桶信息
    filled_tokens = filled_tokens - requested
    redis.call('HMSET', key, 'tokens', filled_tokens, 'timestamp', now)
    return 1	-- 请求成功
end
`

var luaScriptSHA string

func Init() {
	// 预加载lua脚本到redis
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	luaScriptSHA, _ = resource.RedisClient.ScriptLoad(ctx, luaScript).Result()
	cancel()
}
