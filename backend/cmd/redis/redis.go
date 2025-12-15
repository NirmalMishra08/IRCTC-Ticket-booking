package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func RedisConnect(REDIS_DB_URL, REDIS_PASSWORD string) *redis.Client {
	ctx := context.Background()
	fmt.Print(REDIS_DB_URL, REDIS_PASSWORD)
	rdb := redis.NewClient(&redis.Options{
		Addr:     REDIS_DB_URL,
		Username: "default",
		Password: REDIS_PASSWORD,
		DB:       0,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("Failed to connect to Redis: %v\n", err)
		return nil
	}

	fmt.Printf("Redis database started at %s\n", REDIS_DB_URL)

	return rdb

}
