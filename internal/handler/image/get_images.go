package image

import (
	"fmt"
	"net/http"

	"image-uploader/internal/http/request"
	"image-uploader/internal/http/response"
	"image-uploader/internal/logger"
)

type GetImagesResponse struct {
	Images []ImageDTO `json:"images"`
	Count  int        `json:"count"`
}

// GetImages godoc
// @Summary List images
// @Tags images
// @Produce json
// @Param limit query int false "Maximum number of images" minimum(0)
// @Param offset query int false "Number of images to skip" minimum(0)
// @Success 200 {object} GetImagesResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /images [get]
func (h *ImageHandler) GetImages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	limit, offset, err := getLimitOffsetParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "Failed to get limit/offset query params")
		return
	}

	images, count, err := h.imageService.GetImages(ctx, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, "Failed to get image")
		return
	}

	responseItems := make([]ImageDTO, 0, len(images))
	for _, image := range images {
		imageURL := h.urlBuilder.Build(image.Bucket, image.ObjectKey)
		responseItems = append(responseItems, imageDTOFromDomain(image, imageURL))
	}

	response := GetImagesResponse{
		Images: responseItems,
		Count:  count,
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
