package image

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"image-uploader/internal/domain"
	"image-uploader/internal/http/server"
)

type ImageService interface {
	GetImage(ctx context.Context, id uuid.UUID) (domain.Image, error)
	GetImages(ctx context.Context, limit *int, offset *int) ([]domain.Image, int, error)
	DeleteImage(ctx context.Context, id uuid.UUID) error
}

type ImageURLBuilder interface {
	Build(bucket string, objectKey string) string
}

type ImageHandler struct {
	imageService ImageService
	urlBuilder   ImageURLBuilder
}

func NewImageHandler(imageService ImageService, urlBuilder ImageURLBuilder) *ImageHandler {
	return &ImageHandler{
		imageService: imageService,
		urlBuilder:   urlBuilder,
	}
}

func (h *ImageHandler) Routes() []server.Route {
	return []server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/images",
			Handler: h.GetImages,
		},
		{
			Method:  http.MethodGet,
			Path:    "/images/{id}",
			Handler: h.GetImage,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/images/{id}",
			Handler: h.DeleteImage,
		},
	}
}
