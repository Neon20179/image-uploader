package ocr

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"image-uploader/internal/apperr"
	"image-uploader/internal/domain"
	"image-uploader/internal/postgres"
)

func (r *OCRRepository) GetOCRResult(ctx context.Context, uploadID uuid.UUID) (domain.OCRResult, error) {
	querier := r.querierForContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, querier.OpTimeout())
	defer cancel()

	query := `
		SELECT upload_id, text, block_detected, created_at
		FROM ocr_results
		WHERE upload_id=$1;
	`

	row := querier.QueryRow(ctx, query, uploadID)

	var model OCRResultModel

	err := row.Scan(
		&model.UploadID,
		&model.Text,
		&model.BlockDetected,
		&model.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, postgres.ErrNoRows) {
			return domain.OCRResult{}, apperr.ErrNotFound
		}
		return domain.OCRResult{}, fmt.Errorf("scan OCR result: %w", err)
	}

	return ocrResultDomainFromModel(model), nil
}
