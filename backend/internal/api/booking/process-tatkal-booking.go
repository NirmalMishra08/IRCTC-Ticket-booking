package booking

import (
	"better-uptime/common/util"
	db "better-uptime/internal/db/sqlc"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"
)

func (h *Handler) ProcessTatkalBooking(ctx context.Context,
	data BookingRequest,
	userId string,
	bookingId string) error {
	// redis for handling massive crowd
	LockKey := fmt.Sprintf("tatkal:%d:%s", data.JourneyId, data.CoachType)
	ok, err := h.Redis.SetNX(ctx, LockKey, userId, 5*time.Second).Result()
	if err != nil {
		return fmt.Errorf("failed to acquire lock: %w", err)
	}
	if !ok {
		return errors.New("booking in progress, please retry")
	}
	defer h.Redis.Del(ctx, LockKey)

	return h.store.ExecTx(ctx, func(q *db.Queries) error {
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

		bookingIdInt, err := strconv.Atoi(bookingId)
		if err != nil {
			return err
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
}
