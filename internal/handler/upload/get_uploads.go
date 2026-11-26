package upload

import (
	"fmt"
	"net/http"

	"image-uploader/internal/http/request"
	"image-uploader/internal/http/response"
	"image-uploader/internal/logger"
)

type GetUploadsResponse struct {
	Uploads []UploadDTO `json:"uploads"`
	Count   int         `json:"count"`
}

// GetUploads godoc
// @Summary List uploads
// @Tags uploads
// @Produce json
// @Param limit query int false "Maximum number of uploads" minimum(0)
// @Param offset query int false "Number of uploads to skip" minimum(0)
// @Success 200 {object} GetUploadsResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /uploads [get]
func (h *UploadHandler) GetUploads(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	limit, offset, err := getLimitOffsetParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "Failed to get limit/offset query params")
		return
	}

	uploads, count, err := h.uploadService.GetUploads(ctx, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, "Failed to get upload")
		return
	}

	responseItems := make([]UploadDTO, 0, len(uploads))
	for _, upload := range uploads {
		uploadURL := h.urlBuilder.Build(upload.Bucket, upload.RawObjectKey)
		responseItems = append(responseItems, uploadDTOFromDomain(upload, uploadURL))
	}

	response := GetUploadsResponse{
		Uploads: responseItems,
		Count:   count,
	}

	responseHandler.JSONResponse(response, http.StatusOK)
}

func getLimitOffsetParams(r *http.Request) (*int, *int, error) {
	limit, err := request.GetIntQueryParam(r, "limit")
	if err != nil {
		return nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
	}

	offset, err := request.GetIntQueryParam(r, "offset")
	if err != nil {
		return nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}

	return limit, offset, nil
}
