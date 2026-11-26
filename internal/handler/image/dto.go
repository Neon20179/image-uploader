package image

import (
	"time"

	"github.com/google/uuid"

	"image-uploader/internal/domain"
)

type ImageDTO struct {
	ID uuid.UUID `json:"id"`

	URL string `json:"url"`

	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Size     int    `json:"size"`
	MIMEType string `json:"mimeType"`

	CreatedAt time.Time `json:"createdAt"`
}

func imageDTOFromDomain(image domain.Image, imageURL string) ImageDTO {
	return ImageDTO{
		ID:        image.ID,
		URL:       imageURL,
		Width:     image.Width,
		Height:    image.Height,
		Size:      image.Size,
		MIMEType:  image.MIMEType,
		CreatedAt: image.CreatedAt,
	}
}
