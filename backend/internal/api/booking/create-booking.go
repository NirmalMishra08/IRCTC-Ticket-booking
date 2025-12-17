package booking

import (
	"better-uptime/common/middleware"
	"better-uptime/common/util"
	"net/http"
	"time"
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
	userId := payload.UserId

}
