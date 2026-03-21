package booking

import (
	"better-uptime/common/logger"
	"better-uptime/common/middleware"
	"better-uptime/common/stripe"
	"better-uptime/common/util"
	db "better-uptime/internal/db/sqlc"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type BookingRequest struct {
	JourneyId   int            `json:"journey_id,omitempty"`
	BookingType db.BookingType `json:"booking_type,omitempty"`
	SeatCount   int            `json:"seat_count,omitempty"`
	CoachType   db.CoachType   `json:"coach_type,omitempty"`
}

type PublishJob struct {
	BookingID string         `json:"booking_id"`
	UserID    string         `json:"user_id"`
	Data      BookingRequest `json:"data"`
}

func (h *Handler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	payload, err := middleware.GetFirebasePayloadFromContext(ctx)
	if err != nil {
		util.ErrorJson(w, util.ErrUnauthorized)
		return
	}

	var data BookingRequest
	err = util.ReadJsonAndValidate(w, r, &data)
	if err != nil {
		util.ErrorJson(w, util.ErrNotValidRequest)
		return
	}

	userId := payload.UserId

	if data.SeatCount <= 0 {
		util.ErrorJson(w, fmt.Errorf("not enought seats"))
		return
	}

	// rate limiter for upto 20 request in 10 minutes
	err = h.RateLimitUser(ctx, userId.String(), 10, 20)
	if err != nil {
		util.ErrorJson(w, util.ErrRateLimiting)
		return
	}

	train_journey, err := h.store.GetTrainJourneyById(ctx, int32(data.JourneyId))
	if err != nil {
		util.ErrorJson(w, errors.New("not able to get train journey details"))
	}

	if (train_journey.Status != db.NullJourneyStatus{JourneyStatus: "OPEN", Valid: true}) {
		util.ErrorJson(w, fmt.Errorf("Not opened for booking"))
		return
	}

	// check only if the booking type is tatkal
	if data.BookingType == "TATKAL" {
		tatkal_data, err := h.store.ValidateTatkalWindow(ctx, util.ToPgInt4(train_journey.TrainID.Int32))
		if err != nil {
			util.ErrorJson(w, fmt.Errorf("not able to get tatkal data"))
			return
		}

		if time.Now().Before(tatkal_data.TatkalStartTime.Time) || time.Now().After(tatkal_data.TatkalEndTime.Time) {
			util.ErrorJson(w, fmt.Errorf("current time is outside tatkal booking window"))
			return
		}
	}

	holdToken := string(uuid.New().String())

	if data.BookingType == db.BookingTypeTATKAL {

		booking, err := h.store.CreateBooking(ctx, db.CreateBookingParams{
			Userid:    pgtype.UUID{Bytes: userId, Valid: true},
			JourneyID: util.ToPgInt4(int32(data.JourneyId)),
			Holdtoken: pgtype.Text{String: holdToken, Valid: true},
		})

		job := PublishJob{
			BookingID: fmt.Sprintf("%d", booking.ID),
			UserID:    payload.UserId.String(),
			Data:      data,
		}

		value, err := json.Marshal(job)
		if err != nil {
			util.ErrorJson(w, fmt.Errorf("not able to parase the data"))
			return
		}

		partionKey := fmt.Sprintf("%d:%s:%s", data.JourneyId, data.CoachType, data.BookingType)
		// partition should be according to the journeyId coach type booking type for the ordering reason
		err = h.Kafka.Publish(ctx, "tatkal_booking", partionKey, value)
		if err != nil {
			util.ErrorJson(w, err)
			return
		}

		util.WriteJson(w, http.StatusAccepted, map[string]string{
			"message": "Tatkal booking request accepted. Processing in background.",
		})
		return
	} else {
		var bookingId int

		err = h.store.ExecTx(ctx, func(q *db.Queries) error {
			booking, err := q.CreateBooking(ctx, db.CreateBookingParams{
				Userid:    pgtype.UUID{Bytes: userId, Valid: true},
				JourneyID: util.ToPgInt4(int32(data.JourneyId)),
				Holdtoken: pgtype.Text{String: holdToken, Valid: true},
			})
			if err != nil {
				return fmt.Errorf("not able to book seats: %w", err)
			}

			bookingId = int(booking.ID)

			seatIDs, err := q.LockAvailableSeats(ctx, db.LockAvailableSeatsParams{
				JourneyID: int32(data.JourneyId),
				CoachType: data.CoachType,
				Quota:     db.SeatQuota(data.BookingType),
				SeatLimit: int32(data.SeatCount),
			})
			if err != nil {
				return fmt.Errorf("not able to lock seats: %w", err)
			}

			if len(seatIDs) < data.SeatCount {
				err := q.UpdateBookingStatus(ctx, db.UpdateBookingStatusParams{
					ID:     booking.ID,
					Status: db.BookingStatusWAITLIST,
				})
				if err != nil {
					return err
				}

				wlNumber, err := q.GetNextWaitlistNumber(ctx, util.ToPgInt4(int32(data.JourneyId)))
				if err != nil {
					return err
				}

				err = q.InsertWaitlist(ctx, db.InsertWaitlistParams{
					JourneyID:      util.ToPgInt4(int32(data.JourneyId)),
					Bookingid:      util.ToPgInt4(int32(bookingId)),
					WaitlistNumber: int32(wlNumber),
				})
				if err != nil {
					return err
				}

				return nil
			} else {
				for _, seatID := range seatIDs {
					err := q.HoldSeat(ctx, db.HoldSeatParams{
						JourneyID: int32(data.JourneyId),
						SeatID:    seatID,
						BookingID: util.ToPgInt4(booking.ID),
					})
					if err != nil {
						return fmt.Errorf("failed to hold seat %d: %w", seatID, err)
					}
				}

				for _, seatID := range seatIDs {
					_, err := q.CreateBookingItem(ctx, db.CreateBookingItemParams{
						Bookingid: util.ToPgInt4(booking.ID),
						Seatid:    util.ToPgInt4(seatID),
					})
					if err != nil {
						return fmt.Errorf("failed to create booking item: %w", err)
					}
				}

			}

			return nil

		})

		if err != nil {
			util.ErrorJson(w, err)
			return
		}

		amount := CalculateFare(int32(data.SeatCount), data.CoachType, data.BookingType)

		paymentIntent, err := stripe.StripeSession(ctx, userId.String(), fmt.Sprintf("%d", amount), "booking_train", h.config.STRIPE_SECRET_KEY, int(bookingId), holdToken)
		if err != nil {
			updateErr := h.store.UpdateBookingStatus(ctx, db.UpdateBookingStatusParams{
				ID:     int32(bookingId),
				Status: db.BookingStatusEXPIRED,
			})

			if updateErr != nil {
				logger.Error("failed to expire booking %d: %v", bookingId, updateErr)
			}

			_ = h.store.ReleaseSeatsByBooking(ctx, util.ToPgInt4(int32(bookingId)))

			util.ErrorJson(w, errors.New("not able to create booking intent"))
			return
		}

		_, err = h.store.CreatePayment(ctx, db.CreatePaymentParams{
			Bookingid:     util.ToPgInt4(int32(bookingId)),
			Amount:        float64(amount),
			Transactionid: paymentIntent.SessionURL.SessionID,
		})

		//Send user notification

		if err != nil {
			util.ErrorJson(w, fmt.Errorf("not able to create payment"))
		}

		booking, err := h.store.GetBookingById(ctx, int32(bookingId))
		if err != nil {
			util.ErrorJson(w, err)
			return
		}

		if booking.Status == db.BookingStatusWAITLIST {
			util.WriteJson(w, http.StatusOK, map[string]interface{}{
				"bookingId": bookingId,
				"status":    "WAITLIST",
				"message":   "Seats not available. You are on waitlist.",
			})
			return
		}

		response := map[string]interface{}{
			"bookingId":  bookingId,
			"sessionUrl": paymentIntent.SessionURL,
			"expires_in": 600,
		}

		util.WriteJson(w, http.StatusOK, response)

	}
}

func CalculateFare(seatIds int32, coachType db.CoachType, bookingType db.BookingType) int32 {

	return seatIds * 100
}
