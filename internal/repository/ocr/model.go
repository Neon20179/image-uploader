package ocr

import (
	"time"

	"github.com/google/uuid"

	"image-uploader/internal/domain"
)

type OCRResultModel struct {
	UploadID      uuid.UUID
	Text          *string
	BlockDetected bool
	CreatedAt     time.Time
}

func ocrResultDomainFromModel(model OCRResultModel) domain.OCRResult {
	var text string
	if model.Text != nil {
		text = *model.Text
	}
	return domain.NewOCRResult(
		model.UploadID,
		text,
		model.BlockDetected,
		model.CreatedAt,
	)
}
