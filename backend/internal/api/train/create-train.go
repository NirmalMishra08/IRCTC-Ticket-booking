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
	ArrivalTime time.Time `json:"arrival_time"`
}


type CreateTrainResponse struct {
	Train         db.Train         `json:"train"`
	TrainSchedule db.Trainschedule `json:"train_schedule"`
}

func (h *Handler) CreateTrain(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	_, err := middleware.GetFirebasePayloadFromContext(ctx)
	if err != nil {
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

	arrivalGoTime := time.Time(data.ArrivalTime)

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
		Arrivaltime:   arrivalGoTime,
		Departuretime: departureTime,
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
