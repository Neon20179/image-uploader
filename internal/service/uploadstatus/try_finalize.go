package uploadstatus

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"image-uploader/internal/domain"
)

func (s *UploadStatusService) TryFinalize(ctx context.Context, id uuid.UUID) error {
	return s.txManager.WithTx(ctx, func(ctx context.Context) error {
		upload, err := s.uploadRepository.GetUploadForUpdate(ctx, id)
		if err != nil {
			return err
		}
		if upload.Status != domain.UploadStatusProcessing {
			return nil
		}
		if upload.EmbeddingStatus == domain.ProcessStatusPending ||
			upload.OCRStatus == domain.ProcessStatusPending ||
			upload.NormalizationStatus == domain.ProcessStatusPending {
			return nil
		}

		if upload.NormalizationStatus == domain.ProcessStatusFailed ||
			upload.OCRStatus == domain.ProcessStatusFailed ||
			upload.EmbeddingStatus == domain.ProcessStatusFailed {
			return s.uploadRepository.UpdateStatus(ctx, id, domain.UploadStatusProcessing, domain.UploadStatusFailed)
		}

		ocr, err := s.ocrRepository.GetOCRResult(ctx, id)
		if err != nil {
			return fmt.Errorf("get OCR result: %w", err)
		}
		embedding, err := s.embeddingRepository.GetEmbedding(ctx, id)
		if err != nil {
			return fmt.Errorf("get embedding: %w", err)
		}

		if ocr.BlockDetected || embedding.BlockDetected {
			return s.uploadRepository.UpdateStatus(ctx, id, domain.UploadStatusProcessing, domain.UploadStatusBlocked)
		}

		normalizationResult, err := s.normalizationRepository.GetNormalizationResult(ctx, id)
		if err != nil {
			return fmt.Errorf("get normalization result: %w", err)
		}

		imageID, err := uuid.NewV7()
		if err != nil {
			return fmt.Errorf("generating uuid: %w", err)
		}

		image := domain.NewImageUninitialized(
			imageID,
			upload.Bucket,
			*upload.ObjectKey,
			normalizationResult.Width,
			normalizationResult.Height,
			normalizationResult.MIMEType,
			normalizationResult.Size,
		)
		if _, err := s.imageRepository.CreateImage(ctx, image); err != nil {
			return fmt.Errorf("create image: %w", err)
		}

		return s.uploadRepository.UpdateStatus(ctx, id, domain.UploadStatusProcessing, domain.UploadStatusReady)
	})
}
