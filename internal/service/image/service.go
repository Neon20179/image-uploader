package image

import (
	"context"

	"github.com/google/uuid"

	"image-uploader/internal/domain"
	"image-uploader/internal/s3"
)

type ImageRepository interface {
	GetImage(ctx context.Context, id uuid.UUID) (domain.Image, error)
	GetImages(ctx context.Context, limit *int, offset *int) ([]domain.Image, error)
	GetImagesCount(ctx context.Context) (int, error)
	DeleteImage(ctx context.Context, id uuid.UUID) error
}

type ImageService struct {
	imageRepository ImageRepository
	s3Client        s3.Client
	bucket          string
}

func NewImageService(imageRepository ImageRepository, s3Client s3.Client, bucket string) *ImageService {
	return &ImageService{
		imageRepository: imageRepository,
		s3Client:        s3Client,
		bucket:          bucket,
	}
}
