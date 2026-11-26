package upload

import (
	"context"

	"image-uploader/internal/postgres"
)

type UploadRepository struct {
	querier postgres.Querier
}

func NewUploadRepository(querier postgres.Querier) *UploadRepository {
	return &UploadRepository{
		querier: querier,
	}
}

func (r *UploadRepository) querierForContext(ctx context.Context) postgres.Querier {
	return postgres.QuerierFromContext(ctx, r.querier)
}
