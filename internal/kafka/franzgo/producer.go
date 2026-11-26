package franzgo

import (
	"context"
	"fmt"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
	"go.opentelemetry.io/otel/propagation"

	"image-uploader/internal/kafka"
)

type ProducerConfig struct {
	Broker          string
	DeliveryTimeout time.Duration
	Linger          time.Duration
	ShutdownTimeout time.Duration
}

type Producer struct {
	client  *kgo.Client
	cfg     ProducerConfig
	tracing Tracing
}

func NewProducer(cfg ProducerConfig, tracing Tracing) (*Producer, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.Broker),
		kgo.AllowAutoTopicCreation(),
		kgo.RecordDeliveryTimeout(cfg.DeliveryTimeout),
		kgo.ProducerLinger(cfg.Linger),
	)
	if err != nil {
		return nil, fmt.Errorf("initializing broker: %w", err)
	}

	return &Producer{
		client:  client,
		cfg:     cfg,
		tracing: tracing,
	}, nil
}

func (p *Producer) Send(ctx context.Context, msg kafka.Message) error {
	record := msgToRecord(msg)

	ctx = p.tracing.propagator.Extract(ctx, propagation.MapCarrier(msg.Headers))
	ctx, span := p.tracing.startPublishing(ctx, record)
	defer span.End()

	err := p.client.ProduceSync(ctx, record).FirstErr()
	markSpanError(span, err)
	if err != nil {
		return fmt.Errorf("produce to %s: %w", msg.Topic, err)
	}
	return nil
}

func (p *Producer) Close(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), p.cfg.ShutdownTimeout)
	defer cancel()

	err := p.client.Flush(ctx)
	p.client.Close()
	if err != nil {
		return fmt.Errorf("flush producer: %w", err)
	}
	return nil
}

func msgToRecord(msg kafka.Message) *kgo.Record {
	record := &kgo.Record{
		Key:   msg.Key,
		Value: msg.Value,
		Topic: msg.Topic,
	}
	for key, value := range msg.Headers {
		record.Headers = append(record.Headers, kgo.RecordHeader{Key: key, Value: []byte(value)})
	}
	return record
}
