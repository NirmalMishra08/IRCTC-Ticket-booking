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
	fmt.Println("🔥 webhook hit")
	sig := r.Header.Get("Stripe-Signature")
	fmt.Println(sig)
	fmt.Println("🔥 secret", h.config.STRIPE_WEBHOOK_SECRET)
	event, err := webhook.ConstructEventWithOptions(payload, sig, h.config.STRIPE_WEBHOOK_SECRET,
	webhook.ConstructEventOptions{
		IgnoreAPIVersionMismatch: true,
	})
	
	if err != nil {
		util.ErrorJson(w, fmt.Errorf("not able to create event", err),)
		return
	}

	fmt.Println("✅ event type:", event.Type)

	switch event.Type {
	case "checkout.session.completed":
		h.handlePaymentSuccess(event, ctx)

	case "checkout.session.expired":
		h.handlePaymentExpired(event, ctx)

	case "payment_intent.payment_failed":
		h.handlePaymentExpired(event, ctx)
	}

	util.WriteJson(w, http.StatusOK, nil)

}

func (h *Handler) handlePaymentSuccess(event stripe.Event, ctx context.Context) {
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

	bookingId, err := strconv.Atoi(bookingIDStr)
	if err != nil {
		logger.Error("invalid booking_id", err)
		return
	}

	if err := h.store.ExecTx(ctx, func(q *db.Queries) error {
		booking, err := q.GetBookingById(ctx, int32(bookingId))
		if err != nil {
			return err
		}

		if booking.Status == db.BookingStatusCONFIRMED {
			return nil
		}

		err = q.UpdateBookingStatus(ctx, db.UpdateBookingStatusParams{
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

		err = q.ConfirmSeat(ctx, util.ToPgInt4(int32(bookingId)))
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

		booking, err := q.GetBookingById(ctx, int32(bookingID))
		if err != nil {
			return err
		}

		if booking.Status == db.BookingStatusCONFIRMED {
			return nil
		}

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

}
