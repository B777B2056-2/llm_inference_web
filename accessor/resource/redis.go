package resource

import (
	"context"
	"fmt"
	"llm_inference_web/accessor/confparser"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
)

var (
	RedisClient *redis.Client
	RedisLocker *redsync.Redsync
)

func initRedis() {
	conf := confparser.ResourceConfig
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", conf.Redis.Host, conf.Redis.Port),
		Password:     conf.Redis.Password,
		DB:           conf.Redis.DB,
		PoolSize:     conf.Redis.PoolSize,
		DialTimeout:  time.Duration(conf.Redis.DialTimeoutSecond) * time.Second,
		ReadTimeout:  time.Duration(conf.Redis.ReadTimeoutSecond) * time.Second,
		WriteTimeout: time.Duration(conf.Redis.WriteTimeoutSecond) * time.Second,
		MaxRetries:   conf.Redis.ConnMaxRetries,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
	RedisClient = rdb
	RedisLocker = redsync.New(goredis.NewPool(rdb))
	if RedisLocker == nil {
		panic("redis locker is nil")
	}
}
