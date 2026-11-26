package normalization

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"image-uploader/internal/apperr"
	"image-uploader/internal/domain"
	"image-uploader/internal/postgres"
)

func (r *NormalizationRepository) GetNormalizationResult(ctx context.Context, uploadID uuid.UUID) (domain.NormalizationResult, error) {
	querier := r.querierForContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, querier.OpTimeout())
	defer cancel()

	query := `
		SELECT upload_id, width, height, size, mime_type, created_at
		FROM normalization_results
		WHERE upload_id = $1;
	`
	row := querier.QueryRow(ctx, query, uploadID)

	var model NormalizationResultModel

	if err := row.Scan(
		&model.UploadID,
		&model.Width,
		&model.Height,
		&model.Size,
		&model.MIMEType,
		&model.CreatedAt,
	); err != nil {
		if errors.Is(err, postgres.ErrNoRows) {
			return domain.NormalizationResult{}, apperr.ErrNotFound
		}
		return domain.NormalizationResult{}, fmt.Errorf("scan normalization result: %w", err)
	}
	return normalizationResultDomainFromModel(model), nil
}
