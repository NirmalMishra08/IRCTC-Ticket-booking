package train

import (
	"better-uptime/common/logger"
	"better-uptime/common/util"
	db "better-uptime/internal/db/sqlc"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type GetAvailableSeatsRequest struct {
	TrainID    int32        `json:"train_id"`
	TravelDate string       `json:"travel_date"`
	Quota      db.SeatQuota `json:"quota,omitempty"`
}

func (h *Handler) GetAvailableSeats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var data GetAvailableSeatsRequest
	json.NewDecoder(r.Body).Decode(&data)

	travelDate, err := time.Parse("2006-01-02", data.TravelDate)
	if err != nil {
		util.ErrorJson(w, err)
		return
	}

	logger.Debug(travelDate.String())

	seats, err := h.store.GetAvailableSeats(ctx, db.GetAvailableSeatsParams{
		TrainID:     util.ToPgInt4(data.TrainID),
		JourneyDate: pgtype.Date{Valid: true, Time: travelDate},
		Quota:       data.Quota,
	})
	if err != nil {
		util.ErrorJson(w, err)
	}

	response := map[string]interface{}{
		"message": "All the trains with empty seats",
		"data":    seats,
	}

	util.WriteJson(w, http.StatusFound, response)
}
