package upload

import (
	"time"

	"github.com/google/uuid"

	"image-uploader/internal/domain"
)

type UploadModel struct {
	ID uuid.UUID

	ImageID *uuid.UUID

	Status domain.UploadStatus

	Bucket       string
	ObjectKey    *string
	RawObjectKey string

	NormalizationStatus domain.ProcessStatus
	EmbeddingStatus     domain.ProcessStatus
	OCRStatus           domain.ProcessStatus

	CreatedAt time.Time
	UpdatedAt time.Time
}

func uploadDomainFromModel(model UploadModel) domain.Upload {
	return domain.NewUpload(
		model.ID,
		model.ImageID,
		model.Status,
		model.Bucket,
		model.ObjectKey,
		model.RawObjectKey,
		model.NormalizationStatus,
		model.EmbeddingStatus,
		model.OCRStatus,
		model.CreatedAt,
		model.UpdatedAt,
	)
}
