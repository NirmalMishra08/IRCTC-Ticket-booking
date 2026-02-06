package tatkal

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

func (t *Handler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var job TatkalJob

		if err := json.Unmarshal(msg.Value, &job); err != nil {
			log.Println("invalid tatkal message:", err)
			session.MarkMessage(msg, "")
			continue
		}

		err := t.ProcessTatkalUser(
			context.Background(),
			job.Data,
			job.UserID,
		)
		if err != nil {
			log.Println("booking failed:", err)
			continue // kafka will retry
		}
		session.MarkMessage(msg, "")

	}

	return nil
}
