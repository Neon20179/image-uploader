package domain

import (
	"time"

	"github.com/google/uuid"
)

type UploadStatus string

const (
	UploadStatusPending    UploadStatus = "pending"
	UploadStatusUploaded   UploadStatus = "uploaded"
	UploadStatusProcessing UploadStatus = "processing"
	UploadStatusReady      UploadStatus = "ready"
	UploadStatusBlocked    UploadStatus = "blocked"
	UploadStatusFailed     UploadStatus = "failed"
)

type ProcessStatus string

const (
	ProcessStatusPending   ProcessStatus = "pending"
	ProcessStatusCompleted ProcessStatus = "completed"
	ProcessStatusFailed    ProcessStatus = "failed"
)

type Upload struct {
	ID uuid.UUID

	ImageID *uuid.UUID

	Status UploadStatus

	Bucket       string
	ObjectKey    *string
	RawObjectKey string

	NormalizationStatus ProcessStatus
	EmbeddingStatus     ProcessStatus
	OCRStatus           ProcessStatus

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUpload(
	id uuid.UUID,
	imageID *uuid.UUID,
	status UploadStatus,
	bucket string,
	objectKey *string,
	rawObjectKey string,
	normalizationStatus ProcessStatus,
	embeddingStatus ProcessStatus,
	ocrStatus ProcessStatus,
	createdAt time.Time,
	updatedAt time.Time,
) Upload {
	return Upload{
		ID:                  id,
		ImageID:             imageID,
		Status:              status,
		Bucket:              bucket,
		ObjectKey:           objectKey,
		RawObjectKey:        rawObjectKey,
		NormalizationStatus: normalizationStatus,
		EmbeddingStatus:     embeddingStatus,
		OCRStatus:           ocrStatus,
		CreatedAt:           createdAt,
		UpdatedAt:           updatedAt,
	}
}

func NewUploadUninitialized(
	id uuid.UUID,
	bucket string,
	rawObjectKey string,
) Upload {
	now := time.Now()

	return NewUpload(
		id,
		nil,
		UploadStatusPending,
		bucket,
		nil,
		rawObjectKey,
		ProcessStatusPending,
		ProcessStatusPending,
		ProcessStatusPending,
		now,
		now,
	)
}

type Embedding struct {
	UploadID uuid.UUID

	Vector []float32

	BlockDetected bool

	CreatedAt time.Time
}

func NewEmbedding(
	uploadID uuid.UUID,
	vector []float32,
	blockDetected bool,
	createdAt time.Time,
) Embedding {
	return Embedding{
		UploadID:      uploadID,
		Vector:        vector,
		BlockDetected: blockDetected,
		CreatedAt:     createdAt,
	}
}

func NewEmbeddingUninitialized(
	uploadID uuid.UUID,
	vector []float32,
	blockDetected bool,
) Embedding {
	now := time.Now()

	return NewEmbedding(uploadID, vector, blockDetected, now)
}

type OCRResult struct {
	UploadID uuid.UUID

	Text string

	BlockDetected bool

	CreatedAt time.Time
}

func NewOCRResult(
	uploadID uuid.UUID,
	text string,
	blockDetected bool,
	createdAt time.Time,
) OCRResult {
	return OCRResult{
		UploadID:      uploadID,
		Text:          text,
		BlockDetected: blockDetected,
		CreatedAt:     createdAt,
	}
}

func NewOCRResultUninitialized(
	uploadID uuid.UUID,
	text string,
	blockDetected bool,
) OCRResult {
	now := time.Now()

	return NewOCRResult(uploadID, text, blockDetected, now)
}
