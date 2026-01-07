package cancellation

import (
	"better-uptime/common/middleware"
	"better-uptime/common/stripe"
	"better-uptime/common/util"
	db "better-uptime/internal/db/sqlc"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type RefundRequest struct {
	TrainId    string `json:"train_id" validate:"required"`
	TravelDate string `json:"travel_date" validate:"required`
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

	TrainId, err := strconv.Atoi(data.TrainId)
	if err != nil {
		util.ErrorJson(w, errors.New("not able to parse string to integer"))
		return
	}

	travelDate, err := time.Parse("2006-01-02", data.TravelDate)
	if err != nil {
		util.ErrorJson(w, util.ErrNotValidRequest)
		return
	}

	trainWithAmount, err := h.store.GetPaymentAndTrain(ctx, db.GetPaymentAndTrainParams{
		Userid:     pgtype.UUID{Bytes: userId, Valid: true},
		Trainid:    util.ToPgInt4(int32(TrainId)),
		Traveldate: pgtype.Date{Time: travelDate, Valid: true},
	})
	if err != nil {
		util.ErrorJson(w, err)
		return
	}

	var amount float64
	amount = trainWithAmount.Amount
	amountStr := fmt.Sprintf("%.2f", amount)

	apiResponse, err := stripe.RefundSession(ctx, userId.String(), amountStr, trainWithAmount.Holdtoken.String, h.config.REDIS_DB_URL)
	if err != nil || apiResponse == nil {
		util.ErrorJson(w, err)
	}

	if apiResponse.Status == "success" {

	}

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
