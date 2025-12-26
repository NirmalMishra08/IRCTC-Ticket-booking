package booking

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func (h *Handler) RateLimiterUser(ctx context.Context, userId string, window time.Duration, maxRequest int64) error {
	pipe := h.Redis.Pipeline()
	key := fmt.Sprintf("rate_limit:%s:%d", userId, time.Now().Unix()/int64(window.Seconds()))

	pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, window)

	results, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}

	currentCount := results[0].(*redis.IntCmd).Val()

	if int64(currentCount) > maxRequest {
		return fmt.Errorf("rate limit exceeded: %d requests in %v", currentCount, window)
	}

	return nil;
}
