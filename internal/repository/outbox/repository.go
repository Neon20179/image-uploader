package outbox

import (
	"context"

	"image-uploader/internal/postgres"
)

type OutboxRepository struct {
	querier postgres.Querier
}

func NewOutboxRepository(querier postgres.Querier) *OutboxRepository {
	return &OutboxRepository{
		querier: querier,
	}
}

func (r *OutboxRepository) querierForContext(ctx context.Context) postgres.Querier {
	return postgres.QuerierFromContext(ctx, r.querier)
}
