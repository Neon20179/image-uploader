package upload

import (
	"github.com/google/uuid"

	"image-uploader/internal/domain"
)

type CreateUploadResult struct {
	UploadID  uuid.UUID
	Status    domain.UploadStatus
	UploadURL string
	Fields    map[string]string
}
