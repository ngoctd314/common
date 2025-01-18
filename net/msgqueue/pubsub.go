package queue

import "context"

type Publisher interface {
	Publish(ctx context.Context, msg *Message) error
}

type Subscriber interface {
	Subscribe(ctx context.Context) (*Message, error)
}
