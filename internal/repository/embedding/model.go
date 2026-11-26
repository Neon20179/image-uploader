package embedding

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"image-uploader/internal/domain"
)

type EmbeddingModel struct {
	UploadID      uuid.UUID
	Vector        string
	BlockDetected bool
	CreatedAt     time.Time
}

func embeddingDomainFromModel(model EmbeddingModel) (domain.Embedding, error) {
	var vector []float32
	if err := json.Unmarshal([]byte(model.Vector), &vector); err != nil {
		return domain.Embedding{}, fmt.Errorf("decode embedding vector: %w", err)
	}
	return domain.NewEmbedding(
		model.UploadID,
		vector,
		model.BlockDetected,
		model.CreatedAt,
	), nil
}
