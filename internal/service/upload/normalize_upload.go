package upload

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/google/uuid"

	"image-uploader/internal/apperr"
	"image-uploader/internal/domain"
	"image-uploader/internal/events"
	"image-uploader/internal/imageproc"
)

func (s *UploadService) NormalizeUpload(ctx context.Context, id uuid.UUID) error {
	return s.txManager.WithTx(ctx, func(ctx context.Context) error {
		upload, err := s.uploadRepository.GetUpload(ctx, id)
		if err != nil {
			return fmt.Errorf("get upload: %w", err)
		}
		if upload.NormalizationStatus != domain.ProcessStatusPending {
			return nil
		}

		objectInfo, err := s.s3Client.StatObject(
			ctx,
			upload.Bucket,
			upload.RawObjectKey,
		)
		if err != nil {
			return fmt.Errorf("get object stats: %w", err)
		}
		if objectInfo.Size > s.cfg.MaxRawSizeBytes {
			return fmt.Errorf("object is too big: %w", apperr.ErrInvalidArgument)
		}

		contentType, err := s.s3Client.ObjectContentType(
			ctx,
			upload.Bucket,
			upload.RawObjectKey,
		)
		if err != nil {
			return fmt.Errorf("object content type: %w", err)
		}
		if !slices.Contains(s.cfg.AllowedMIMETypes, contentType) {
			return fmt.Errorf("unsupported content type: %s: %w", contentType, apperr.ErrInvalidArgument)
		}

		rawObject, err := s.s3Client.GetObject(ctx, upload.Bucket, upload.RawObjectKey)
		if err != nil {
			return fmt.Errorf("get object: %w", err)
		}
		defer func() {
			_ = rawObject.Close()
		}()

		normalizeCtx, cancel := context.WithTimeout(ctx, s.cfg.NormalizeTimeout)
		defer cancel()

		normalizedImage, err := imageproc.Normalize(
			normalizeCtx,
			rawObject,
			s.cfg.Image.Options(),
		)
		if err != nil {
			return fmt.Errorf("normalize: %w", err)
		}

		normalizedObjectKey := fmt.Sprintf("%s/%s.webp", s.cfg.NormalizedObjectKeyPrefix, upload.ID.String())
		normalizedSize := normalizedImage.Buffer.Len()

		if err := s.s3Client.PutObject(
			ctx,
			upload.Bucket,
			normalizedObjectKey,
			normalizedImage.Buffer,
			int64(normalizedSize),
			"image/webp",
		); err != nil {
			return fmt.Errorf("put object: %w", err)
		}

		err = s.uploadRepository.MarkNormalized(
			ctx,
			id,
			normalizedObjectKey,
		)
		if err != nil {
			return err
		}

		result := domain.NewNormalizationResultUninitialized(
			id,
			normalizedImage.Width,
			normalizedImage.Height,
			normalizedSize,
			"image/webp",
		)
		if _, err := s.normalizationRepository.CreateNormalizationResult(ctx, result); err != nil {
			return fmt.Errorf("save normalization result: %w", err)
		}

		payload, err := json.Marshal(events.ProcessCompleted{UploadID: id})
		if err != nil {
			return fmt.Errorf("marshal image normalized event: %w", err)
		}
		outboxEvent := domain.NewOutboxEventUninitialized(s.cfg.NormalizedTopic, []byte(id.String()), payload, nil)
		_, err = s.outboxRepository.Add(ctx, outboxEvent)
		return err
	})
}
