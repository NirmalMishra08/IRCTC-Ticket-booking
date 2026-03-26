package booking

import (
	db "better-uptime/internal/db/sqlc"
	"context"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
)

func (h *Handler) PromoteWaitlist(ctx context.Context, JourneyId string, CoachType db.CoachType) error {
	for {
		err := h.store.ExecTx(ctx, func(q *db.Queries) error {
			JourneyID, err := strconv.Atoi(JourneyId)
			train_journey, err := q.GetTrainJourneyById(ctx, int32(JourneyID))
			if err != nil {
				return err
			}

			availability, err := q.GetAvailableSeats(ctx, db.GetAvailableSeatsParams{
				TrainID:     train_journey.TrainID,
				JourneyDate: train_journey.JourneyDate,
				Quota:       db.SeatQuotaNORMAL,
			})

			if err != nil {
				return err
			}
			var available int32
			for _, a := range availability {
				if a.CoachType == CoachType {
					available = int32(a.AvailableSeats)
					break
				}
			}

			if available == 0 {
				return nil
			}

			wl, err := q.GetNextWaitlist(ctx, pgtype.Int4{Int32: int32(JourneyID), Valid:  true})
			if err != nil {
				return nil
			}

			required := wl.WaitlistNumber

			if available < required {
				return nil
			}

			err = q.ConfirmSeat(ctx, wl.Bookingid)
			if err != nil {
				return err
			}

			// update booking
			err = q.UpdateBookingStatus(ctx, db.UpdateBookingStatusParams{
				ID:     wl.Bookingid.Int32,
				Status: db.BookingStatusCONFIRMED,
			})
			if err != nil {
				return err
			}

			// update waitlist
			err = q.UpdateWaitlistStatus(ctx, db.UpdateWaitlistStatusParams{
				ID:     wl.ID,
				Status: db.WaitingStatusCONFIRMED,
			})
			if err != nil {
				return err
			}

			return nil

		})

		if err != nil {
			return err
		}
	}
}
