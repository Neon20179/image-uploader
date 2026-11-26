package upload

import (
	"context"

	"github.com/google/uuid"

	"image-uploader/internal/domain"
)

func (s *UploadService) FailNormalization(ctx context.Context, id uuid.UUID) error {
	return s.txManager.WithTx(ctx, func(ctx context.Context) error {
		if err := s.uploadRepository.MarkProcessFailed(ctx, id, domain.UploadProcessNormalization); err != nil {
			return err
		}
		return s.uploadRepository.MarkFailed(ctx, id)
	})
}
