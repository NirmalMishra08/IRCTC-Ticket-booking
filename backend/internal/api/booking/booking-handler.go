package booking

import (
	db "better-uptime/internal/db/sqlc"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
)

type ConsumerHandler struct {
	store db.Store
}

func (h *ConsumerHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerHandler) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {

	for msg := range claim.Messages() {

		switch msg.Topic {

		case "tatkal_booking":
			h.handleTatkal(session, msg)

		case "seat_upgradation":
			h.handleSeat(session, msg)
		}
	}

	return nil
}

type TatkalJob struct {
	BookingID string         `json:"booking_id"`
	UserID    string         `json:"user_id"`
	Data      BookingRequest `json:"data"`
}

func (h *ConsumerHandler) handleTatkal(session sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage) {
	var job TatkalJob

	if err := json.Unmarshal(msg.Value, &job); err != nil {
		log.Println("invalid tatkal message:", err)
		session.MarkMessage(msg, "")
		return
	}

	// call DB or service
	session.MarkMessage(msg, "")
}

func (h *ConsumerHandler) handleSeat(session sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage) {
	var data SeatReleasedEvent

	if err := json.Unmarshal(msg.Value, &data); err != nil {
		log.Println("invalid seat message:", err)
		session.MarkMessage(msg, "")
		return
	}

	// call DB logic
	session.MarkMessage(msg, "")
}
