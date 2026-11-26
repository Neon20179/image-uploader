package embedding

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"image-uploader/internal/apperr"
	"image-uploader/internal/domain"
	"image-uploader/internal/postgres"
)

func (r *EmbeddingRepository) GetEmbedding(ctx context.Context, uploadID uuid.UUID) (domain.Embedding, error) {
	querier := r.querierForContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, querier.OpTimeout())
	defer cancel()

	query := `
		SELECT upload_id, vector::text, block_detected, created_at
		FROM embeddings
		WHERE upload_id=$1;
	`

	row := querier.QueryRow(ctx, query, uploadID)

	var model EmbeddingModel

	err := row.Scan(
		&model.UploadID,
		&model.Vector,
		&model.BlockDetected,
		&model.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, postgres.ErrNoRows) {
			return domain.Embedding{}, apperr.ErrNotFound
		}
		return domain.Embedding{}, fmt.Errorf("scan embedding: %w", err)
	}

	return embeddingDomainFromModel(model)
}
