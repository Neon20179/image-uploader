package embedding

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

type EmbeddingService interface {
	EmbedUpload(ctx context.Context, id uuid.UUID) error
}

type Handler struct {
	embeddingService EmbeddingService
}

func NewHandler(embeddingService EmbeddingService) *Handler {
	return &Handler{embeddingService: embeddingService}
}

func (h *Handler) Handle(ctx context.Context, msg kafka.Message) error {
	var event events.ProcessCompleted
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		return kafka.NewPermanentError(fmt.Errorf("decode: %w", err))
	}

	if err := h.embeddingService.EmbedUpload(ctx, event.UploadID); err != nil {
		if errors.Is(err, apperr.ErrNotFound) ||
			errors.Is(err, apperr.ErrConflict) ||
			errors.Is(err, apperr.ErrInvalidArgument) {
			return kafka.NewPermanentError(err)
		}
		if apperr.IsUnavailableError(err) && !kafka.IsPermanentError(err) && !kafka.IsRetryableError(err) && !kafka.IsWaitError(err) {
			return kafka.NewWaitError(err)
		}
		return err
	}

	return nil
}
