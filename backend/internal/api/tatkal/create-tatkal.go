package tatkal

import (
	"better-uptime/common/middleware"
	"better-uptime/common/util"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type TatkalRequest struct {
	trainId    string
	travelDate string
	coach_type string
	passengers []PassenngerDetails
}

type PassenngerDetails struct {
	name   string
	age    int
	gender string
}

func (h *Handler) CreateTatkalBooking(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	payload, err := middleware.GetFirebasePayloadFromContext(ctx)

	var data TatkalRequest
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		util.ErrorJson(w, errors.New("not able to parse reques"))
	}

	userId := payload.UserId

	TrainId, err := strconv.Atoi(data.trainId)
	if err != nil {
		util.ErrorJson(w, fmt.Errorf("not able to parse string to integer"))
	}

	// get tatkalDate and time
	tatkalData, err := h.store.GetTatkaData(ctx, util.ToPgInt4(int32(TrainId)))
	if err != nil {
		util.ErrorJson(w, fmt.Errorf("not able to get tatkalData"))
		return
	}

	if tatkalData.TatkalStartTime.Valid && tatkalData.TatkalStartTime.Time.Before(time.Now()) {
		util.ErrorJson(w, errors.New("tatkal booking window is not open"))
		return
	}
	if tatkalData.TatkalEndTime.Valid && tatkalData.TatkalEndTime.Time.After(time.Now()) {
		util.ErrorJson(w, errors.New("tatkal booking window is not open"))
		return
	}

	err = h.RateLimiter.RateLimitUser(ctx, userId.String(), 5, 5)
	if err != nil {
		util.ErrorJson(w, fmt.Errorf("rate limit exceed for this user %d", userId))
		return
	}

	err = h.EnqueueTatkalUser(ctx, TrainId, data.travelDate, data.coach_type, userId.String())
	if err != nil {
		util.ErrorJson(w, fmt.Errorf("not able to add user into queue"))
	}

}

func (h *Handler) EnqueueTatkalUser(ctx context.Context, trainId int, date string, coachType string, userId string) error {
	key := fmt.Sprintf("tatkal:queue:%d:%s:%s", trainId, date, coachType)
	score := float64(time.Now().UnixNano())

	return h.Redis.ZAdd(ctx, key, redis.Z{
		Score:  score,
		Member: userId,
	}).Err()
}
