package upload

import (
	"context"
	"fmt"

	"image-uploader/internal/domain"
)

func (r *UploadRepository) CreateUpload(ctx context.Context, upload domain.Upload) (domain.Upload, error) {
	querier := r.querierForContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, querier.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO uploads (id, image_id, status, bucket, object_key, raw_object_key, normalization_status, embedding_status, ocr_status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, image_id, status, bucket, object_key, raw_object_key, normalization_status, embedding_status, ocr_status, created_at, updated_at;
	`

	row := querier.QueryRow(
		ctx,
		query,
		upload.ID,
		upload.ImageID,
		upload.Status,
		upload.Bucket,
		upload.ObjectKey,
		upload.RawObjectKey,
		upload.NormalizationStatus,
		upload.EmbeddingStatus,
		upload.OCRStatus,
	)

	var model UploadModel
	err := row.Scan(
		&model.ID,
		&model.ImageID,
		&model.Status,
		&model.Bucket,
		&model.ObjectKey,
		&model.RawObjectKey,
		&model.NormalizationStatus,
		&model.EmbeddingStatus,
		&model.OCRStatus,
		&model.CreatedAt,
		&model.UpdatedAt,
	)
	if err != nil {
		return domain.Upload{}, fmt.Errorf("exec query: %w", err)
	}

	return uploadDomainFromModel(model), nil
}
