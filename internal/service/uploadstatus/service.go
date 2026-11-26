package uploadstatus

import (
	"context"

	"github.com/google/uuid"

	"image-uploader/internal/domain"
	"image-uploader/internal/postgres"
)

type OutboxRepository interface {
	Add(ctx context.Context, event domain.OutboxEvent) (domain.OutboxEvent, error)
}

type UploadRepository interface {
	MarkFailed(ctx context.Context, id uuid.UUID) error
	GetUpload(ctx context.Context, id uuid.UUID) (domain.Upload, error)
	GetUploadForUpdate(ctx context.Context, id uuid.UUID) (domain.Upload, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, from domain.UploadStatus, to domain.UploadStatus) error
	MarkReady(ctx context.Context, id uuid.UUID, imageID uuid.UUID) error
}

type ImageRepository interface {
	CreateImage(ctx context.Context, image domain.Image) (domain.Image, error)
}

type OCRRepository interface {
	GetOCRResult(ctx context.Context, uploadID uuid.UUID) (domain.OCRResult, error)
}

type EmbeddingRepository interface {
	GetEmbedding(ctx context.Context, uploadID uuid.UUID) (domain.Embedding, error)
}

type NormalizationRepository interface {
	GetNormalizationResult(ctx context.Context, uploadID uuid.UUID) (domain.NormalizationResult, error)
}

type UploadStatusService struct {
	uploadRepository        UploadRepository
	outboxRepository        OutboxRepository
	ocrRepository           OCRRepository
	embeddingRepository     EmbeddingRepository
	normalizationRepository NormalizationRepository
	imageRepository         ImageRepository
	txManager               *postgres.TxManager
	cfg                     Config
}

func NewUploadStatusService(
	uploadRepository UploadRepository,
	outboxRepository OutboxRepository,
	ocrRepository OCRRepository,
	embeddingRepository EmbeddingRepository,
	normalizationRepository NormalizationRepository,
	imageRepository ImageRepository,
	txManager *postgres.TxManager,
	cfg Config,
) *UploadStatusService {
	return &UploadStatusService{
		uploadRepository:        uploadRepository,
		outboxRepository:        outboxRepository,
		ocrRepository:           ocrRepository,
		embeddingRepository:     embeddingRepository,
		normalizationRepository: normalizationRepository,
		imageRepository:         imageRepository,
		txManager:               txManager,
		cfg:                     cfg,
	}
}
