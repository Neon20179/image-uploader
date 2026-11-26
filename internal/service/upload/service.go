package upload

import (
	"context"

	"github.com/google/uuid"

	"image-uploader/internal/domain"
	"image-uploader/internal/postgres"
	"image-uploader/internal/s3"
)

type UploadRepository interface {
	MarkFailed(ctx context.Context, id uuid.UUID) error
	MarkProcessFailed(ctx context.Context, id uuid.UUID, process domain.UploadProcess) error
	GetUpload(ctx context.Context, id uuid.UUID) (domain.Upload, error)
	GetUploads(ctx context.Context, limit *int, offset *int) ([]domain.Upload, error)
	GetUploadsCount(ctx context.Context) (int, error)
	DeleteUpload(ctx context.Context, id uuid.UUID) error
	CreateUpload(ctx context.Context, upload domain.Upload) (domain.Upload, error)
	MarkNormalized(ctx context.Context, id uuid.UUID, objectKey string) error
}

type OutboxRepository interface {
	Add(ctx context.Context, event domain.OutboxEvent) (domain.OutboxEvent, error)
}

type NormalizationRepository interface {
	CreateNormalizationResult(ctx context.Context, result domain.NormalizationResult) (domain.NormalizationResult, error)
}

type UploadService struct {
	uploadRepository        UploadRepository
	outboxRepository        OutboxRepository
	normalizationRepository NormalizationRepository
	txManager               *postgres.TxManager
	s3Client                s3.Client
	cfg                     Config
}

func NewUploadService(
	uploadRepository UploadRepository,
	outboxRepository OutboxRepository,
	normalizationRepository NormalizationRepository,
	txManager *postgres.TxManager,
	s3Client s3.Client,
	cfg Config,
) *UploadService {
	return &UploadService{
		uploadRepository:        uploadRepository,
		outboxRepository:        outboxRepository,
		normalizationRepository: normalizationRepository,
		txManager:               txManager,
		s3Client:                s3Client,
		cfg:                     cfg,
	}
}
