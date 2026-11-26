package image

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"image-uploader/internal/apperr"
)

func (r *ImageRepository) DeleteImage(ctx context.Context, id uuid.UUID) error {
	querier := r.querierForContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, querier.OpTimeout())
	defer cancel()

	query := `
		DELETE FROM images
		WHERE id=$1;
	`

	tag, err := querier.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("image with id=%s: %w", id, apperr.ErrNotFound)
	}

	return nil
}
