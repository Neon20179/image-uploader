package embedding

import (
	"context"
	"io"

	"github.com/google/uuid"

	"image-uploader/internal/clip"
	"image-uploader/internal/domain"
)

type UploadRepository interface {
	MarkFailed(ctx context.Context, id uuid.UUID) error
	MarkProcessFailed(ctx context.Context, id uuid.UUID, process domain.UploadProcess) error
	GetUploadForUpdate(ctx context.Context, id uuid.UUID) (domain.Upload, error)
	MarkEmbedded(ctx context.Context, id uuid.UUID) error
	UpdateStatus(ctx context.Context, id uuid.UUID, from domain.UploadStatus, to domain.UploadStatus) error
}

type EmbeddingRepository interface {
	CreateEmbedding(ctx context.Context, embedding domain.Embedding) (domain.Embedding, error)
}

type EmbeddingBlacklistRepository interface {
	HasSimilarEmbedding(ctx context.Context, vector []float32, similarityThreshold float64) (bool, error)
}

type OutboxRepository interface {
	Add(ctx context.Context, event domain.OutboxEvent) (domain.OutboxEvent, error)
}

type TxManager interface {
	WithTx(ctx context.Context, fn func(context.Context) error) error
}

type ObjectStorage interface {
	GetObject(ctx context.Context, bucket, objectKey string) (io.ReadCloser, error)
}

type CLIPClient interface {
	EncodeImage(ctx context.Context, image io.Reader) (clip.Result, error)
}

type EmbeddingService struct {
	uploadRepository    UploadRepository
	embeddingRepository EmbeddingRepository
	blacklistRepository EmbeddingBlacklistRepository
	outboxRepository    OutboxRepository
	txManager           TxManager
	s3Client            ObjectStorage
	clipClient          CLIPClient
	cfg                 Config
}

func NewEmbeddingService(
	uploadRepository UploadRepository,
	embeddingRepository EmbeddingRepository,
	blacklistRepository EmbeddingBlacklistRepository,
	outboxRepository OutboxRepository,
	txManager TxManager,
	s3Client ObjectStorage,
	clipClient CLIPClient,
	cfg Config,
) *EmbeddingService {
	return &EmbeddingService{
		uploadRepository:    uploadRepository,
		embeddingRepository: embeddingRepository,
		blacklistRepository: blacklistRepository,
		outboxRepository:    outboxRepository,
		txManager:           txManager,
		s3Client:            s3Client,
		clipClient:          clipClient,
		cfg:                 cfg,
	}
}
