package upload

import (
	"net/http"

	"image-uploader/internal/apperr"
	"image-uploader/internal/http/request"
	"image-uploader/internal/http/response"
	"image-uploader/internal/logger"
)

// GetUpload godoc
// @Summary Get an upload
// @Tags uploads
// @Produce json
// @Param id path string true "Upload UUID" Format(uuid)
// @Success 200 {object} UploadDTO
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /uploads/{id} [get]
func (h *UploadHandler) GetUpload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	id, err := request.GetUUIDPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(apperr.ErrInvalidArgument, "Failed to get upload id from path value")
		return
	}

	upload, err := h.uploadService.GetUpload(ctx, id)
	if err != nil {
		responseHandler.ErrorResponse(err, "Failed to get upload")
		return
	}

	uploadURL := h.urlBuilder.Build(upload.Bucket, upload.RawObjectKey)
	responseBody := uploadDTOFromDomain(upload, uploadURL)

	responseHandler.JSONResponse(responseBody, http.StatusOK)
}
