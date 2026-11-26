package upload

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"image-uploader/internal/domain"
)

func (r *UploadRepository) MarkProcessFailed(ctx context.Context, id uuid.UUID, process domain.UploadProcess) error {
	querier := r.querierForContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, querier.OpTimeout())
	defer cancel()

	var column string
	switch process {
	case domain.UploadProcessNormalization:
		column = "normalization_status"
	case domain.UploadProcessOCR:
		column = "ocr_status"
	case domain.UploadProcessEmbedding:
		column = "embedding_status"
	default:
		return fmt.Errorf("unknown upload process: %q", process)
	}

	query := fmt.Sprintf(`
		UPDATE uploads SET %s = 'failed', updated_at = now()
		WHERE id = $1
	`, column)

	if _, err := querier.Exec(ctx, query, id); err != nil {
		return fmt.Errorf("mark %s process failed: %w", process, err)
	}
	return nil
}
