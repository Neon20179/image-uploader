package embedding

import (
	"context"

	"github.com/google/uuid"

	"image-uploader/internal/domain"
)

func (s *EmbeddingService) FailEmbedding(ctx context.Context, id uuid.UUID) error {
	return s.txManager.WithTx(ctx, func(ctx context.Context) error {
		if err := s.uploadRepository.MarkProcessFailed(ctx, id, domain.UploadProcessEmbedding); err != nil {
			return err
		}
		return s.uploadRepository.MarkFailed(ctx, id)
	})
}
