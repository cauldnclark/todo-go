package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/cauldnclark/todo-go/internal/redis"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
}

func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.GetClient().Set(ctx, key, data, expiration).Err()
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	return r.client.GetClient().Del(ctx, key).Err()
}

func (r *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	exists, err := r.client.GetClient().Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}

var ErrCacheMiss = errors.New("cache: key not found")
