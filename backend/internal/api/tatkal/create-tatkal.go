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
)

type TatkalRequest struct {
	trainId    string              `json:"train_id,omitempty"`
	travelDate string              `json:"travel_date,omitempty"`
	coach_type string              `json:"coach_type,omitempty"`
	passengers []PassengerDetails `json:"passengers,omitempty"`
}

type PassengerDetails struct {
	name   string `json:"name,omitempty"`
	age    int    `json:"age,omitempty"`
	gender string `json:"gender,omitempty"`
}

func (h *Handler) CreateTatkalBooking(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	payload, err := middleware.GetFirebasePayloadFromContext(ctx)
	if err != nil {
		util.ErrorJson(w, util.ErrUnauthorized)
		return
	}

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

	if tatkalData.TatkalStartTime.Valid && tatkalData.TatkalStartTime.Time.After(time.Now()) {
		util.ErrorJson(w, errors.New("tatkal booking window is not open"))
		return
	}
	if tatkalData.TatkalEndTime.Valid && tatkalData.TatkalEndTime.Time.Before(time.Now()) {
		util.ErrorJson(w, errors.New("tatkal booking window is not open"))
		return
	}

	err = h.RateLimiter.RateLimitUser(ctx, userId.String(), 5, 5)
	if err != nil {
		util.ErrorJson(w, fmt.Errorf("rate limit exceed for this user %d", userId))
		return
	}

	err = h.ProcessTatkalUser(ctx, data, userId.String())
	if err != nil {
		util.ErrorJson(w, fmt.Errorf("not able to add user into queue"))
	}

	response := map[string]interface{}{
		"message": "user in queue",
	}

	util.WriteJson(w, http.StatusAccepted, response)
}

func (h *Handler) ProcessTatkalUser(ctx context.Context, data TatkalRequest, userId string) error {
	msg := map[string]interface{}{
		"user_id":     userId,
		"train_id":    data.trainId,
		"travel_date": data.travelDate,
		"coach_type":  data.coach_type,
		"passengers":  data.passengers,
	}

	bytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return h.Kafka.Publish(ctx, "tatkal-requests", userId, bytes)
}
