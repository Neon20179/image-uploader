package kafka

import "context"

type ConsumerHandler interface {
	Handle(ctx context.Context, msg Message) error
}

type DeadLetterHandler interface {
	OnDeadLetter(ctx context.Context, msg Message) error
}
