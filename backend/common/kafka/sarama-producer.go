package kafka

import (
	"context"

	"github.com/IBM/sarama"
)

type SaramaProducer struct {
	producer sarama.SyncProducer
}

func NewSaramaProducer(brokerURL []string) (*SaramaProducer, error) {

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	p, err := sarama.NewSyncProducer(brokerURL, config)
	if err != nil {
		return nil, err
	}

	return &SaramaProducer{
		producer: p,
	}, nil
}

func (h *SaramaProducer) Publish(ctx context.Context, topic string, key string, value []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(value),
	}

	if key != "" {
		msg.Key = sarama.StringEncoder(key)
	}

	_, _, err := h.producer.SendMessage(msg)

	return err

}
func (h *SaramaProducer) Close() error {
	err := h.producer.Close()
	return err
}


// do it in main.go file
// kafkaProducer, err := kafka.NewSaramaProducer(
// 	[]string{"localhost:9092"},
// )
// if err != nil {
// 	log.Fatal(err)
// }
// defer kafkaProducer.Close()

// handler := &Handler{
// 	Kafka: kafkaProducer,
// }