package ratelimiter

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRateLimiter struct {
	Redis *redis.Client
}

func NewRedisRateLimiter(redisClient *redis.Client) *RedisRateLimiter {
	return &RedisRateLimiter{
		Redis: redisClient,
	}
}



func (h *RedisRateLimiter) RateLimitUser(ctx context.Context, userId string, window time.Duration, maxRequests int) error {
	key := fmt.Sprintf("rate_limit:%s:%d", userId, time.Now().Unix()/int64(window.Seconds()))

	// Using pipeline for atomic operation
	pipe := h.Redis.Pipeline()
	pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, window)

	results, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}

	currentCount := results[0].(*redis.IntCmd).Val()
	if currentCount > int64(maxRequests) {
		return fmt.Errorf("rate limit exceeded: %d requests in %v", currentCount, window)
	}

	return nil
}

// Enhanced version with sliding window
func (h *RedisRateLimiter) RateLimitUserSlidingWindow(ctx context.Context, userId string, window time.Duration, maxRequests int) error {
	now := time.Now().UnixMilli()
	windowMs := window.Milliseconds()
	key := fmt.Sprintf("rate_limit_sliding:%s", userId)

	pipe := h.Redis.Pipeline()

	// Remove old timestamps
	pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", now-windowMs))

	// Add current timestamp
	pipe.ZAdd(ctx, key, redis.Z{
		Score:  float64(now),
		Member: fmt.Sprintf("%d", now),
	})

	// Get count
	pipe.ZCard(ctx, key)

	// Set expiry
	pipe.Expire(ctx, key, window)

	results, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}

	count := results[2].(*redis.IntCmd).Val()
	if count > int64(maxRequests) {
		return fmt.Errorf("rate limit exceeded: %d requests in sliding window %v", count, window)
	}

	return nil
}
