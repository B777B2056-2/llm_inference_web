package resource

import (
	"context"
	"fmt"
	"llm_online_interence/llmgateway/confparser"
	"time"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func initRedis() {
	conf := confparser.ResourceConfig
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", conf.Redis.Host, conf.Redis.Port),
		Password:     conf.Redis.Password,
		DB:           conf.Redis.DB,
		PoolSize:     conf.Redis.PoolSize,
		DialTimeout:  time.Duration(conf.Redis.DialTimeout) * time.Second,
		ReadTimeout:  time.Duration(conf.Redis.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(conf.Redis.WriteTimeout) * time.Second,
		MaxRetries:   conf.Redis.ConnMaxRetries,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
	RedisClient = rdb
}
