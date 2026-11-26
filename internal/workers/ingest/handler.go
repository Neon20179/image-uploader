package ingest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"image-uploader/internal/apperr"
	"image-uploader/internal/events"
	"image-uploader/internal/kafka"
)

type UploadStatusService interface {
	MarkUploaded(ctx context.Context, id uuid.UUID) error
}

type Handler struct {
	log                 *slog.Logger
	uploadStatusService UploadStatusService
}

func NewHandler(log *slog.Logger, uploadStatusService UploadStatusService) *Handler {
	return &Handler{
		log:                 log,
		uploadStatusService: uploadStatusService,
	}
}

func (h *Handler) Handle(ctx context.Context, msg kafka.Message) error {
	var event events.S3Uploaded
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		return kafka.NewPermanentError(fmt.Errorf("decode: %w", err))
	}
	uploadID, err := event.UploadID()
	if err != nil {
		return kafka.NewPermanentError(err)
	}

	trace.SpanFromContext(ctx).SetAttributes(attribute.String("upload.id", uploadID.String()))
	if err := h.uploadStatusService.MarkUploaded(ctx, uploadID); err != nil {
		if errors.Is(err, apperr.ErrNotFound) || errors.Is(err, apperr.ErrConflict) {
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
