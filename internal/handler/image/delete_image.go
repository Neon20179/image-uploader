package image

import (
	"net/http"

	"image-uploader/internal/apperr"
	"image-uploader/internal/http/request"
	"image-uploader/internal/http/response"
	"image-uploader/internal/logger"
)

// DeleteImage godoc
// @Summary Delete an image
// @Tags images
// @Produce json
// @Param id path string true "Image UUID" Format(uuid)
// @Success 204 "Image deleted"
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /images/{id} [delete]
func (h *ImageHandler) DeleteImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	id, err := request.GetUUIDPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(apperr.ErrInvalidArgument, "Failed to get id from path value")
		return
	}

	if err := h.imageService.DeleteImage(ctx, id); err != nil {
		responseHandler.ErrorResponse(err, "Failed to delete image")
		return
	}

	responseHandler.NoContentResponse()
}
