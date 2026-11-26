package image

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"image-uploader/internal/apperr"
	"image-uploader/internal/domain"
	"image-uploader/internal/postgres"
)

func (r *ImageRepository) GetImage(ctx context.Context, id uuid.UUID) (domain.Image, error) {
	querier := r.querierForContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, querier.OpTimeout())
	defer cancel()

	query := `
		SELECT id, bucket, object_key, width, height, size, mime_type, created_at
		FROM images
		WHERE id=$1;
	`

	row := querier.QueryRow(ctx, query, id)

	var imageModel ImageModel

	err := row.Scan(
		&imageModel.ID,
		&imageModel.Bucket,
		&imageModel.ObjectKey,
		&imageModel.Width,
		&imageModel.Height,
		&imageModel.Size,
		&imageModel.MIMEType,
		&imageModel.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, postgres.ErrNoRows) {
			return domain.Image{}, apperr.ErrNotFound
		}
		return domain.Image{}, fmt.Errorf("scan image: %w", err)
	}

	image := imageDomainFromModel(imageModel)

	return image, nil
}
