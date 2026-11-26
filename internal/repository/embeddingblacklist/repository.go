package embeddingblacklist

import (
	"context"

	"image-uploader/internal/postgres"
)

type EmbeddingBlacklistRepository struct {
	querier postgres.Querier
}

func NewEmbeddingBlacklistRepository(querier postgres.Querier) *EmbeddingBlacklistRepository {
	return &EmbeddingBlacklistRepository{querier: querier}
}

func (r *EmbeddingBlacklistRepository) querierForContext(ctx context.Context) postgres.Querier {
	return postgres.QuerierFromContext(ctx, r.querier)
}
