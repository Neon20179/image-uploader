package ocr

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"image-uploader/internal/apperr"
	"image-uploader/internal/events"
	"image-uploader/internal/kafka"
)

type OCRService interface {
	RecognizeUpload(ctx context.Context, id uuid.UUID) error
}

type Handler struct {
	ocrService OCRService
}

func NewHandler(ocrService OCRService) *Handler {
	return &Handler{ocrService: ocrService}
}

func (h *Handler) Handle(ctx context.Context, msg kafka.Message) error {
	var event events.ProcessCompleted
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		return kafka.NewPermanentError(fmt.Errorf("decode: %w", err))
	}

	if err := h.ocrService.RecognizeUpload(ctx, event.UploadID); err != nil {
		if errors.Is(err, apperr.ErrNotFound) ||
			errors.Is(err, apperr.ErrConflict) ||
			errors.Is(err, apperr.ErrInvalidArgument) {
			return kafka.NewPermanentError(err)
		}
		if apperr.IsUnavailableError(err) && !errors.Is(err, context.Canceled) &&
			!kafka.IsPermanentError(err) && !kafka.IsRetryableError(err) && !kafka.IsWaitError(err) {
			return kafka.NewWaitError(err)
		}
		return err
	}

	return nil
}
