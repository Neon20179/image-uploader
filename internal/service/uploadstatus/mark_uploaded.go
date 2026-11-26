package uploadstatus

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"image-uploader/internal/domain"
	"image-uploader/internal/events"
)

func (s *UploadStatusService) MarkUploaded(ctx context.Context, id uuid.UUID) error {
	return s.txManager.WithTx(ctx, func(ctx context.Context) error {
		upload, err := s.uploadRepository.GetUploadForUpdate(ctx, id)
		if err != nil {
			return err
		}
		if upload.Status != domain.UploadStatusPending {
			return nil
		}

		err = s.uploadRepository.UpdateStatus(ctx, id, domain.UploadStatusPending, domain.UploadStatusUploaded)
		if err != nil {
			return err
		}

		payload, err := json.Marshal(events.ProcessCompleted{UploadID: id})
		if err != nil {
			return fmt.Errorf("marshal image uploaded event: %w", err)
		}
		outboxEvent := domain.NewOutboxEventUninitialized(s.cfg.UploadedTopic, []byte(id.String()), payload, nil)
		if _, err := s.outboxRepository.Add(ctx, outboxEvent); err != nil {
			return fmt.Errorf("add mark uploaded event: %w", err)
		}
		return nil
	})
}
