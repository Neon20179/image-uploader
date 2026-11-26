package embedding

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"image-uploader/internal/apperr"
	"image-uploader/internal/domain"
	"image-uploader/internal/events"
)

func (s *EmbeddingService) EmbedUpload(ctx context.Context, id uuid.UUID) error {
	return s.txManager.WithTx(ctx, func(ctx context.Context) error {
		upload, err := s.uploadRepository.GetUploadForUpdate(ctx, id)
		if err != nil {
			return fmt.Errorf("get upload: %w", err)
		}
		if upload.EmbeddingStatus == domain.ProcessStatusCompleted ||
			upload.EmbeddingStatus == domain.ProcessStatusFailed {
			return nil
		}
		if upload.NormalizationStatus != domain.ProcessStatusCompleted ||
			upload.ObjectKey == nil || *upload.ObjectKey == "" {
			return fmt.Errorf("upload is not normalized: %w", apperr.ErrConflict)
		}

		image, err := s.s3Client.GetObject(ctx, upload.Bucket, *upload.ObjectKey)
		if err != nil {
			return fmt.Errorf("get normalized image: %w", err)
		}
		defer func() {
			_ = image.Close()
		}()

		result, err := s.clipClient.EncodeImage(ctx, image)
		if err != nil {
			return fmt.Errorf("encode image: %w", err)
		}

		blockDetected, err := s.blacklistRepository.HasSimilarEmbedding(ctx, result.Vector, s.cfg.BlacklistSimilarityThreshold)
		if err != nil {
			return fmt.Errorf("check embedding blacklist: %w", err)
		}

		embedding := domain.NewEmbeddingUninitialized(id, result.Vector, blockDetected)
		if _, err := s.embeddingRepository.CreateEmbedding(ctx, embedding); err != nil {
			return fmt.Errorf("save embedding: %w", err)
		}
		if err := s.uploadRepository.MarkEmbedded(ctx, id); err != nil {
			return fmt.Errorf("mark embedded: %w", err)
		}
		if blockDetected {
			if err := s.uploadRepository.UpdateStatus(ctx, id, upload.Status, domain.UploadStatusBlocked); err != nil {
				return fmt.Errorf("block upload: %w", err)
			}
		}

		payload, err := json.Marshal(events.ProcessCompleted{UploadID: id})
		if err != nil {
			return fmt.Errorf("marshal image embedded event: %w", err)
		}
		event := domain.NewOutboxEventUninitialized(s.cfg.EmbeddingTopic, []byte(id.String()), payload, nil)
		if _, err := s.outboxRepository.Add(ctx, event); err != nil {
			return fmt.Errorf("add image embedded event: %w", err)
		}
		return nil
	})
}
