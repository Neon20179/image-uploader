package outbox

import (
	"context"
	"fmt"

	"image-uploader/internal/domain"
)

func (r OutboxRepository) GetUnpublished(ctx context.Context, limit int) ([]domain.OutboxEvent, error) {
	querier := r.querierForContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, querier.OpTimeout())
	defer cancel()

	query := `
		SELECT id, topic, key, payload, headers
		FROM outbox
		WHERE published_at IS NULL
		ORDER BY id
		LIMIT $1
		FOR UPDATE SKIP LOCKED
	`

	rows, err := querier.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("select outbox events: %w", err)
	}
	defer rows.Close()

	var outboxEvents []domain.OutboxEvent
	for rows.Next() {
		var model OutboxEventModel
		err := rows.Scan(
			&model.ID,
			&model.Topic,
			&model.Key,
			&model.Payload,
			&model.Headers,
		)
		if err != nil {
			return nil, fmt.Errorf("scan outbox event: %w", err)
		}
		outboxEvent := outboxEventDomainFromModel(model)
		outboxEvents = append(outboxEvents, outboxEvent)
	}

	return outboxEvents, rows.Err()
}
