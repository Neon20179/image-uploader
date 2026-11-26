package image

import (
	"context"
	"fmt"

	"image-uploader/internal/domain"
)

func (r *ImageRepository) GetImages(ctx context.Context, limit *int, offset *int) ([]domain.Image, error) {
	querier := r.querierForContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, querier.OpTimeout())
	defer cancel()

	query := `
		SELECT id, bucket, object_key, width, height, size, mime_type, created_at
		FROM images
		ORDER BY created_at DESC
		LIMIT $1
		OFFSET $2;
	`

	rows, err := querier.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("select images: %w", err)
	}
	defer rows.Close()

	var images []domain.Image

	for rows.Next() {
		var imageModel ImageModel

		err := rows.Scan(
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
			return nil, fmt.Errorf("scan images: %w", err)
		}

		image := imageDomainFromModel(imageModel)
		images = append(images, image)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("select images rows: %w", err)
	}

	return images, nil
}
