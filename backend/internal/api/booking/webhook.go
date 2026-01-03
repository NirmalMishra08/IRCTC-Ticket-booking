package booking

import (
	"better-uptime/common/logger"
	"better-uptime/common/util"
	db "better-uptime/internal/db/sqlc"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/stripe/stripe-go/v84"
	"github.com/stripe/stripe-go/v84/webhook"
)

type BookingLockContext struct {
	TrainID    int32
	TravelDate time.Time
	SeatID     int32
	HoldToken  string
}

func (h *Handler) StripeWebhook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		util.ErrorJson(w, errors.New("bad request"))
		return
	}

	sig := r.Header.Get("Stripe-Signature")
	event, err := webhook.ConstructEvent(payload, sig, h.config.STRIPE_WEBHOOK_SECRET)
	if err != nil {
		util.ErrorJson(w, fmt.Errorf("not able to create event"))
		return
	}

	switch event.Type {
	case "checkout.session.completed":
		h.handlePaymentSuccess(event, ctx)

	case "checkout.session.expired":
		h.handlePaymentExpired(event, ctx)
	}

	util.WriteJson(w, http.StatusOK, nil)

}

func (h *Handler) handlePaymentSuccess(event stripe.Event, ctx context.Context) {
	var session stripe.CheckoutSession

	json.Unmarshal(event.Data.Raw, &session)

	bookingId, _ := strconv.Atoi(session.Metadata["booking_id"])

	if err := h.store.ExecTx(ctx, func(q *db.Queries) error {
		err := q.UpdateBookingStatus(ctx, db.UpdateBookingStatusParams{
			ID:     int32(bookingId),
			Status: db.BookingStatusCONFIRMED,
		})
		if err != nil {
			return err
		}

		err = q.UpdateBookingItemStatus(ctx, db.UpdateBookingItemStatusParams{
			Bookingid:     util.ToPgInt4(int32(bookingId)),
			Bookingstatus: db.BookingStatusCONFIRMED,
		})
		if err != nil {
			return err
		}

		err = q.UpdatePaymentStatus(ctx, db.UpdatePaymentStatusParams{
			Bookingid: util.ToPgInt4(int32(bookingId)),
			Status:    db.NullPaymentStatus{PaymentStatus: db.PaymentStatusSUCCESS, Valid: true},
		})
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		logger.Error("error occurred while updating booking status", err)
	}

}

func (h *Handler) handlePaymentExpired(
	event stripe.Event,
	ctx context.Context,
) {
	var session stripe.CheckoutSession
	if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
		logger.Error("failed to parse checkout session", err)
		return
	}

	bookingIDStr := session.Metadata["booking_id"]
	if bookingIDStr == "" {
		logger.Error("missing booking_id in metadata")
		return
	}

	bookingID, err := strconv.Atoi(bookingIDStr)
	if err != nil {
		logger.Error("invalid booking_id", err)
		return
	}

	err = h.store.ExecTx(ctx, func(q *db.Queries) error {

		// 1. Booking → EXPIRED
		if err := q.UpdateBookingStatus(ctx, db.UpdateBookingStatusParams{
			ID:     int32(bookingID),
			Status: db.BookingStatusEXPIRED,
		}); err != nil {
			return err
		}

		// 2. BookingItems → EXPIRED
		if err := q.UpdateBookingItemStatus(ctx, db.UpdateBookingItemStatusParams{
			Bookingid:     util.ToPgInt4(int32(bookingID)),
			Bookingstatus: db.BookingStatusEXPIRED,
		}); err != nil {
			return err
		}

		// 3. Payment → FAILED
		if err := q.UpdatePaymentStatus(ctx, db.UpdatePaymentStatusParams{
			Bookingid: util.ToPgInt4(int32(bookingID)),
			Status:    db.NullPaymentStatus{PaymentStatus: db.PaymentStatusFAILED, Valid: true},
		}); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		logger.Error("failed to expire booking", err)
		return
	}

	// 4. Release seat locks (outside DB)
	holdToken := session.Metadata["hold_token"]
	if holdToken != "" {
		h.ReleaseLocksByBooking(ctx, int32(bookingID))
	}
}

func (h *Handler) ReleaseLocksByBooking(
	ctx context.Context,
	bookingID int32,
) error {

	rows, err := h.store.GetBookingLockContext(ctx, bookingID)
	if err != nil {
		return err
	}

	if len(rows) == 0 {
		return nil
	}

	trainId := strconv.Itoa(int(rows[0].Trainid.Int32))
	travelDate := rows[0].Traveldate.Time.Format("2006-01-02")
	holdToken := rows[0].Holdtoken.String

	var seatIds []string
	for _, row := range rows {
		seatIds = append(seatIds, strconv.Itoa(int(row.Seatid.Int32)))
	}

	return h.ReleaseLocks(
		ctx,
		trainId,
		travelDate,
		seatIds,
		holdToken,
	)
}
