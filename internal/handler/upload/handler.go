package upload

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"image-uploader/internal/domain"
	"image-uploader/internal/http/server"
	service "image-uploader/internal/service/upload"
)

type UploadService interface {
	GetUpload(ctx context.Context, id uuid.UUID) (domain.Upload, error)
	GetUploads(ctx context.Context, limit *int, offset *int) ([]domain.Upload, int, error)
	DeleteUpload(ctx context.Context, id uuid.UUID) error
	CreateUpload(ctx context.Context) (service.CreateUploadResult, error)
}

type UploadURLBuilder interface {
	Build(bucket string, objectKey string) string
}

type UploadHandler struct {
	uploadService UploadService
	urlBuilder    UploadURLBuilder
}

func NewUploadHandler(uploadService UploadService, urlBuilder UploadURLBuilder) *UploadHandler {
	return &UploadHandler{
		uploadService: uploadService,
		urlBuilder:    urlBuilder,
	}
}

func (h *UploadHandler) Routes() []server.Route {
	return []server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/uploads",
			Handler: h.GetUploads,
		},
		{
			Method:  http.MethodGet,
			Path:    "/uploads/{id}",
			Handler: h.GetUpload,
		},
		{
			Method:  http.MethodPost,
			Path:    "/uploads",
			Handler: h.CreateUpload,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/uploads/{id}",
			Handler: h.DeleteUpload,
		},
	}
}
