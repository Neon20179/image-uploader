package normalization

import (
	"context"
	"fmt"

	"image-uploader/internal/domain"
)

func (r *NormalizationRepository) CreateNormalizationResult(ctx context.Context, result domain.NormalizationResult) (domain.NormalizationResult, error) {
	querier := r.querierForContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, querier.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO normalization_results (upload_id, width, height, size, mime_type, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING upload_id, width, height, size, mime_type, created_at;
	`

	row := querier.QueryRow(
		ctx,
		query,
		result.UploadID,
		result.Width,
		result.Height,
		result.Size,
		result.MIMEType,
		result.CreatedAt,
	)

	var model NormalizationResultModel

	if err := row.Scan(
		&model.UploadID,
		&model.Width,
		&model.Height,
		&model.Size,
		&model.MIMEType,
		&model.CreatedAt,
	); err != nil {
		return domain.NormalizationResult{}, fmt.Errorf("create normalization result: %w", err)
	}

	return normalizationResultDomainFromModel(model), nil
}
