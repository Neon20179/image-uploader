package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"image-uploader/internal/kafka/franzgo"
	"image-uploader/internal/logger"
	"image-uploader/internal/ocr"
	"image-uploader/internal/postgres"
	"image-uploader/internal/postgres/pgx"
	ocrrepository "image-uploader/internal/repository/ocr"
	outboxrepository "image-uploader/internal/repository/outbox"
	uploadrepository "image-uploader/internal/repository/upload"
	"image-uploader/internal/s3/minio"
	ocrservice "image-uploader/internal/service/ocr"
	"image-uploader/internal/telemetry"
	ocrworker "image-uploader/internal/workers/ocr"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	logCfg := logger.MustNewConfig()
	log, err := logger.NewLogger(logCfg)
	if err != nil {
		panic(fmt.Errorf("init app logger: %w", err))
	}
	defer func() {
		if err := log.Close(); err != nil {
			slog.Error("failed to close app logger", slog.Any("error", err))
		}
	}()

	tracing, stopTracing := telemetry.StartTracing(ctx, log.Logger, "image-uploader-ocr")
	defer stopTracing()

	log.Debug("initializing postgres connection pool")
	postgresCfg := pgx.MustNewConfig()
	pool, err := pgx.NewPool(
		ctx,
		postgresCfg,
	)
	if err != nil {
		log.Error("failed to init postgres connection pool", slog.Any("error", err))
		os.Exit(1)
	}
	txManager := postgres.NewTxManager(pool)
	defer pool.Close()

	log.Debug("initializing MinIO")
	s3ClientCfg := minio.MustNewConfig()
	s3Client, err := minio.NewClient(ctx, s3ClientCfg)
	if err != nil {
		log.Error("failed to init minio client", slog.Any("error", err))
		os.Exit(1)
	}

	outboxRepo := tracing.TraceOutboxWrites(outboxrepository.NewOutboxRepository(pool))
	uploadRepo := uploadrepository.NewUploadRepository(pool)
	ocrRepo := ocrrepository.NewOCRRepository(pool)
	ocrSvcCfg := ocrservice.MustNewConfig()
	ocrClient, err := ocr.NewClient(ocrSvcCfg.OCR.Options())
	if err != nil {
		log.Error("failed to init OCR client", slog.Any("error", err))
		os.Exit(1)
	}
	ocrSvc := ocrservice.NewOCRService(uploadRepo, ocrRepo, outboxRepo, txManager, s3Client, ocrClient, ocrSvcCfg)

	log.Debug("initializing consumer")
	handler := ocrworker.NewHandler(ocrSvc)
	kafkaCfg := franzgo.MustNewConfig()
	workerCfg := ocrworker.MustNewConfig()
	consumer, err := franzgo.NewConsumer(
		kafkaCfg.Consumer(workerCfg.Group, []string{workerCfg.Topic}, workerCfg.DLQTopic),
		log.Logger,
		handler,
		ocrworker.NewDeadLetterHandler(ocrSvc),
		franzgo.NewTracing(tracing.Provider, tracing.Propagator),
	)
	if err != nil {
		log.Error("failed to initialize consumer", slog.Any("error", err))
		os.Exit(1)
	}

	log.Debug("running consumer")
	if err = consumer.Run(ctx); err != nil {
		log.Error("consumer run error", slog.Any("error", err))
	}
}
