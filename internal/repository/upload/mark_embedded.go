package upload

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"image-uploader/internal/apperr"
)

func (r *UploadRepository) MarkEmbedded(ctx context.Context, id uuid.UUID) error {
	querier := r.querierForContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, querier.OpTimeout())
	defer cancel()

	query := `
		UPDATE uploads
		SET embedding_status = 'completed', updated_at = now()
		WHERE id = $1 AND embedding_status = 'pending';
	`
	tag, err := querier.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("complete embedding: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("upload with ID=%s: %w", id, apperr.ErrConflict)
	}
	return nil
}
