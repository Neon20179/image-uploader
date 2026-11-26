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
	"image-uploader/internal/postgres"
	"image-uploader/internal/postgres/pgx"
	embeddingrepository "image-uploader/internal/repository/embedding"
	imagerepository "image-uploader/internal/repository/image"
	normalizationrepository "image-uploader/internal/repository/normalization"
	ocrrepository "image-uploader/internal/repository/ocr"
	outboxrepository "image-uploader/internal/repository/outbox"
	uploadrepository "image-uploader/internal/repository/upload"
	uploadstatusservice "image-uploader/internal/service/uploadstatus"
	"image-uploader/internal/telemetry"
	"image-uploader/internal/workers/finalizer"
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

	tracing, stopTracing := telemetry.StartTracing(ctx, log.Logger, "image-uploader-finalizer")
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

	outboxRepo := tracing.TraceOutboxWrites(outboxrepository.NewOutboxRepository(pool))
	uploadRepo := uploadrepository.NewUploadRepository(pool)
	ocrRepo := ocrrepository.NewOCRRepository(pool)
	embeddingRepo := embeddingrepository.NewEmbeddingRepository(pool)
	normalizationRepo := normalizationrepository.NewNormalizationRepository(pool)
	imageRepo := imagerepository.NewImageRepository(pool)
	uploadSvcConfig := uploadstatusservice.MustNewConfig()
	uploadSvc := uploadstatusservice.NewUploadStatusService(uploadRepo, outboxRepo, ocrRepo, embeddingRepo, normalizationRepo, imageRepo, txManager, uploadSvcConfig)

	log.Debug("initializing consumer")
	kafkaCfg := franzgo.MustNewConfig()
	handler := finalizer.NewHandler(log.Logger, uploadSvc)
	workerCfg := finalizer.MustNewConfig()
	consumerCfg := kafkaCfg.Consumer(workerCfg.Group, []string{workerCfg.EmbeddingTopic, workerCfg.OCRTopic}, workerCfg.DLQTopic)
	consumer, err := franzgo.NewConsumer(
		consumerCfg,
		log.Logger,
		handler,
		finalizer.NewDeadLetterHandler(uploadSvc),
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
