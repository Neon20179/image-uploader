package upload

import (
	"time"

	"github.com/google/uuid"

	"image-uploader/internal/domain"
)

type UploadDTO struct {
	ID uuid.UUID `json:"id"`

	Status domain.UploadStatus `json:"status"`

	URL string `json:"url"`

	NormalizationStatus domain.ProcessStatus `json:"normalizationStatus"`
	EmbeddingStatus     domain.ProcessStatus `json:"embeddingStatus"`
	OCRStatus           domain.ProcessStatus `json:"ocrStatus"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func uploadDTOFromDomain(upload domain.Upload, uploadURL string) UploadDTO {
	return UploadDTO{
		ID:                  upload.ID,
		Status:              upload.Status,
		URL:                 uploadURL,
		NormalizationStatus: upload.NormalizationStatus,
		EmbeddingStatus:     upload.EmbeddingStatus,
		OCRStatus:           upload.OCRStatus,
		CreatedAt:           upload.CreatedAt,
		UpdatedAt:           upload.UpdatedAt,
	}
}
