package image

import (
	"context"
	"fmt"

	"image-uploader/internal/apperr"
	"image-uploader/internal/domain"
)

func (s *ImageService) GetImages(ctx context.Context, limit *int, offset *int) ([]domain.Image, int, error) {
	if limit != nil && *limit < 0 {
		return nil, 0, fmt.Errorf("negative 'limit' argument: %w", apperr.ErrInvalidArgument)
	}

	if offset != nil && *offset < 0 {
		return nil, 0, fmt.Errorf("negative 'offset' argument: %w", apperr.ErrInvalidArgument)
	}

	images, err := s.imageRepository.GetImages(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("get images from repository: %w", err)
	}

	count, err := s.imageRepository.GetImagesCount(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("get images count from repository: %w", err)
	}

	return images, count, nil
}
