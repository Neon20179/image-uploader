package normalization

import (
	"time"

	"github.com/google/uuid"

	"image-uploader/internal/domain"
)

type NormalizationResultModel struct {
	UploadID  uuid.UUID
	Width     int
	Height    int
	Size      int
	MIMEType  string
	CreatedAt time.Time
}

func normalizationResultDomainFromModel(model NormalizationResultModel) domain.NormalizationResult {
	return domain.NewNormalizationResult(
		model.UploadID,
		model.Width,
		model.Height,
		model.Size,
		model.MIMEType,
		model.CreatedAt,
	)
}
