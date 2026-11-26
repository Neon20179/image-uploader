package outbox

import (
	"context"
	"log/slog"
	"time"

	"image-uploader/internal/domain"
	"image-uploader/internal/kafka"
	"image-uploader/internal/postgres"
)

type OutboxRepository interface {
	GetUnpublished(ctx context.Context, limit int) ([]domain.OutboxEvent, error)
	MarkPublished(ctx context.Context, publishedEventIDs []int64) error
}

type Relay struct {
	outboxRepository OutboxRepository
	txManager        *postgres.TxManager
	producer         kafka.Producer
	log              *slog.Logger
	batchSize        int
	pollInterval     time.Duration
}

func NewRelay(
	outboxRepository OutboxRepository,
	txManager *postgres.TxManager,
	producer kafka.Producer,
	log *slog.Logger,
	cfg Config,
) *Relay {
	return &Relay{
		outboxRepository: outboxRepository,
		txManager:        txManager,
		producer:         producer,
		log:              log,
		batchSize:        cfg.BatchSize,
		pollInterval:     cfg.PollInterval,
	}
}

func (r *Relay) Run(ctx context.Context) error {
	ticker := time.NewTicker(r.pollInterval)
	defer ticker.Stop()

	for {
		for {
			publishedCount, err := r.publishBatch(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return nil
				}
				r.log.Error("outbox relay", slog.Any("error", err))
				break
			}
			if publishedCount < r.batchSize {
				break
			}
		}

		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
		}
	}
}

func (r *Relay) publishBatch(ctx context.Context) (int, error) {
	var publishedCount int
	err := r.txManager.WithTx(ctx, func(ctx context.Context) error {
		events, err := r.outboxRepository.GetUnpublished(ctx, r.batchSize)
		if err != nil {
			return err
		}
		if len(events) == 0 {
			return nil
		}

		var publishedEventIDs []int64
		var sendErr error
		for _, event := range events {
			sendErr = r.producer.Send(ctx, kafka.Message{
				Topic: event.Topic, Key: event.Key, Value: event.Payload, Headers: event.Headers,
			})
			if sendErr != nil {
				break
			}
			publishedEventIDs = append(publishedEventIDs, event.ID)
		}

		if len(publishedEventIDs) > 0 {
			if err := r.outboxRepository.MarkPublished(ctx, publishedEventIDs); err != nil {
				return err
			}
		}
		publishedCount = len(publishedEventIDs)
		return sendErr
	})
	return publishedCount, err
}
