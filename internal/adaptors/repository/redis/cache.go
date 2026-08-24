package redis

import (
	"context"
	"hamrahTask1/internal/core/ports"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	rdb *redis.Client
}

func NewRedisCache(rdb *redis.Client) ports.SessionCache {
	return &RedisCache{rdb: rdb}
}

func (r *RedisCache) CreateSession(ctx context.Context, token, userID string) error {
	return r.rdb.Set(ctx, token, userID, 24*time.Hour).Err()
}

func (r *RedisCache) GetUserID(ctx context.Context, token string) (string, error) {
	return r.rdb.Get(ctx, token).Result()
}

func (r *RedisCache) DeleteSession(ctx context.Context, token string) error {
	return r.rdb.Del(ctx, token).Err()
}
