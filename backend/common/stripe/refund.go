package stripe

import (
	"context"
	"errors"
	"fmt"

	"github.com/stripe/stripe-go/v84"
	"github.com/stripe/stripe-go/v84/refund"
)

type apiResponse struct {
	Status  string
	Message string
}

func RefundSession(ctx context.Context, userId, amount, holdToken, StripeKey string) (*apiResponse, error) {
	price, err := ConvertTheAmount(amount, "usd")
	if err != nil {
		return nil, fmt.Errorf("failed to convert the amount %w", err)
	}
	stripe.Key = StripeKey
	params := &stripe.RefundParams{
		Amount: stripe.Int64(price),
		Charge: stripe.String(holdToken),
	}
	params.SetIdempotencyKey(fmt.Sprintf("refund_%s_%s", userId, holdToken))

	result, err := refund.New(params)
	if err != nil {
		// Start checking specific Stripe Error types
		var stripeErr *stripe.Error
		if errors.As(err, &stripeErr) {
			switch stripeErr.Code {
			case stripe.ErrorCodeChargeAlreadyRefunded:
				return nil, fmt.Errorf("this payment has already been refunded")
			case stripe.ErrorCodeAmountTooLarge:
				return nil, fmt.Errorf("refund amount exceeds the original payment")
			case stripe.ErrorCodeBalanceInsufficient:
				return nil, fmt.Errorf("your Stripe balance is too low to process this refund")
			}
		}
		return nil, fmt.Errorf("stripe error: %v", err.Error())
	}
	
	return &apiResponse{
		Status:  "success",
		Message: fmt.Sprintf("Refund of %v processed for ID %s", price, result.ID),
	}, nil

}
