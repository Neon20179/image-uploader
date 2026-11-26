package events

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type S3Uploaded struct {
	EventName string `json:"eventName"`
	Key       string `json:"key"`
}

func (ev S3Uploaded) UploadID() (uuid.UUID, error) {
	parts := strings.Split(ev.Key, "/")
	if len(parts) == 0 {
		return uuid.Nil, fmt.Errorf("upload key not found")
	}
	uploadID, err := uuid.Parse(parts[len(parts)-1])
	if err != nil {
		return uuid.Nil, fmt.Errorf("parse upload ID: %w", err)
	}
	return uploadID, nil
}
