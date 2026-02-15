package booking

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func (h *Handler) RateLimitUser(ctx context.Context, userId string, window time.Duration, maxRequests int) error {
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
