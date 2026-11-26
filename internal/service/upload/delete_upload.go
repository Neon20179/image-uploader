package upload

import (
	"context"

	"github.com/google/uuid"
)

func (s *UploadService) DeleteUpload(ctx context.Context, id uuid.UUID) error {
	return s.uploadRepository.DeleteUpload(ctx, id)
}
