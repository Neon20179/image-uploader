package kafka

import "context"

type Producer interface {
	Send(ctx context.Context, msg Message) error
	Close(ctx context.Context) error
}
