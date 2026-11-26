package upload

import (
	"net/http"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"image-uploader/internal/domain"
	"image-uploader/internal/http/response"
	"image-uploader/internal/logger"
)

type CreateUploadResponse struct {
	UploadID  uuid.UUID           `json:"uploadId"`
	Status    domain.UploadStatus `json:"status"`
	UploadURL string              `json:"uploadUrl"`
	Fields    map[string]string   `json:"fields"`
}

// CreateUpload godoc
// @Summary Create an upload
// @Tags uploads
// @Produce json
// @Success 201 {object} CreateUploadResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /uploads [post]
func (h *UploadHandler) CreateUpload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	upload, err := h.uploadService.CreateUpload(ctx)
	if err != nil {
		responseHandler.ErrorResponse(err, "Failed to create upload")
		return
	}

	trace.SpanFromContext(ctx).SetAttributes(attribute.String("upload.id", upload.UploadID.String()))
	responseBody := CreateUploadResponse{
		UploadID:  upload.UploadID,
		Status:    upload.Status,
		UploadURL: upload.UploadURL,
		Fields:    upload.Fields,
	}

	responseHandler.JSONResponse(responseBody, http.StatusCreated)
}
