package train

import (
	"better-uptime/common/util"
	db "better-uptime/internal/db/sqlc"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type JourneyHandler struct {
	Queries *db.Queries
}

type CreateJourneyRequest struct {
	TrainID     int    `json:"train_id"`
	JourneyDate string `json:"journey_date"` // YYYY-MM-DD
	ScheduleID  int    `json:"schedule_id"`
	Status      string `json:"status"` // OPEN / CHARTED / CANCELLED
}

func (h *Handler) CreateJourney(w http.ResponseWriter, r *http.Request) {
	var req CreateJourneyRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	journeyDate, err := time.Parse("2006-01-02", req.JourneyDate)
	if err != nil {
		http.Error(w, "invalid date format", http.StatusBadRequest)
		return
	}

	if req.Status != "OPEN" && req.Status != "CHARTED" && req.Status != "CANCELLED" {
		http.Error(w, "invalid status", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	var result db.TrainJourney

	err = h.store.ExecTx(ctx, func(q *db.Queries) error {

		// ✅ Step 1: Create Journey
		journey, err := q.CreateTrainJourney(ctx, db.CreateTrainJourneyParams{
			TrainID:     util.ToPgInt4(int32(req.TrainID)),
			JourneyDate: pgtype.Date{Time: journeyDate, Valid: true},
			ScheduleID:  util.ToPgInt4(int32(req.ScheduleID)),
			Status:      db.NullJourneyStatus{JourneyStatus: db.JourneyStatus(req.Status), Valid: true},
		})
		if err != nil {
			return err
		}

		// ✅ Step 2: Initialize Seat Inventory
		err = q.InitializeSeatInventory(ctx, db.InitializeSeatInventoryParams{
			JourneyID: journey.ID,
			Trainid:   util.ToPgInt4(int32(req.TrainID)),
		})
		if err != nil {
			return err
		}

		result = journey
		return nil
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}