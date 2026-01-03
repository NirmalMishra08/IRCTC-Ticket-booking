package stripe

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/stripe/stripe-go/v84"
	"github.com/stripe/stripe-go/v84/checkout/session"
)

type PaymentRequest struct {
	UserUUID string `json:"user_uuid,omitempty"`
	Price    string `json:"price,omitempty"`
	PlanName string `json:"plan_name,omitempty"`
}

type APIResponse struct {
	status     int
	Message    string
	StatusCode int
	SessionURL StripeModel
}

type StripeModel struct {
	UserUUID   string `json:"user_uuid,omitempty"`
	SessionURL string `json:"session_url,omitempty"`
	SessionID  string `json:"session_id,omitempty"`
	PlanName   string `json:"plan_name,omitempty"`
	Price      string `json:"price,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
}

func ConvertTheAmount(price string, currency string) (int64, error) {
	var amount float64

	_, err := fmt.Sscanf(price, "%f", &amount)
	if err != nil {
		return 0, fmt.Errorf("invalid price format")
	}

	return int64(amount * 100), nil

}

func StripeSession(ctx context.Context, userUUID, price, planName, StripeKey string, bookingId int, holdToken string) (*APIResponse, error) {
	convertedAmount, err := ConvertTheAmount(price, "usd")
	if err != nil {
		return nil, fmt.Errorf("failed to convert the amount %w", err)
	}

	fmt.Print(convertedAmount)

	stripe.Key = StripeKey

	params := &stripe.CheckoutSessionParams{
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency: stripe.String("inr"),
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name: stripe.String(string(planName)),
				},
				UnitAmount: stripe.Int64(convertedAmount),
			},
			Quantity: stripe.Int64(1),
		}},
		AllowPromotionCodes: stripe.Bool(true),
		SuccessURL:          stripe.String("http://127.0.0.1:5500/Stripe-Payment-Go/payment-success.html"),
		CancelURL:           stripe.String("http://127.0.0.1:5500/Stripe-Payment-Go/payment-failed.html"),
		ExpiresAt:           stripe.Int64(time.Now().Add(30 * time.Minute).Unix()),
		Metadata: map[string]string{
			"booking_id": strconv.Itoa(bookingId),
			"hold_token": holdToken,
		},
	}
	params.AddMetadata("api_version", "2024-05-01")
	s, err := session.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create the session: %w", err)
	}

	stripeRecord := StripeModel{
		UserUUID:   userUUID,
		SessionURL: s.URL,
		SessionID:  s.ID,
		PlanName:   planName,
		Price:      price,
		CreatedAt:  time.Now().UTC().Format(time.RFC3339),
	}

	fmt.Println("Stripe Record is ------->>>>", stripeRecord)

	response := &APIResponse{
		status:     1,
		Message:    "Checkout session created successful",
		StatusCode: 200,
		SessionURL: StripeModel{
			SessionURL: s.URL,
			SessionID:  s.ID,
			UserUUID:   userUUID,
			PlanName:   planName,
			Price:      price,
			CreatedAt:  stripeRecord.CreatedAt,
		},
	}

	return response, nil

}
