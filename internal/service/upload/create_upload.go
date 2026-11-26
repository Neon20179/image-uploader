package upload

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"image-uploader/internal/domain"
	"image-uploader/internal/s3"
)

func (s *UploadService) CreateUpload(ctx context.Context) (CreateUploadResult, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return CreateUploadResult{}, fmt.Errorf("generating uuid: %w", err)
	}

	rawObjectKey := fmt.Sprintf("%s/%s", s.cfg.RawObjectKeyPrefix, id.String())

	upload := domain.NewUploadUninitialized(id, s.cfg.Bucket, rawObjectKey)

	uploadURL, formFields, err := s.s3Client.PresignedPostPolicy(ctx, s.cfg.Bucket, rawObjectKey, s3.PostPolicyOptions{
		Expiry:            s.cfg.URLExpiry,
		MinSizeBytes:      s.cfg.MinRawSizeBytes,
		MaxSizeBytes:      s.cfg.MaxRawSizeBytes,
		ContentTypePrefix: "image/",
	})
	if err != nil {
		return CreateUploadResult{}, err
	}

	upload, err = s.uploadRepository.CreateUpload(ctx, upload)
	if err != nil {
		return CreateUploadResult{}, fmt.Errorf("create upload: %w", err)
	}

	result := CreateUploadResult{
		UploadID:  upload.ID,
		Status:    upload.Status,
		UploadURL: uploadURL.String(),
		Fields:    formFields,
	}

	return result, nil
}
