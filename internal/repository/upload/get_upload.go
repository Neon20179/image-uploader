package upload

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"image-uploader/internal/apperr"
	"image-uploader/internal/domain"
	"image-uploader/internal/postgres"
)

func (r *UploadRepository) GetUpload(
	ctx context.Context,
	id uuid.UUID,
) (domain.Upload, error) {
	return r.getUpload(ctx, id, false)
}

func (r *UploadRepository) GetUploadForUpdate(
	ctx context.Context,
	id uuid.UUID,
) (domain.Upload, error) {
	return r.getUpload(ctx, id, true)
}

func (r *UploadRepository) getUpload(
	ctx context.Context,
	id uuid.UUID,
	forUpdate bool,
) (domain.Upload, error) {
	querier := r.querierForContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, querier.OpTimeout())
	defer cancel()

	query := `
		SELECT id, image_id, status, bucket, object_key, raw_object_key, normalization_status, embedding_status, ocr_status, created_at, updated_at
		FROM uploads
		WHERE id=$1
		%s;
	`

	if forUpdate {
		query = fmt.Sprintf(query, "FOR UPDATE")
	} else {
		query = fmt.Sprintf(query, "")
	}

	row := querier.QueryRow(ctx, query, id)

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
		if errors.Is(err, postgres.ErrNoRows) {
			return domain.Upload{}, apperr.ErrNotFound
		}
		return domain.Upload{}, fmt.Errorf("scan upload: %w", err)
	}

	return uploadDomainFromModel(model), nil
}
