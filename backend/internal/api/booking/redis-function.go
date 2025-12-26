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

func (h *Handler) CheckLockForSeat(ctx context.Context, trainId, travelDate string, seatId string, holdToken string) (bool, error) {
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

func (h *Handler) ValidateSeatLocks(ctx context.Context, trainId, travelDate string, seatIds []string, holdToken string) (bool, error) {
	pipe := h.Redis.Pipeline()
	cmds := make([]*redis.StringCmd, len(seatIds))

	for i, seatId := range seatIds {
		key := fmt.Sprintf("seats:%s:%s:%s", trainId, travelDate, seatId)
		cmds[i] = pipe.Get(ctx, key)
	}

	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return false, err
	}

	for _, cmd := range cmds {
		val, err := cmd.Result()
		if err != nil {
			return false, err
		}
		if val != holdToken {
			return false, err
		}
	}

	return true, nil
}

// extend the pipeline by ttl time
func (h *Handler) ExtendSeatLocksTTL(ctx context.Context, trainId, travelDate string, seatIds []string, ttl time.Duration) (bool, error) {
	pipe := h.Redis.Pipeline()
	cmds := make([]*redis.IntCmd, len(seatIds))

	for i, seatId := range seatIds {
		key := fmt.Sprintf("seats:%s:%s:%s", trainId, travelDate, seatId)

		cmds[i] = pipe.Exists(ctx, key)
	}

	results, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return false, err
	}

	for _, cmd := range results {
		if cmd.(*redis.IntCmd).Val() == 0 {
			return false, nil // Some lock expired/missing
		}
	}

	extendPipe := h.Redis.Pipeline()

	for _, seatId := range seatIds {
		key := fmt.Sprintf("seats:%s:%s:%s", trainId, travelDate, seatId)
		extendPipe.Expire(ctx, key, ttl)
	}

	_, err = extendPipe.Exec(ctx)

	return err == nil, err
}

func (h *Handler) AreSeatsTemporarilyLocked(ctx context.Context, trainId, travelDate string, seatIds []string) (map[string]bool, error) {
	pipe := h.Redis.Pipeline()
	cmds := make(map[string]*redis.StringCmd)

	for _, seatId := range seatIds {
		key := fmt.Sprintf("seats:%s:%s:%s", trainId, travelDate, seatId)
		cmds[seatId] = pipe.Get(ctx, key)
	}

	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, err
	}

	result := make(map[string]bool)

	for seatId, cmd := range cmds {
		_, err := cmd.Result()

		result[seatId] = (err != redis.Nil)
	}

	return result, nil

}
