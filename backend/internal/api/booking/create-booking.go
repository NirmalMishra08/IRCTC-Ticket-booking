package booking

import (
	"better-uptime/common/middleware"
	"better-uptime/common/util"
	db "better-uptime/internal/db/sqlc"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type SeatRequest struct {
	CoachId   int     `json:"coach_id" validate:"required,min=1"`
	SeatId    int     `json:"seat_id" validate:"required,min=1"`
	BerthType string  `json:"berth_type" validagte:"required,oneof=UP DOWN MID"`
	HoldToken *string `json:"holdToken,omitempty"`
}

type bookingRequest struct {
	TrainId    int           `json:"train_id" validate:"required"`
	TravelDate string        `json:"travel_date" validate:"required"`
	Seats      []SeatRequest `json:"coach_id" validate:"required,min=1"`
}

type BookingResponse struct {
	BookingID   int                `json:"bookingId"`
	Status      string             `json:"status"`
	TotalAmount float64            `json:"totalAmount"`
	Seats       []BookedSeatDetail `json:"seats"`
	PaymentID   *int               `json:"paymentId,omitempty"`
	HoldToken   *string            `json:"holdToken,omitempty"`
	CreatedAt   time.Time          `json:"createdAt"`
}

// BookedSeatDetail represents details of a booked seat
type BookedSeatDetail struct {
	SeatID      int    `json:"seatId"`
	CoachType   string `json:"coachType"`
	CoachNumber int    `json:"coachNumber"`
	SeatNumber  int    `json:"seatNumber"`
	BerthType   string `json:"berthType"`
}

func (h *Handler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	payload, err := middleware.GetFirebasePayloadFromContext(ctx)
	if err != nil {
		util.ErrorJson(w, util.ErrUnauthorized)
		return
	}

	var data bookingRequest

	err = util.ReadJsonAndValidate(w, r, &data)
	if err != nil {
		util.ErrorJson(w, util.ErrNotValidRequest)
		return
	}
	// expire all the old bookings who are pending or cancelled
	err = h.store.ExpireOldBooking(ctx)
	if err != nil {
		util.ErrorJson(w, util.ErrInternal)
		return
	}
	// check if the user have some active bookings
	booking, err := h.store.GetActiveBookingByUser(ctx, pgtype.UUID{Bytes: payload.UserId, Valid: true})
	if err != nil {
		util.ErrorJson(w, util.ErrInternal)
		return
	}

	if booking.ID != 0 {
		util.ErrorJson(w, fmt.Errorf("you already have an acting booking"))
		return
	}
	// check whether that train exists
	count, err := h.store.ValidateTrain(ctx, int32(data.TrainId))
	if err != nil {
		util.ErrorJson(w, util.ErrInternal)
		return
	}
	if count == 0 {
		util.ErrorJson(w, fmt.Errorf("no such train exists"))
		return
	}
	// parse the travelDate from request into time.Time
	travelDate, err := time.Parse("2006-01-02", data.TravelDate)
	if err != nil {
		util.ErrorJson(w, util.ErrNotValidRequest)
		return
	}
	// check whether that train on that scheedule exists
	trainScheduleCount, err := h.store.ValidateSchedule(ctx, db.ValidateScheduleParams{
		Trainid: util.ToPgInt4(int32(data.TrainId)),
		Column2: pgtype.Date{Time: travelDate},
	})
	if err != nil {
		util.ErrorJson(w, util.ErrInternal)
		return
	}
	if trainScheduleCount == 0 {
		util.ErrorJson(w, fmt.Errorf("no such train  schedule exists"))
		return
	}

	seatLength := len(data.Seats)

	var seatIDs []int32

	seatIDs, err = ValidateSeatId(data.Seats)
	if err != nil {
		util.ErrorJson(w, fmt.Errorf("not able to validate seats"))
	}

	_, err = h.store.ValidateSeatsBelongToTrain(ctx, db.ValidateSeatsBelongToTrainParams{
		Column1: int32(seatLength),
		Column2: seatIDs,
	})
	if err != nil {
		util.ErrorJson(w, fmt.Errorf("seat does not belong to that train"))
	}

	// get all the available seats from that train
	availableSeats, err := h.store.CurrentAvailabeSeats(ctx, db.CurrentAvailabeSeatsParams{
		Column1:    seatIDs,
		Trainid:    pgtype.Int4{Int32: int32(data.TrainId)},
		Traveldate: pgtype.Date{Time: travelDate, Valid: true},
	})
	if err != nil {
		util.ErrorJson(w, util.ErrInternal)
		return
	}
	// compare all the available setas with the users'requested seats
	if slices.Equal(availableSeats, seatIDs) {
		util.ErrorJson(w, fmt.Errorf("not all requested seats are available"))
		return
	}

	holdToken := uuid.New().String()
	

}

func ValidateSeatId(Seats []SeatRequest) ([]int32, error) {

	if len(Seats) == 0 {
		return nil, errors.New("no seats selected")
	}

	var Seatids []int32

	seen := make(map[int32]bool)

	for _, seat := range Seats {
		if seat.SeatId <= 0 {
			return nil, errors.New("invalid seat")
		}

		if seen[int32(seat.SeatId)] {
			return nil, errors.New("duplicate seat")
		}

		seen[int32(seat.SeatId)] = true

		Seatids = append(Seatids, int32(seat.SeatId))
	}

	return Seatids, nil

}
