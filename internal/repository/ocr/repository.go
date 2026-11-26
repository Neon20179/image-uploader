package ocr

import (
	"context"

	"image-uploader/internal/postgres"
)

type OCRRepository struct {
	querier postgres.Querier
}

func NewOCRRepository(querier postgres.Querier) *OCRRepository {
	return &OCRRepository{querier: querier}
}

func (r *OCRRepository) querierForContext(ctx context.Context) postgres.Querier {
	return postgres.QuerierFromContext(ctx, r.querier)
}
