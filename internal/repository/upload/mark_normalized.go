package upload

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"image-uploader/internal/apperr"
)

func (r *UploadRepository) MarkNormalized(ctx context.Context, id uuid.UUID, objectKey string) error {
	querier := r.querierForContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, querier.OpTimeout())
	defer cancel()

	query := `
		UPDATE uploads
		SET status = 'processing', object_key = $2, normalization_status = 'completed', updated_at = now()
		WHERE id = $1 AND normalization_status = 'pending';
	`
	tag, err := querier.Exec(ctx, query, id, objectKey)
	if err != nil {
		return fmt.Errorf("complete normalization: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("upload with ID=%s: %w", id, apperr.ErrConflict)
	}
	return nil
}
