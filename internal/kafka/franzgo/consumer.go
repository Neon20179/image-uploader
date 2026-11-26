package franzgo

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"

	"image-uploader/internal/kafka"
)

type ConsumerConfig struct {
	Broker          string
	Group           string
	Topics          []string
	MaxAttempts     int
	RetryTimeout    time.Duration
	WaitTimeout     time.Duration
	DeliveryTimeout time.Duration
	DLQTopic        string
	ShutdownTimeout time.Duration
}

type Consumer struct {
	client            *kgo.Client
	handler           kafka.ConsumerHandler
	deadLetterHandler kafka.DeadLetterHandler
	log               *slog.Logger
	cfg               ConsumerConfig
	tracing           Tracing
}

func NewConsumer(cfg ConsumerConfig, log *slog.Logger, handler kafka.ConsumerHandler, deadLetterHandler kafka.DeadLetterHandler, tracing Tracing) (*Consumer, error) {
	if cfg.MaxAttempts <= 0 || cfg.RetryTimeout <= 0 || cfg.WaitTimeout <= 0 || cfg.DeliveryTimeout <= 0 || cfg.ShutdownTimeout <= 0 {
		return nil, fmt.Errorf("consumer attempts and timeouts must be positive")
	}
	client, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.Broker),
		kgo.AllowAutoTopicCreation(),
		kgo.ConsumerGroup(cfg.Group),
		kgo.ConsumeTopics(cfg.Topics...),
		kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
		kgo.AutoCommitMarks(),
		kgo.BlockRebalanceOnPoll(),
		kgo.OnPartitionsRevoked(
			func(ctx context.Context, client *kgo.Client, _ map[string][]int32) {
				if err := client.CommitMarkedOffsets(ctx); err != nil {
					log.Error("commit offsets on partition rebalance", slog.Any("error", err))
				}
			}),
	)
	if err != nil {
		return nil, fmt.Errorf("initializing consumer: %w", err)
	}

	return &Consumer{
		client:            client,
		handler:           handler,
		deadLetterHandler: deadLetterHandler,
		log:               log,
		cfg:               cfg,
		tracing:           tracing,
	}, nil
}

func (c *Consumer) Run(ctx context.Context) error {
	c.log.Info(
		"consumer started",
		slog.String("topic", strings.Join(c.cfg.Topics, ", ")),
	)

	err := c.loop(ctx)

	c.shutdown(ctx)

	if err != nil {
		return fmt.Errorf("consumer error: %w", err)
	}

	return nil
}

func (c *Consumer) loop(ctx context.Context) error {
	for {
		fetches := c.client.PollFetches(ctx)
		if fetches.IsClientClosed() || ctx.Err() != nil {
			return nil
		}

		fetches.EachError(func(t string, p int32, err error) {
			c.log.Error("consumer fetch", slog.String("topic", t), slog.Any("partition", p), slog.Any("error", err))
		})

		err := c.processBatch(ctx, fetches)
		c.client.AllowRebalance()
		if err != nil {
			return err // Could't send message to DLQ
		}
	}
}

func (c *Consumer) processBatch(ctx context.Context, fetches kgo.Fetches) error {
	iter := fetches.RecordIter()
	for !iter.Done() {
		if ctx.Err() != nil {
			return nil
		}

		if err := c.processRecord(ctx, iter.Next()); err != nil {
			return err
		}
	}
	return nil
}

func (c *Consumer) processRecord(ctx context.Context, record *kgo.Record) error {
	ctx, span := c.tracing.startProcessing(ctx, record, c.cfg.Group)
	defer span.End()

	if err := c.handle(ctx, record); err != nil {
		markSpanError(span, err)
		if ctx.Err() != nil {
			return nil
		}
		c.log.Error("failed to handle record, sending to DLQ",
			slog.String("topic", record.Topic),
			slog.Any("partition", record.Partition),
			slog.Int64("offset", record.Offset),
			slog.Any("error", err),
		)
		if dlqErr := c.sendToDLQ(ctx, record, err); dlqErr != nil {
			if ctx.Err() != nil {
				return nil
			}
			return fmt.Errorf("%s[%d]@%d: %w (handler: %w)", record.Topic, record.Partition, record.Offset, dlqErr, err)
		}
	}

	c.client.MarkCommitRecords(record)
	return nil
}

