package image

import (
	"context"

	"github.com/google/uuid"

	"image-uploader/internal/domain"
)

func (s *ImageService) GetImage(ctx context.Context, id uuid.UUID) (domain.Image, error) {
	return s.imageRepository.GetImage(ctx, id)
}
