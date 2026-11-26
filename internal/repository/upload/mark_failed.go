package upload

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *UploadRepository) MarkFailed(ctx context.Context, id uuid.UUID) error {
	querier := r.querierForContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, querier.OpTimeout())
	defer cancel()

	query := `
		UPDATE uploads SET status = 'failed', updated_at = now()
		WHERE id = $1;
	`

	if _, err := querier.Exec(ctx, query, id); err != nil {
		return fmt.Errorf("mark upload failed: %w", err)
	}

	return nil
}
