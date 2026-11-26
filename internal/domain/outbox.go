package domain

type OutboxEvent struct {
	ID      int64
	Topic   string
	Key     []byte
	Payload []byte
	Headers map[string]string
}

func NewOutboxEvent(
	id int64,
	topic string,
	key []byte,
	payload []byte,
	headers map[string]string,
) OutboxEvent {
	return OutboxEvent{
		ID:      id,
		Topic:   topic,
		Key:     key,
		Payload: payload,
		Headers: headers,
	}
}

func NewOutboxEventUninitialized(
	topic string,
	key []byte,
	payload []byte,
	headers map[string]string,
) OutboxEvent {
	return NewOutboxEvent(-1, topic, key, payload, headers)
}
