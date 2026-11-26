package upload

import (
	"context"
	"fmt"

	"image-uploader/internal/domain"
)

func (r *UploadRepository) GetUploads(ctx context.Context, limit *int, offset *int) ([]domain.Upload, error) {
	querier := r.querierForContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, querier.OpTimeout())
	defer cancel()

	query := `
		SELECT id, image_id, status, bucket, raw_object_key, object_key, normalization_status, embedding_status, ocr_status, created_at, updated_at
		FROM uploads
		ORDER BY created_at DESC
		LIMIT $1
		OFFSET $2;
	`

	rows, err := querier.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("select uploads: %w", err)
	}
	defer rows.Close()

	var uploads []domain.Upload

	for rows.Next() {
		var model UploadModel

		err := rows.Scan(
			&model.ID,
			&model.ImageID,
			&model.Status,
			&model.Bucket,
			&model.RawObjectKey,
			&model.ObjectKey,
			&model.NormalizationStatus,
			&model.EmbeddingStatus,
			&model.OCRStatus,
			&model.CreatedAt,
			&model.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan uploads: %w", err)
		}

		savedUpload := uploadDomainFromModel(model)
		uploads = append(uploads, savedUpload)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("select uploads rows: %w", err)
	}

	return uploads, nil
}
