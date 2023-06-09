package redis

import (
	"AuthService/config"
	"github.com/go-redis/redis/v8"
)

func InitRedisClient(cfg *config.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		//Addr:     cfg.Redis.Addr,
		//Password: cfg.Redis.Password,
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return rdb
}
