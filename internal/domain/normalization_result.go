package domain

import (
	"time"

	"github.com/google/uuid"
)

type NormalizationResult struct {
	UploadID  uuid.UUID
	Width     int
	Height    int
	Size      int
	MIMEType  string
	CreatedAt time.Time
}

func NewNormalizationResult(uploadID uuid.UUID, width, height, size int, mimeType string, createdAt time.Time) NormalizationResult {
	return NormalizationResult{
		UploadID:  uploadID,
		Width:     width,
		Height:    height,
		Size:      size,
		MIMEType:  mimeType,
		CreatedAt: createdAt,
	}
}

func NewNormalizationResultUninitialized(uploadID uuid.UUID, width, height, size int, mimeType string) NormalizationResult {
	return NewNormalizationResult(uploadID, width, height, size, mimeType, time.Now())
}
