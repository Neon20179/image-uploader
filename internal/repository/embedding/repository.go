package embedding

import (
	"context"

	"image-uploader/internal/postgres"
)

type EmbeddingRepository struct {
	querier postgres.Querier
}

func NewEmbeddingRepository(querier postgres.Querier) *EmbeddingRepository {
	return &EmbeddingRepository{querier: querier}
}

func (r *EmbeddingRepository) querierForContext(ctx context.Context) postgres.Querier {
	return postgres.QuerierFromContext(ctx, r.querier)
}
