package queue

import (
	"context"
)

type KafkaWriter struct {
}

func (k *KafkaWriter) Publish(ctx context.Context, msg *Message) (err error) {
	return nil
}

func (k *KafkaWriter) Close() error {
	return nil
}
