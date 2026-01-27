package tatkal

import (
	"better-uptime/common/middleware"
	"better-uptime/common/util"
	"encoding/json"
	"errors"
	"net/http"
)

type TatkalRequest struct {
	trainId string
	travelDate string
	coach_type string
	passengers []PassenngerDetails
}

type PassenngerDetails struct {
	name string
	age int
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

	// get tatkalDate and time

}
