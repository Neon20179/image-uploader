package outbox

import (
	"context"
	"fmt"
)

func (r *OutboxRepository) MarkPublished(ctx context.Context, ids []int64) error {
	querier := r.querierForContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, querier.OpTimeout())
	defer cancel()

	query := `
		UPDATE outbox
		SET published_at = now()
		WHERE id = ANY($1)
	`

	_, err := querier.Exec(ctx, query, ids)
	if err != nil {
		return fmt.Errorf("mark published outbox events: %w", err)
	}

	return nil
}
