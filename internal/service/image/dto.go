package image

import "image-uploader/internal/domain"

type ImageUpload struct {
	Image     domain.Image
	UploadURL string
	Fields    map[string]string
}