func (c *Consumer) handle(ctx context.Context, record *kgo.Record) error {
	msg := kafka.Message{
		Key:   record.Key,
		Value: record.Value,
		Topic: record.Topic,
	}
	return c.handleWithRetry(ctx, record, func() error { return c.handler.Handle(ctx, msg) })
}

func (c *Consumer) handleWithRetry(ctx context.Context, record *kgo.Record, handle func() error) error {
	boundedFailures := 0
	for {
		err := handle()
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if err == nil {
			return nil
		}
		c.log.Error("handle record",
			slog.String("topic", record.Topic),
			slog.Any("partition", record.Partition),
			slog.Int64("offset", record.Offset),
			slog.Int("bounded_failures", boundedFailures),
			slog.Any("error", err),
		)
		if kafka.IsPermanentError(err) {
			c.log.Warn("permanent error, skipping retries",
				slog.String("topic", record.Topic),
				slog.Any("partition", record.Partition),
				slog.Int64("offset", record.Offset),
				slog.Any("error", err),
			)
			return err
		}

		delay := c.cfg.WaitTimeout
		if !kafka.IsWaitError(err) {
			boundedFailures++
			if boundedFailures >= c.cfg.MaxAttempts {
				return err
			}
			delay = c.cfg.RetryTimeout
		}
		if err := pauseProcessing(ctx, delay); err != nil {
			return err
		}
	}
}

func pauseProcessing(ctx context.Context, delay time.Duration) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func (c *Consumer) sendToDLQ(ctx context.Context, record *kgo.Record, cause error) error {
	headers := make([]kgo.RecordHeader, 0, len(record.Headers)+6)
	headers = append(headers, record.Headers...)
	msg := &kgo.Record{
		Topic: c.cfg.DLQTopic,
		Key:   record.Key,
		Value: record.Value,
		Headers: append(headers,
			kgo.RecordHeader{Key: "dlq.error", Value: []byte(cause.Error())},
			kgo.RecordHeader{Key: "dlq.consumer.group", Value: []byte(c.cfg.Group)},
			kgo.RecordHeader{Key: "dlq.source.topic", Value: []byte(record.Topic)},
			kgo.RecordHeader{Key: "dlq.source.partition", Value: []byte(fmt.Sprint(record.Partition))},
			kgo.RecordHeader{Key: "dlq.source.offset", Value: []byte(fmt.Sprint(record.Offset))},
			kgo.RecordHeader{Key: "dlq.source.timestamp", Value: []byte(record.Timestamp.UTC().Format(time.RFC3339Nano))},
		),
	}
	publishCtx, span := c.tracing.startPublishing(ctx, msg)
	err := c.client.ProduceSync(publishCtx, msg).FirstErr()
	markSpanError(span, err)
	span.End()
	if err != nil {
		c.log.Error("failed to send record to DLQ",
			slog.String("record_topic", record.Topic),
			slog.String("dlq_topic", c.cfg.DLQTopic),
			slog.Any("error", err),
		)
		return fmt.Errorf("send to DLQ: %w", err)
	}
	c.log.Error("record sent to DLQ",
		slog.String("record_topic", record.Topic),
		slog.String("dlq_topic", c.cfg.DLQTopic),
	)
	if c.deadLetterHandler != nil {
		if err := c.handleWithRetry(ctx, record, func() error {
			return c.deadLetterHandler.OnDeadLetter(ctx, kafka.Message{Key: record.Key, Value: record.Value, Topic: record.Topic})
		}); err != nil {
			return fmt.Errorf("update status after DLQ publication: %w", err)
		}
	}
	return nil
}

func (c *Consumer) shutdown(ctx context.Context) {
	c.log.Warn("shutting down consumer")
	shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), c.cfg.ShutdownTimeout)
	defer cancel()

	if err := c.client.CommitMarkedOffsets(shutdownCtx); err != nil {
		c.log.Error("shutdown commit", slog.Any("error", err))
	}
	c.client.CloseAllowingRebalance()
	c.log.Warn("consumer stopped")
}
