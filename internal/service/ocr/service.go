package ocr

import (
	"context"
	"io"

	"github.com/google/uuid"

	"image-uploader/internal/domain"
	ocrclient "image-uploader/internal/ocr"
)

type UploadRepository interface {
	MarkFailed(ctx context.Context, id uuid.UUID) error
	MarkProcessFailed(ctx context.Context, id uuid.UUID, process domain.UploadProcess) error
	GetUploadForUpdate(ctx context.Context, id uuid.UUID) (domain.Upload, error)
	MarkOCRCompleted(ctx context.Context, id uuid.UUID) error
	UpdateStatus(ctx context.Context, id uuid.UUID, from domain.UploadStatus, to domain.UploadStatus) error
}

type OCRRepository interface {
	CreateOCRResult(ctx context.Context, result domain.OCRResult) (domain.OCRResult, error)
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

type OCRClient interface {
	RecognizeImage(ctx context.Context, image io.Reader) (ocrclient.Result, error)
}

type OCRService struct {
	uploadRepository UploadRepository
	ocrRepository    OCRRepository
	outboxRepository OutboxRepository
	txManager        TxManager
	s3Client         ObjectStorage
	ocrClient        OCRClient
	cfg              Config
}

func NewOCRService(
	uploadRepository UploadRepository,
	ocrRepository OCRRepository,
	outboxRepository OutboxRepository,
	txManager TxManager,
	s3Client ObjectStorage,
	ocrClient OCRClient,
	cfg Config,
) *OCRService {
	return &OCRService{
		uploadRepository: uploadRepository,
		ocrRepository:    ocrRepository,
		outboxRepository: outboxRepository,
		txManager:        txManager,
		s3Client:         s3Client,
		ocrClient:        ocrClient,
		cfg:              cfg,
	}
}
