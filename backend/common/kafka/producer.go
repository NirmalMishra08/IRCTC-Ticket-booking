package kafka


import "context"

type Producer interface {
	Publish(ctx context.Context, topic string, key string, value []byte) error
	Close() error
}