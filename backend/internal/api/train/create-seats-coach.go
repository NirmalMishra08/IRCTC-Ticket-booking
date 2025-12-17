// handler_simple.go
package train

import (
	"context"
	"fmt"
	"net/http"

	"better-uptime/common/logger"
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

	// Verify train exists
	train, err := h.store.GetTrainById(ctx, int32(req.TrainID))
	if err != nil {
		util.ErrorJson(w, fmt.Errorf("train with ID %d not found", req.TrainID))
		return
	}

	// Check if coaches already exist
	existingCoaches, _ := h.store.GetCoachesByTrain(ctx, pgtype.Int4{Int32: int32(req.TrainID),Valid: true})
	if len(existingCoaches) > 0 {
		// Optional: Delete existing coaches first
		// Or return error asking to delete first
		util.ErrorJson(w, fmt.Errorf("train already has %d coaches. Delete them first", len(existingCoaches)))
		return
	}

	// Track statistics
	totalCoaches := 0
	totalSeats := 0
	coachBreakdown := []CoachTypeSummary{}

	// Create coaches and seats
	for _, config := range req.Configurations {
		// Get berth allocation for this coach type
		berthAllocation := getBerthAllocationSimple(config)

		coachSummary := CoachTypeSummary{
			CoachType:       string(config.CoachType),
			NumberOfCoaches: config.NumberOfCoaches,
			SeatsPerCoach:   config.SeatsPerCoach,
			TotalSeats:      config.NumberOfCoaches * config.SeatsPerCoach,
		}
		coachBreakdown = append(coachBreakdown, coachSummary)

		for coachNum := 1; coachNum <= config.NumberOfCoaches; coachNum++ {
			// Create coach
			coach, err := h.store.CreateCoach(ctx, db.CreateCoachParams{
				Trainid:     pgtype.Int4{Int32: int32(req.TrainID), Valid: true},
				Coachtype:   db.CoachType(config.CoachType),
				Coachnumber: int32(coachNum),
			})
			if err != nil {
				logger.Error("Failed to create coach %d: %v", coachNum, err)
				util.ErrorJson(w, fmt.Errorf("failed to create coach %d: %v", coachNum, err))
				return
			}
			totalCoaches++

			// Create seats for this coach
			seatsCreated, err := h.createSeatsForCoachSimple(ctx, coach.ID, berthAllocation)
			if err != nil {
				logger.Error("Failed to create seats for coach %d: %v", coach.ID, err)
				util.ErrorJson(w, fmt.Errorf("failed to create seats: %v", err))
				return
			}
			totalSeats += seatsCreated

			logger.Info("Created coach %d (type: %s, number: %d) with %d seats",
				coach.ID, config.CoachType, coachNum, seatsCreated)
		}
	}

	// Prepare response
	response := CoachSeatResponse{
		TrainID:        req.TrainID,
		TrainNumber:    int(train.Trainnumber),
		TrainName:      train.Trainname,
		TotalCoaches:   totalCoaches,
		TotalSeats:     totalSeats,
		CoachBreakdown: coachBreakdown,
		Message:        fmt.Sprintf("Successfully created %d coaches with %d seats", totalCoaches, totalSeats),
	}

	util.WriteJson(w, http.StatusCreated, response)
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
