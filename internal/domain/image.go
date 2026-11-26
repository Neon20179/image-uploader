package domain

import (
	"time"

	"github.com/google/uuid"
)

type Image struct {
	ID uuid.UUID

	Bucket    string
	ObjectKey string

	Width    int
	Height   int
	Size     int
	MIMEType string

	CreatedAt time.Time
}

func NewImage(
	id uuid.UUID,
	bucket string,
	objectKey string,
	width int,
	height int,
	mimeType string,
	size int,
	createdAt time.Time,
) Image {
	return Image{
		ID:        id,
		Bucket:    bucket,
		ObjectKey: objectKey,
		Width:     width,
		Height:    height,
		MIMEType:  mimeType,
		Size:      size,
		CreatedAt: createdAt,
	}
}

func NewImageUninitialized(
	id uuid.UUID,
	bucket string,
	objectKey string,
	width int,
	height int,
	mimeType string,
	size int,
) Image {
	return NewImage(
		id,
		bucket,
		objectKey,
		width,
		height,
		mimeType,
		size,
		time.Now(),
	)
}
