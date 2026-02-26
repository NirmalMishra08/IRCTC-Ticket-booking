package booking

import (
	"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
)

func (t *Handler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (t *Handler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

type TatkalJob struct {
	BookingID string         `json:"booking_id"`
	UserID    string         `json:"user_id"`
	Data      BookingRequest `json:"data"`
}

func (t *Handler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var job TatkalJob

		if err := json.Unmarshal(msg.Value, &job); err != nil {
			log.Println("invalid tatkal message:", err)
			session.MarkMessage(msg, "")
			continue
		}

		

		err := t.ProcessTatkalBooking(
			context.Background(),
			job.Data,
			job.UserID,
			job.BookingID,
		)
		if err != nil {
			log.Println("booking failed:", err)
			continue // kafka will retry
		}
		session.MarkMessage(msg, "")

	}

	return nil
}
