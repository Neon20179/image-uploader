package events

import "github.com/google/uuid"

type ProcessCompleted struct {
	UploadID uuid.UUID `json:"uploadId"`
}
