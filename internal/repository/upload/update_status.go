package upload

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"image-uploader/internal/apperr"
	"image-uploader/internal/domain"
)

func (r *UploadRepository) UpdateStatus(
	ctx context.Context,
	id uuid.UUID,
	from domain.UploadStatus,
	to domain.UploadStatus,
) error {
	querier := r.querierForContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, querier.OpTimeout())
	defer cancel()

	query := `
		UPDATE uploads SET
		status = $3, updated_at = now()
		WHERE id = $1 AND status = $2 
	`
	tag, err := querier.Exec(ctx, query, id, from, to)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("upload with ID=%s: %w", id, apperr.ErrConflict)
	}

	return nil
}
