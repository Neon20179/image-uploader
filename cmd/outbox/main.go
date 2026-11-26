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
	outboxrepository "image-uploader/internal/repository/outbox"
	"image-uploader/internal/telemetry"
	outboxrelay "image-uploader/internal/workers/outbox"
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

	tracing, stopTracing := telemetry.StartTracing(ctx, log.Logger, "image-uploader-outbox")
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

	outboxRepo := outboxrepository.NewOutboxRepository(pool)

	log.Debug("initializing producer")
	kafkaCfg := franzgo.MustNewConfig()
	producer, err := franzgo.NewProducer(kafkaCfg.Producer(), franzgo.NewTracing(tracing.Provider, tracing.Propagator))
	if err != nil {
		log.Error("failed to initialize producer", slog.Any("error", err))
		os.Exit(1)
	}
	defer func() {
		if err := producer.Close(ctx); err != nil {
			log.Error("closing producer", slog.Any("error", err))
		}
	}()

	log.Debug("initializing outbox relay")
	relayCfg := outboxrelay.MustNewConfig()
	relay := outboxrelay.NewRelay(outboxRepo, txManager, producer, log.Logger, relayCfg)

	if err = relay.Run(ctx); err != nil {
		log.Error("relay run error", slog.Any("error", err))
	}
}
