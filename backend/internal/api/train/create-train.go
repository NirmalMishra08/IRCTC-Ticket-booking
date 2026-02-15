package train

import (
	"better-uptime/common/middleware"
	"better-uptime/common/util"
	db "better-uptime/internal/db/sqlc"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type CreateTrainRequest struct {
	TrainNumber int       `json:"train_number"`
	TrainName   string    `json:"train_name"`
	Source      string    `json:"source"`
	Destination string    `json:"destination"`
	Day         string    `json:"day"`
	ArrivalTime string `json:"arrival_time"`
}


type CreateTrainResponse struct {
	Train         db.Train         `json:"train"`
	TrainSchedule db.TrainSchedule `json:"train_schedule"`
}

func (h *Handler) CreateTrain(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	payload, err := middleware.GetFirebasePayloadFromContext(ctx)
	if err != nil {
		util.ErrorJson(w, util.ErrUnauthorized)
		return
	}

	role:= payload.Role
	if role != "ADMIN" {
		util.ErrorJson(w, util.ErrUnauthorized)
		return
	}

	var data CreateTrainRequest

	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		util.ErrorJson(w, util.ErrNotValidRequest)
		return
	}

	r.Body.Close()

	arrivalGoTime, err := time.Parse("15:04", data.ArrivalTime)
	if err != nil {
		util.ErrorJson(w, util.ErrNotValidRequest)
		return
	}

	const stopDuration = 5 * time.Minute
	departureTime := arrivalGoTime.Add(stopDuration)

	train, err := h.store.CreateTrain(ctx, db.CreateTrainParams{
		Trainnumber: int32(data.TrainNumber),
		Trainname:   data.TrainName,
		Source:      data.Source,
		Destination: data.Destination,
	})
	if err != nil {
		util.ErrorJson(w, err)
		return
	}

	trainSchedule, err := h.store.CreateTrainSchedule(ctx, db.CreateTrainScheduleParams{
		Trainid:       pgtype.Int4{Int32: train.ID, Valid: true},
		Day:           db.DayOfWeek(data.Day),
		Arrivaltime:   pgtype.Time{Microseconds: int64(arrivalGoTime.Hour()*3600000 + arrivalGoTime.Minute()*60000 + arrivalGoTime.Second()*1000000), Valid: true},
		Departuretime: pgtype.Time{Microseconds: int64(departureTime.Hour()*3600000 + departureTime.Minute()*60000 + departureTime.Second()*1000000), Valid: true},
	})
	if err != nil {
		util.ErrorJson(w, err)
		return
	}

	result := CreateTrainResponse{
		Train:         train,
		TrainSchedule: trainSchedule,
	}

	util.WriteJson(w, http.StatusCreated, result)

}
