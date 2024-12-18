package util

import (
	"go-graphql-product-svc/config"

	"github.com/redis/go-redis/v9"
)

func GetRedisClient(cfg *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisConfig.Address,
		Password: cfg.RedisConfig.Password,
		DB:       cfg.RedisConfig.DB,
	})
	return client
}
