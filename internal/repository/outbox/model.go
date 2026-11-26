package outbox

import "image-uploader/internal/domain"

type OutboxEventModel struct {
	ID      int64
	Topic   string
	Key     []byte
	Payload []byte
	Headers map[string]string
}

func outboxEventDomainFromModel(model OutboxEventModel) domain.OutboxEvent {
	return domain.NewOutboxEvent(model.ID, model.Topic, model.Key, model.Payload, model.Headers)
}
