package ocr

import (
	"context"
	"fmt"

	"image-uploader/internal/domain"
)

func (r *OCRRepository) CreateOCRResult(ctx context.Context, result domain.OCRResult) (domain.OCRResult, error) {
	querier := r.querierForContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, querier.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO ocr_results (upload_id, text, block_detected, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING upload_id, text, block_detected, created_at;
	`

	row := querier.QueryRow(
		ctx,
		query,
		result.UploadID,
		result.Text,
		result.BlockDetected,
		result.CreatedAt,
	)

	var model OCRResultModel

	err := row.Scan(
		&model.UploadID,
		&model.Text,
		&model.BlockDetected,
		&model.CreatedAt,
	)
	if err != nil {
		return domain.OCRResult{}, fmt.Errorf("create OCR result: %w", err)
	}

	return ocrResultDomainFromModel(model), nil
}
