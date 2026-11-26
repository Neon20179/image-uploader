package upload

import (
	"context"
	"fmt"

	"image-uploader/internal/apperr"
	"image-uploader/internal/domain"
)

func (s *UploadService) GetUploads(ctx context.Context, limit *int, offset *int) ([]domain.Upload, int, error) {
	if limit != nil && *limit < 0 {
		return nil, 0, fmt.Errorf("negative 'limit' argument: %w", apperr.ErrInvalidArgument)
	}

	if offset != nil && *offset < 0 {
		return nil, 0, fmt.Errorf("negative 'offset' argument: %w", apperr.ErrInvalidArgument)
	}

	uploads, err := s.uploadRepository.GetUploads(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("get uploads from repository: %w", err)
	}

	count, err := s.uploadRepository.GetUploadsCount(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("get uploads count from repository: %w", err)
	}

	return uploads, count, nil
}
