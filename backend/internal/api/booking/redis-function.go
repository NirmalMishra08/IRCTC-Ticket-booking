package booking

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func (h *Handler) TrySeatLock(ctx context.Context, trainId, travelDate string, seatIds []string, holdToken string, ttl time.Duration) (bool, error) {
	pipe := h.Redis.TxPipeline()

	for _, seatId := range seatIds {
		key := fmt.Sprintf("seats:%s:%s:%s", trainId, travelDate, seatId)

		pipe.SetNX(ctx, key, holdToken, ttl)
	}
	// Execute all commands atomically
	cmds, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	// check if all the locks are acquired
	for _, cmd := range cmds {
		if !cmd.(*redis.BoolCmd).Val() {
			// One of the seats is already locked, release any acquired locks
			h.ReleaseLocks(ctx, trainId, travelDate, seatIds, holdToken)
			return false, nil
		}
	}

	return true, nil

}

func (h *Handler) ReleaseLocks(ctx context.Context, trainId, travelDate string, seatIds []string, holdToken string) error {
	pipe := h.Redis.TxPipeline()

	for _, seatId := range seatIds {
		key := fmt.Sprintf("seats:%s:%s:%s", trainId, travelDate, seatId)
		getCmd := pipe.Get(ctx, key)
		log.Print(getCmd)
	}

	cmds, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return err
	}

	// Second pipeline for deletions
	delPipe := h.Redis.Pipeline()
	deletedCount := 0

	for i, cmd := range cmds {
		if getCmd, ok := cmd.(*redis.StringCmd); ok {
			currentToken, err := getCmd.Result()
			if err == nil && currentToken == holdToken {
				seatId := seatIds[i]
				key := fmt.Sprintf("seats:%s:%s:%s", trainId, travelDate, seatId)
				delPipe.Del(ctx, key)
				deletedCount++
			}
		}
	}

	if deletedCount > 0 {
		_, err := delPipe.Exec(ctx)
		return err
	}

	return nil
}

func (h *Handler) CheckLock(ctx context.Context, trainId, travelDate string, seatId string, holdToken string) (bool, error) {
	key := fmt.Sprintf("seats:%s:%s:%s", trainId, travelDate, seatId)

	val, err := h.Redis.Get(ctx, key).Result()
	if err != nil {
		return false, err
	}
	if err == redis.Nil {
		// no seat lock for this seat
		return false, nil
	}

	return holdToken == val, nil
}
