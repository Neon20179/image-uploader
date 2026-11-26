package embedding

import (
	"context"
	"encoding/json"
	"fmt"

	"image-uploader/internal/domain"
)

func (r *EmbeddingRepository) CreateEmbedding(ctx context.Context, embedding domain.Embedding) (domain.Embedding, error) {
	vector, err := json.Marshal(embedding.Vector)
	if err != nil {
		return domain.Embedding{}, fmt.Errorf("encode embedding vector: %w", err)
	}

	querier := r.querierForContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, querier.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO embeddings (upload_id, vector, block_detected, created_at)
		VALUES ($1, $2::text::vector, $3, $4)
		RETURNING upload_id, vector::text, block_detected, created_at;
	`

	row := querier.QueryRow(ctx,
		query,
		embedding.UploadID,
		string(vector),
		embedding.BlockDetected,
		embedding.CreatedAt,
	)

	var model EmbeddingModel

	err = row.Scan(
		&model.UploadID,
		&model.Vector,
		&model.BlockDetected,
		&model.CreatedAt,
	)
	if err != nil {
		return domain.Embedding{}, fmt.Errorf("create embedding: %w", err)
	}

	return embeddingDomainFromModel(model)
}
