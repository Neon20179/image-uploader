package domain

import (
	"time"

	"github.com/google/uuid"
)

type BlacklistedEmbedding struct {
	ID        uuid.UUID
	Vector    []float32
	CreatedAt time.Time
}

func NewBlacklistedEmbedding(id uuid.UUID, vector []float32, createdAt time.Time) BlacklistedEmbedding {
	return BlacklistedEmbedding{
		ID:        id,
		Vector:    vector,
		CreatedAt: createdAt,
	}
}

func NewBlacklistedEmbeddingUninitialized(id uuid.UUID, vector []float32) BlacklistedEmbedding {
	return NewBlacklistedEmbedding(id, vector, time.Now())
}
