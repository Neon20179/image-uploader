package image

import (
	"context"

	"image-uploader/internal/postgres"
)

type ImageRepository struct {
	querier postgres.Querier
}

func NewImageRepository(querier postgres.Querier) *ImageRepository {
	return &ImageRepository{
		querier: querier,
	}
}

func (r *ImageRepository) querierForContext(ctx context.Context) postgres.Querier {
	return postgres.QuerierFromContext(ctx, r.querier)
}
