// handler_simple.go
package train

import (
	"context"
	"fmt"
	"net/http"

	"better-uptime/common/util"
	db "better-uptime/internal/db/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

// CreateCoachesAndSeats - Simple version without complex transactions
func (h *Handler) CreateCoachesAndSeats(w http.ResponseWriter, r *http.Request) {
	var req CreateCoachSeatRequest
	if err := util.ReadJsonAndValidate(w, r, &req); err != nil {
		util.ErrorJson(w, err)
		return
	}

	ctx := r.Context()

	err := h.store.ExecTx(ctx, func(q *db.Queries) error {
		_, err := q.LockTrainForLayout(ctx, int32(req.TrainID))
		if err != nil {
			return fmt.Errorf("failed to lock train: %w", err)
		}

		existingCoaches, err := q.GetCoachesByTrain(ctx, pgtype.Int4{
			Int32: int32(req.TrainID),
			Valid: true,
		})
		if err != nil {
			return err
		}
		if len(existingCoaches) > 0 {
			return fmt.Errorf("layout already exists for this train")
		}

		for _, config := range req.Configurations {

			berthAllocation := getBerthAllocationSimple(config)

			for i := 0; i < config.NumberOfCoaches; i++ {

				// 🔢 Get next coach number safely
				nextCoachNumber, err := q.GetNextCoachNumber(ctx, pgtype.Int4{
					Int32: int32(req.TrainID),
					Valid: true,
				})
				if err != nil {
					return err
				}
				coach, err := q.CreateCoach(ctx, db.CreateCoachParams{
					Trainid:     pgtype.Int4{Int32: int32(req.TrainID), Valid: true},
					Coachtype:   db.CoachType(config.CoachType),
					Coachnumber: int32(nextCoachNumber),
				})
				if err != nil {
					return fmt.Errorf("failed to create coach: %w", err)
				}
				if err := createSeatsTx(ctx, q, coach.ID, berthAllocation); err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		util.ErrorJson(w, err)
		return
	}

	util.WriteJson(w, http.StatusCreated, map[string]string{
		"message": "seats created successfully",
	})
}

// getBerthAllocationSimple - Simple berth allocation
func getBerthAllocationSimple(config SeatConfiguration) []BerthAllocation {
	switch config.CoachType {
	case "1A", "2A", "3A":
		// AC coaches: alternate berths
		return []BerthAllocation{
			{BerthType: "DOWN", Count: config.SeatsPerCoach / 2},
			{BerthType: "UP", Count: config.SeatsPerCoach / 2},
		}
	case "SL":
		// Sleeper: mixed berths
		return []BerthAllocation{
			{BerthType: "DOWN", Count: config.SeatsPerCoach / 3},
			{BerthType: "MID", Count: config.SeatsPerCoach / 3},
			{BerthType: "UP", Count: config.SeatsPerCoach / 3},
		}
	case "GN":
		// General: all lower berth (for simplicity)
		return []BerthAllocation{
			{BerthType: "DOWN", Count: config.SeatsPerCoach},
		}
	default:
		// Default: all lower berth
		return []BerthAllocation{
			{BerthType: "DOWN", Count: config.SeatsPerCoach},
		}
	}
}

// createSeatsForCoachSimple - Simple seat creation
func (h *Handler) createSeatsForCoachSimple(ctx context.Context, coachID int32,
	berthAllocation []BerthAllocation) (int, error) {

	seatNumber := 1
	totalSeatsCreated := 0

	for _, allocation := range berthAllocation {
		for i := 0; i < allocation.Count; i++ {
			_, err := h.store.CreateSeat(ctx, db.CreateSeatParams{
				Coachid: pgtype.Int4{Int32: coachID, Valid: true},
				Seatno:  int32(seatNumber),
				Berth:   db.BerthType(allocation.BerthType),
			})
			if err != nil {
				return totalSeatsCreated, fmt.Errorf("failed to create seat %d: %v", seatNumber, err)
			}
			seatNumber++
			totalSeatsCreated++
		}
	}

	return totalSeatsCreated, nil
}

func createSeatsTx(
	ctx context.Context,
	q *db.Queries,
	coachID int32,
	allocations []BerthAllocation,
) error {

	seatNumber := int32(1)

	for _, alloc := range allocations {
		for i := 0; i < alloc.Count; i++ {
			_, err := q.CreateSeat(ctx, db.CreateSeatParams{
				Coachid: pgtype.Int4{Int32: coachID, Valid: true},
				Seatno:  seatNumber,
				Berth:   db.BerthType(alloc.BerthType),
			})
			if err != nil {
				return err
			}
			seatNumber++
		}
	}
	return nil
}
