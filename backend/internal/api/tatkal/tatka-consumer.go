package tatkal

import "github.com/IBM/sarama"

type TatkalConsumer struct {
	Handler *Handler
}

func (t *TatkalConsumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (t *TatkalConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}
