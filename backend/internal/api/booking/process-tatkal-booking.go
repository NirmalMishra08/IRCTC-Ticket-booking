package booking

import (
	"better-uptime/common/util"
	db "better-uptime/internal/db/sqlc"
	"context"
	"errors"
	"fmt"
	"strconv"
)

func (h *Handler) ProcessTatkalBooking(ctx context.Context,
	data BookingRequest,
	userId string,
	bookingId string) error {

	bookingIdInt, err := strconv.Atoi(bookingId)
	if err != nil {
		return err
	}

	booking, err := h.store.GetBookingById(ctx, int32(bookingIdInt))
	if err == nil && (booking.Status == db.BookingStatusCONFIRMED || booking.Status == db.BookingStatusEXPIRED || booking.Status == db.BookingStatusPENDING ) {
		return nil 
	}

	// redis for handling massive crowd
	// LockKey := fmt.Sprintf("tatkal:%d:%s", data.JourneyId, data.CoachType)
	// ok, err := h.Redis.SetNX(ctx, LockKey, userId, 5*time.Second).Result()
	// if err != nil {
	// 	return fmt.Errorf("failed to acquire lock: %w", err)
	// }
	// if !ok {
	// 	return errors.New("booking in progress, please retry")
	// }
	// defer h.Redis.Del(ctx, LockKey)

	// removed the setNx as it was unable to handle scalability 2ms for locks * 1 million user = 2000s = 33 minutes for a seat

	ok, err := h.reserveSeats(ctx, data)
	if err != nil {
		return err // retryable
	}
	if !ok {
		return errors.New("not enough tatkal seats available") // non-retryable
	}

	err = h.store.ExecTx(ctx, func(q *db.Queries) error {
		journey, err := q.GetTrainJourneyById(ctx, int32(data.JourneyId))
		if err != nil {
			return err
		}

		if !journey.Status.Valid || journey.Status.JourneyStatus != db.JourneyStatusOPEN {
			return errors.New("journey not open for booking")
		}

		seatIDs, err := q.LockAvailableSeats(ctx, db.LockAvailableSeatsParams{
			JourneyID: int32(data.JourneyId),
			CoachType: data.CoachType,
			Quota:     db.SeatQuotaTATKAL,
			SeatLimit: int32(data.SeatCount),
		})
		if err != nil {
			return fmt.Errorf("not able to lock seats: %w", err)
		}

		if len(seatIDs) < data.SeatCount {
			return errors.New("not enough tatkal seats available")
		}

		for _, seatID := range seatIDs {
			err := q.HoldSeat(ctx, db.HoldSeatParams{
				JourneyID: int32(data.JourneyId),
				SeatID:    seatID,
				BookingID: util.ToPgInt4(int32(bookingIdInt)),
			})
			if err != nil {
				return err
			}
		}

		for _, seatID := range seatIDs {
			_, err := q.CreateBookingItem(ctx, db.CreateBookingItemParams{
				Bookingid: util.ToPgInt4(int32(bookingIdInt)),
				Seatid:    util.ToPgInt4(seatID),
			})
			if err != nil {
				return fmt.Errorf("failed to create booking item: %w", err)
			}
		}

		return nil

	})

	if err != nil {
		// compensate
		key := fmt.Sprintf("tatkal:available:%d:%s", data.JourneyId, data.CoachType)
		_ = h.Redis.IncrBy(ctx, key, int64(data.SeatCount))
		return err
	}

	return nil
}

var reserveSeatsLua = `
-- reserve_seats.lua
-- KEYS[1] = seat key
-- ARGV[1] = seats requested
local avail = tonumber(redis.call('GET', KEYS[1]) or '0')
local req = tonumber(ARGV[1])
if avail >= req then
  redis.call('DECRBY', KEYS[1], req)
  return 1
end
return 0
`

func (h *Handler) reserveSeats(ctx context.Context, data BookingRequest) (bool, error) {
	key := fmt.Sprintf("tatkal:available:%d:%s", data.JourneyId, data.CoachType)
	res, err := h.Redis.Eval(ctx, reserveSeatsLua, []string{key}, data.SeatCount).Result()
	if err != nil {
		return false, err
	}
	return res.(int64) == 1, nil
}
