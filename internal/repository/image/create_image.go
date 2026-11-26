package image

import (
	"context"
	"fmt"

	"image-uploader/internal/domain"
)

func (r *ImageRepository) CreateImage(ctx context.Context, image domain.Image) (domain.Image, error) {
	querier := r.querierForContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, querier.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO images (id, bucket, object_key, width, height, size, mime_type, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, bucket, object_key, width, height, size, mime_type, created_at;
	`

	row := querier.QueryRow(
		ctx,
		query,
		image.ID,
		image.Bucket,
		image.ObjectKey,
		image.Width,
		image.Height,
		image.Size,
		image.MIMEType,
		image.CreatedAt,
	)

	var model ImageModel

	err := row.Scan(
		&model.ID,
		&model.Bucket,
		&model.ObjectKey,
		&model.Width,
		&model.Height,
		&model.Size,
		&model.MIMEType,
		&model.CreatedAt,
	)
	if err != nil {
		return domain.Image{}, fmt.Errorf("create image: %w", err)
	}

	return imageDomainFromModel(model), nil
}
