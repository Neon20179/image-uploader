package image

import (
	"time"

	"github.com/google/uuid"

	"image-uploader/internal/domain"
)

type ImageModel struct {
	ID uuid.UUID

	Bucket    string
	ObjectKey string

	Width    int
	Height   int
	MIMEType string
	Size     int

	CreatedAt time.Time
}

func imageDomainFromModel(model ImageModel) domain.Image {
	return domain.NewImage(
		model.ID,
		model.Bucket,
		model.ObjectKey,
		model.Width,
		model.Height,
		model.MIMEType,
		model.Size,
		model.CreatedAt,
	)
}
