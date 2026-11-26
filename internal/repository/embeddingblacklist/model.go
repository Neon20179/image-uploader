package embeddingblacklist

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"image-uploader/internal/domain"
)

type EmbeddingModel struct {
	ID        uuid.UUID
	Vector    string
	CreatedAt time.Time
}

func embeddingDomainFromModel(model EmbeddingModel) (domain.BlacklistedEmbedding, error) {
	var vector []float32
	if err := json.Unmarshal([]byte(model.Vector), &vector); err != nil {
		return domain.BlacklistedEmbedding{}, fmt.Errorf("decode embedding vector: %w", err)
	}
	return domain.NewBlacklistedEmbedding(
		model.ID,
		vector,
		model.CreatedAt,
	), nil
}
