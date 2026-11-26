package ocr

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"

	"image-uploader/internal/apperr"
	"image-uploader/internal/events"
	"image-uploader/internal/kafka"
)

type FailureService interface {
	FailRecognition(context.Context, uuid.UUID) error
}

type DeadLetterHandler struct {
	failureService FailureService
}

func NewDeadLetterHandler(failureService FailureService) *DeadLetterHandler {
	return &DeadLetterHandler{failureService: failureService}
}

func (h *DeadLetterHandler) OnDeadLetter(ctx context.Context, msg kafka.Message) error {
	var event events.ProcessCompleted
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		return nil
	}

	err := h.failureService.FailRecognition(ctx, event.UploadID)
	if apperr.IsUnavailableError(err) && !errors.Is(err, context.Canceled) &&
		!kafka.IsPermanentError(err) && !kafka.IsRetryableError(err) && !kafka.IsWaitError(err) {
		return kafka.NewWaitError(err)
	}

	return err
}
