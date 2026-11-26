package upload

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"image-uploader/internal/apperr"
)

func (r *UploadRepository) MarkReady(ctx context.Context, id uuid.UUID, imageID uuid.UUID) error {
	querier := r.querierForContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, querier.OpTimeout())
	defer cancel()

	query := `
		UPDATE uploads
		SET image_id = $2, status = 'ready', updated_at = now()
		WHERE id = $1 AND status = 'processing' AND image_id IS NULL;
	`
	tag, err := querier.Exec(ctx, query, id, imageID)
	if err != nil {
		return fmt.Errorf("mark upload ready: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("upload with ID=%s: %w", id, apperr.ErrConflict)
	}
	return nil
}
