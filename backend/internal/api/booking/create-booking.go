package booking

import (
	"better-uptime/common/logger"
	"better-uptime/common/middleware"
	"better-uptime/common/stripe"
	"better-uptime/common/util"
	db "better-uptime/internal/db/sqlc"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type SeatRequest struct {
	CoachId   int     `json:"coach_id" validate:"required,min=1"`
	SeatId    int     `json:"seat_id" validate:"required,min=1"`
	BerthType string  `json:"berth_type" validate:"required,oneof=UP DOWN MID"`
	HoldToken *string `json:"holdToken,omitempty"`
}

type bookingRequest struct {
	TrainId    int           `json:"train_id" validate:"required"`
	TravelDate string        `json:"travel_date" validate:"required"`
	Seats      []SeatRequest `json:"seats" validate:"required,min=1"`
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
	if err != nil && err.Error() != "no rows in result set" {
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
	logger.Debug(data.TravelDate, travelDate)
	// check whether that train on that scheedule exists
	trainScheduleCount, err := h.store.ValidateSchedule(ctx, db.ValidateScheduleParams{
		Trainid: util.ToPgInt4(int32(data.TrainId)),
		Column2: pgtype.Date{Time: travelDate, Valid: true},
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
		return
	}

	_, err = h.store.ValidateSeatsBelongToTrain(ctx, db.ValidateSeatsBelongToTrainParams{
		Column1: int32(seatLength),
		Column2: seatIDs,
	})
	if err != nil {
		util.ErrorJson(w, fmt.Errorf("seat does not belong to that train"))
		return
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
	if !slices.Equal(availableSeats, seatIDs) {
		util.ErrorJson(w, fmt.Errorf("not all requested seats are available"))
		return
	}

	holdToken := uuid.New().String()
	TrainId := fmt.Sprintf("%d", data.TrainId)

	seatIDStrs := make([]string, len(seatIDs))
	for i, id := range seatIDs {
		seatIDStrs[i] = fmt.Sprintf("%d", id)
	}

	Lock, err := h.TrySeatLock(ctx, TrainId, data.TravelDate, seatIDStrs, holdToken, 5*time.Minute)
	if err != nil {
		util.ErrorJson(w, util.ErrInternal)
		logger.Debug("here--6>")
		return
	}

	if Lock == false {
		util.ErrorJson(w, fmt.Errorf("not able to lock all the tickets"))
		return
	}

	defer func() {
		if err != nil {
			h.ReleaseLocks(ctx, TrainId, data.TravelDate, seatIDStrs, holdToken)
		}
	}()

	var bookingID int32

	err = h.store.ExecTx(ctx, func(q *db.Queries) error {
		bookingInProcess, err := q.CreateBooking(ctx, db.CreateBookingParams{
			Userid:     pgtype.UUID{Bytes: payload.UserId, Valid: true},
			Trainid:    util.ToPgInt4(int32(data.TrainId)),
			Traveldate: pgtype.Date{Time: travelDate, Valid: true},
			Holdtoken:  pgtype.Text{String: holdToken, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("create booking: %w", err)
		}

		bookingID = bookingInProcess.ID

		
		trainScheduled, err := q.GetTrainScheduleByDay(ctx, db.GetTrainScheduleByDayParams{
			Trainid: util.ToPgInt4(int32(data.TrainId)),
			Column2: pgtype.Date{Time: travelDate,Valid: true},
		})
		if err != nil {
			return fmt.Errorf("get train schedule: %w", err)

		}
		var bookingItems []db.Bookingitem

		for _, seatId := range seatIDs {
			bookingItem, err := q.CreateBookingItem(ctx, db.CreateBookingItemParams{
				Bookingid:       util.ToPgInt4(bookingID),
				Seatid:          util.ToPgInt4(seatId),
				Trainscheduleid: util.ToPgInt4(trainScheduled.ID),
			})
			if err != nil {
				return fmt.Errorf("create booking item for seat %d: %w", seatId, err)
			}
			bookingItems = append(bookingItems, bookingItem)
		}

		return nil

	})

	if err != nil {
		// Release locks on transaction failure
		h.ReleaseLocks(ctx, TrainId, data.TravelDate, seatIDStrs, holdToken)
		util.ErrorJson(w, fmt.Errorf("booking creation failed: %v", err))
		return
	}

	var amount int32
	amount = CalculateFare(int32(len(seatIDs)))

	paymentIntent, err := stripe.StripeSession(ctx, payload.UserId.String(), strconv.Itoa(int(amount)), "seatIds:", h.config.STRIPE_SECRET_KEY, int(bookingID), holdToken)
	if err != nil {

		updateErr := h.store.UpdateBookingStatus(ctx, db.UpdateBookingStatusParams{
			ID:     bookingID,
			Status: db.BookingStatusEXPIRED,
		})
		if updateErr != nil {
			fmt.Printf("Critical: Failed to update booking status: %v\n", updateErr)
		}
		lockErr := h.ReleaseLocks(ctx, TrainId, data.TravelDate, seatIDStrs, holdToken)
		if lockErr != nil {
			fmt.Printf("Critical: Failed to unlock seats: %v\n", lockErr)
		}
		util.ErrorJson(w, fmt.Errorf("not able to create session for payment"))
		return
	}

	transactionId := paymentIntent.SessionURL.SessionID

	_, err = h.store.CreatePayment(ctx, db.CreatePaymentParams{
		Bookingid:     util.ToPgInt4(bookingID),
		Amount:        float64(amount),
		Transactionid: transactionId,
	})
	if err != nil {
		util.ErrorJson(w, fmt.Errorf("not able to create payment for this transaction"))
		return
	}

	response := map[string]interface{}{
		"bookingId":      bookingID,
		"hold_token":     holdToken,
		"redirected_url": paymentIntent.SessionURL.SessionURL,
	}

	util.WriteJson(w, http.StatusAccepted, response)

}

func CalculateFare(seatIds int32) int32 {
	return 100 * seatIds
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
