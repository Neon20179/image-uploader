package image

import (
	"net/http"

	"image-uploader/internal/apperr"
	"image-uploader/internal/http/request"
	"image-uploader/internal/http/response"
	"image-uploader/internal/logger"
)

// GetImage godoc
// @Summary Get an image
// @Tags images
// @Produce json
// @Param id path string true "Image UUID" Format(uuid)
// @Success 200 {object} ImageDTO
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /images/{id} [get]
func (h *ImageHandler) GetImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	id, err := request.GetUUIDPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(apperr.ErrInvalidArgument, "Failed to get image id from path value")
		return
	}

	image, err := h.imageService.GetImage(ctx, id)
	if err != nil {
		responseHandler.ErrorResponse(err, "Failed to get image")
		return
	}

	imageURL := h.urlBuilder.Build(image.Bucket, image.ObjectKey)
	responseBody := imageDTOFromDomain(image, imageURL)

	responseHandler.JSONResponse(responseBody, http.StatusOK)
}
