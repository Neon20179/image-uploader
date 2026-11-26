package outbox

import (
	"context"
	"fmt"

	"image-uploader/internal/domain"
)

func (r *OutboxRepository) Add(ctx context.Context, event domain.OutboxEvent) (domain.OutboxEvent, error) {
	querier := r.querierForContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, querier.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO outbox (topic, key, payload, headers)
		VALUES ($1, $2, $3, $4)
		RETURNING id, topic, key, payload, headers
	`

	row := querier.QueryRow(ctx, query, event.Topic, event.Key, event.Payload, event.Headers)
	var model OutboxEventModel
	err := row.Scan(
		&model.ID,
		&model.Topic,
		&model.Key,
		&model.Payload,
		&model.Headers,
	)
	if err != nil {
		return domain.OutboxEvent{}, fmt.Errorf("scan row: %w", err)
	}

	return outboxEventDomainFromModel(model), nil
}
