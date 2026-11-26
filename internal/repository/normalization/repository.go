package normalization

import (
	"context"

	"image-uploader/internal/postgres"
)

type NormalizationRepository struct {
	querier postgres.Querier
}

func NewNormalizationRepository(querier postgres.Querier) *NormalizationRepository {
	return &NormalizationRepository{querier: querier}
}

func (r *NormalizationRepository) querierForContext(ctx context.Context) postgres.Querier {
	return postgres.QuerierFromContext(ctx, r.querier)
}
