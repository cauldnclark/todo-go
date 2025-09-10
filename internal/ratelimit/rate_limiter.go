package ratelimit

import (
	"context"
	"fmt"
	"time"

	"github.com/cauldnclark/todo-go/internal/redis"
	redislib "github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	client  *redis.Client
	limiter string
}

func NewRateLimiter(client *redis.Client, limiter string) *RateLimiter {
	return &RateLimiter{
		client:  client,
		limiter: limiter,
	}
}

func (r *RateLimiter) Allow(ctx context.Context, userID int, maxRequests int, window time.Duration) (bool, int, int64, error) {
	key := fmt.Sprintf("ratelimit:user:%d", userID)
	windowInSeconds := int64(window.Seconds())

	switch r.limiter {
	case "sliding":
		return r.allowSlidingWindow(ctx, key, maxRequests, windowInSeconds)
	default:
		return r.allowFixedWindow(ctx, key, maxRequests, windowInSeconds)
	}
}

func (r *RateLimiter) allowFixedWindow(ctx context.Context, key string, maxReq int, windowSec int64) (bool, int, int64, error) {
	now := time.Now().Unix()
	resetTime := now + windowSec

	val, err := r.client.GetClient().Incr(ctx, key).Result()
	if err != nil {
		return false, 0, 0, err
	}

	if val == 1 {
		err = r.client.GetClient().Expire(ctx, key, time.Duration(windowSec)*time.Second).Err()
		if err != nil {
			return false, 0, 0, err
		}
	}

	remaining := maxReq - int(val)
	if remaining < 0 {
		remaining = 0
	}

	return val <= int64(maxReq), remaining, resetTime - now, nil
}

func (r *RateLimiter) allowSlidingWindow(ctx context.Context, key string, maxReq int, windowSec int64) (bool, int, int64, error) {
	now := time.Now().Unix()
	pipe := r.client.GetClient().TxPipeline()

	pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", now-windowSec))
	pipe.ZAdd(ctx, key, redislib.Z{
		Score:  float64(now),
		Member: now,
	})

	pipe.Expire(ctx, key, time.Duration(windowSec)*time.Second)

	pipe.ZCard(ctx, key)

	cmders, err := pipe.Exec(ctx)
	if err != nil {
		return false, 0, 0, err
	}

	count := int(cmders[len(cmders)-1].(*redislib.IntCmd).Val())
	remaining := max(maxReq-count, 0)

	resetIn := windowSec - (now % windowSec)
	if resetIn == 0 {
		resetIn = windowSec
	}

	return count <= maxReq, remaining, resetIn, nil
}
