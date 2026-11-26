package upload

import (
	"context"

	"github.com/google/uuid"

	"image-uploader/internal/domain"
)

func (s *UploadService) GetUpload(ctx context.Context, id uuid.UUID) (domain.Upload, error) {
	return s.uploadRepository.GetUpload(ctx, id)
}
