package repository

import (
	"context"
	"go-graphql-product-svc/config"

	"github.com/redis/go-redis/v9"
)

type RedisRepository interface {
	StoreToken(ctx context.Context, key string, token string) error
	GetToken(ctx context.Context, key string) (string, error)
}

type redisRepository struct {
	RC  *redis.Client
	cfg *config.Config
}

func NewRedisRepository(rc *redis.Client, cfg *config.Config) *redisRepository {
	return &redisRepository{
		RC:  rc,
		cfg: cfg,
	}
}

func (ar *redisRepository) StoreToken(ctx context.Context, key string, token string) error {
	return ar.RC.Set(ctx, key, token, ar.cfg.AuthTokenConfig.Duration).Err()
}

func (ar *redisRepository) GetToken(ctx context.Context, key string) (string, error) {
	return ar.RC.Get(ctx, key).Result()
}
