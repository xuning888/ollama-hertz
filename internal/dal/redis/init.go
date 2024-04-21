package redis

import (
	"github.com/redis/go-redis/v9"
	"github.com/xuning888/ollama-hertz/pkg/config"
)

var Client *redis.Client

func Init() {
	Client = redis.NewClient(&redis.Options{
		Addr:     config.DefaultConfig.Redis.Addr,
		Password: config.DefaultConfig.Redis.Password,
		DB:       config.DefaultConfig.Redis.DB,
	})
}
