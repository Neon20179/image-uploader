package image

import (
	"context"

	"github.com/google/uuid"
)

func (s *ImageService) DeleteImage(ctx context.Context, id uuid.UUID) error {
	return s.imageRepository.DeleteImage(ctx, id)
}
