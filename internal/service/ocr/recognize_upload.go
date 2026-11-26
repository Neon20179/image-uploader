package ocr

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"image-uploader/internal/apperr"
	"image-uploader/internal/domain"
	"image-uploader/internal/events"
)

func (s *OCRService) RecognizeUpload(ctx context.Context, id uuid.UUID) error {
	return s.txManager.WithTx(ctx, func(ctx context.Context) error {
		upload, err := s.uploadRepository.GetUploadForUpdate(ctx, id)
		if err != nil {
			return fmt.Errorf("get upload: %w", err)
		}
		if upload.OCRStatus == domain.ProcessStatusCompleted ||
			upload.OCRStatus == domain.ProcessStatusFailed {
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

		result, err := s.ocrClient.RecognizeImage(ctx, image)
		if err != nil {
			return fmt.Errorf("recognize image: %w", err)
		}

		blockDetected := hasBlacklistedWord(result.Text, s.cfg.BlacklistWords, s.cfg.BlacklistSimilarityThreshold)

		ocrResult := domain.NewOCRResultUninitialized(id, result.Text, blockDetected)
		if _, err := s.ocrRepository.CreateOCRResult(ctx, ocrResult); err != nil {
			return fmt.Errorf("save OCR result: %w", err)
		}
		if err := s.uploadRepository.MarkOCRCompleted(ctx, id); err != nil {
			return fmt.Errorf("mark OCR completed: %w", err)
		}
		if blockDetected {
			if err := s.uploadRepository.UpdateStatus(ctx, id, upload.Status, domain.UploadStatusBlocked); err != nil {
				return fmt.Errorf("block upload: %w", err)
			}
		}

		payload, err := json.Marshal(events.ProcessCompleted{UploadID: id})
		if err != nil {
			return fmt.Errorf("marshal OCR completed event: %w", err)
		}
		event := domain.NewOutboxEventUninitialized(s.cfg.OCRTopic, []byte(id.String()), payload, nil)
		if _, err := s.outboxRepository.Add(ctx, event); err != nil {
			return fmt.Errorf("add OCR completed event: %w", err)
		}
		return nil
	})
}
