package embeddingblacklist

import (
	"context"
	"encoding/json"
	"fmt"
)

func (r *EmbeddingBlacklistRepository) HasSimilarEmbedding(ctx context.Context, vector []float32, similarityThreshold float64) (bool, error) {
	encoded, err := json.Marshal(vector)
	if err != nil {
		return false, fmt.Errorf("encode embedding vector: %w", err)
	}

	querier := r.querierForContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, querier.OpTimeout())
	defer cancel()

	query := `
		SELECT EXISTS (
			SELECT 1 FROM embedding_blacklist
			WHERE 1 - (vector <=> $1::text::vector) >= $2
		);
	`

	row := querier.QueryRow(
		ctx,
		query,
		string(encoded),
		similarityThreshold,
	)

	var found bool

	err = row.Scan(&found)
	if err != nil {
		return false, fmt.Errorf("search embedding blacklist: %w", err)
	}

	return found, nil
}
