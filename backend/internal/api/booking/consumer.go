package booking

import (
	db "better-uptime/internal/db/sqlc"
	"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
)

type SeatReleasedEvent struct {
	JourneyId string `json:"journey_id,omitempty"`
	CoachType string `json:"coach_type,omitempty"`
}

func (h *Handler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	for msg := range claim.Messages() {

		switch msg.Topic {

		case "tatkal_booking":
			h.handleTatkal(session, msg)

		case "seat_upgradation":
			h.handleSeatUpgradation(session, msg)
		}
	}

	return nil
}

func (h *Handler) handleTatkal(session sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage) {
	var job TatkalJob

	if err := json.Unmarshal(msg.Value, &job); err != nil {
		log.Println("invalid tatkal message:", err)
		session.MarkMessage(msg, "")
		return
	}

	err := h.ProcessTatkalBooking(context.Background(), job.Data, job.UserID, job.BookingID)
	if err != nil {
		log.Println("booking failed:", err)
		return
	}

	session.MarkMessage(msg, "")
}

func (h *Handler) handleSeatUpgradation(session sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage) {
	var data SeatReleasedEvent

	if err := json.Unmarshal(msg.Value, &data); err != nil {
		log.Println("invalid seat message:", err)
		session.MarkMessage(msg, "")
		return
	}

	err := h.PromoteWaitlist(context.Background(), data.JourneyId, db.CoachType(data.CoachType))
	if err != nil {
		log.Println("seat upgradation failed:", err)
		return
	}

	session.MarkMessage(msg, "")
}
