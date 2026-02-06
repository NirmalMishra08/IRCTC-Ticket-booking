package kafka

import (
	"context"
	"log"

	"github.com/IBM/sarama"
)

type SaramaConsumer struct {
	name    string
	group   sarama.ConsumerGroup
	topic   string
	handler sarama.ConsumerGroupHandler
}

func NewSaramaConsumer(name string, brokerUrl []string, groupId string, topic string, handler sarama.ConsumerGroupHandler) (*SaramaConsumer, error) {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V2_1_0_0
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest

	group, err := sarama.NewConsumerGroup(brokerUrl, groupId, cfg)
	if err != nil {
		return nil, err
	}

	return &SaramaConsumer{
		name:    name,
		group:   group,
		topic:   topic,
		handler: handler,
	}, nil
}

func (c *SaramaConsumer) Start(ctx context.Context) error {
	log.Println("starting consumer:", c.name)
	for {
		if err := c.group.Consume(ctx, []string{c.topic}, c.handler); err != nil {
			log.Println("consumer error:", c.name, err)
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

func (s *SaramaConsumer) Close() error {
	return s.group.Close()
}
