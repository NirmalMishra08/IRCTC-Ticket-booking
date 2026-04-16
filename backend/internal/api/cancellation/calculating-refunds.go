package cancellation

import (
	"better-uptime/common/middleware"
	"better-uptime/common/stripe"
	"better-uptime/common/util"
	db "better-uptime/internal/db/sqlc"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type RefundRequest struct {
	JourneyID string `json:"journey_id" validate:"required"`
}

func (h *Handler) CalculatingRefundAmount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	payload, err := middleware.GetFirebasePayloadFromContext(ctx)
	if err != nil {
		util.ErrorJson(w, util.ErrUnauthorized)
	}

	userId := payload.UserId

	var data RefundRequest
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		util.ErrorJson(w, util.ErrNotValidRequest)
		return
	}

	JourneyId, err := strconv.Atoi(data.JourneyID)
	if err != nil {
		util.ErrorJson(w, fmt.Errorf("not able to parese the journey id"))
		return
	}

	fmt.Println("hello")

	trainWithAmount, err := h.store.GetPaymentAndTrain(ctx, db.GetPaymentAndTrainParams{
		Userid:    pgtype.UUID{Bytes: userId, Valid: true},
		JourneyID: util.ToPgInt4(int32(JourneyId)),
	})
	if err != nil {
		util.ErrorJson(w, err)
		return
	}

	if trainWithAmount.Status_2.PaymentStatus != db.PaymentStatusSUCCESS {
		util.ErrorJson(w, fmt.Errorf("payment was not successfull during booking"))
		return
	}

	if trainWithAmount.Status != db.BookingStatusCONFIRMED {
		util.ErrorJson(w, fmt.Errorf("no seats were confirmed"))
		return
	}

	fmt.Println("hello")
	// calculate the refund
	var amount float64
	amount = trainWithAmount.Amount
	amountStr := fmt.Sprintf("%.2f", amount)

	bookingId := trainWithAmount.Bookingid

	var releasedSeats int64

	apiResponse, err := stripe.RefundSession(ctx, userId.String(), amountStr, trainWithAmount.Holdtoken.String, h.config.STRIPE_SECRET_KEY)
	if err != nil || apiResponse == nil {
		util.ErrorJson(w, err)
	}
	//begin db transaction
	err = h.store.ExecTx(ctx, func(q *db.Queries) error {
		releasedSeats, err = q.CountSeatsByBooking(ctx, util.ToPgInt4(bookingId.Int32))
		if err != nil {
			return fmt.Errorf("failed to count seats: %w", err)
		}
		// update booking -> cancelled
		err := q.UpdateBookingStatus(ctx, db.UpdateBookingStatusParams{
			ID:     bookingId.Int32,
			Status: db.BookingStatusCANCELLED,
		})

		if err != nil {
			return fmt.Errorf("not able to update the db: %w", err)
		}

		if err := q.ReleaseSeatsByBooking(ctx, util.ToPgInt4(bookingId.Int32)); err != nil {
			return fmt.Errorf("failed to release seats: %w", err)
		}

		// update payment

		_, err = q.CreateRefund(ctx, db.CreateRefundParams{
			Userid:    pgtype.UUID{Bytes: userId, Valid: true},
			Bookingid: util.ToPgInt4(bookingId.Int32),
			Amount:    int32(amount),
			Status:    db.RefundStatusSUCCESS,
		})
		if err != nil {
			return fmt.Errorf("not able to create the refund: %w", err)
		}

		err = q.DeleteBookingItem(ctx, util.ToPgInt4(bookingId.Int32))
		if err != nil {
			return fmt.Errorf("not able to delete  the booking items table: %w", err)
		}
		return nil
	})

	if err != nil {
		util.ErrorJson(w, err)
		return
	}

	coachtype, err := h.store.GetCoachTypeByJourneyId(ctx, int32(JourneyId))
	if err != nil {
		util.ErrorJson(w, fmt.Errorf("not able to get the coach type"))
	}

	key := fmt.Sprintf("%d:%s", JourneyId, coachtype)

	message := map[string]interface{}{
		"journey_id":     JourneyId,
		"coach_type":     coachtype,
		"released_seats": releasedSeats,
		"timestamp":      time.Now().Unix(),
	}

	value, _ := json.Marshal(message)
	h.Kafka.Publish(ctx, "seat_released", key, value)

	response := map[string]interface{}{
		"message":         "refund in process",
		"response_stripe": apiResponse,
	}

	util.WriteJson(w, http.StatusOK, response)

}

func CalculatateRefund(cancellationtime, arrivalTime time.Time, amount int) float64 {
	timeDiff := arrivalTime.Sub(cancellationtime)

	switch {
	case timeDiff <= 0:
		return 0
	case timeDiff >= 24*time.Hour:
		return float64(amount)
	case timeDiff < 24*time.Hour:
		return float64(amount / 2)
	}

	return 0
}
