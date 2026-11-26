package franzgo

import (
	"context"
	"fmt"
	"slices"

	"github.com/twmb/franz-go/pkg/kgo"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type Tracing struct {
	tracer     trace.Tracer
	propagator propagation.TextMapPropagator
}

func NewTracing(provider trace.TracerProvider, propagator propagation.TextMapPropagator) Tracing {
	return Tracing{tracer: provider.Tracer("image-uploader/kafka"), propagator: propagator}
}

func (t Tracing) startProcessing(ctx context.Context, record *kgo.Record, group string) (context.Context, trace.Span) {
	ctx = t.propagator.Extract(ctx, recordCarrier{record})
	return t.tracer.Start(ctx, "process "+record.Topic,
		trace.WithSpanKind(trace.SpanKindConsumer),
		trace.WithAttributes(
			attribute.String("messaging.system", "kafka"),
			attribute.String("messaging.operation.type", "process"),
			attribute.String("messaging.destination.name", record.Topic),
			attribute.String("messaging.consumer.group.name", group),
		),
	)
}

func (t Tracing) startPublishing(ctx context.Context, record *kgo.Record) (context.Context, trace.Span) {
	ctx, span := t.tracer.Start(ctx, "send "+record.Topic,
		trace.WithSpanKind(trace.SpanKindProducer),
		trace.WithAttributes(
			attribute.String("messaging.system", "kafka"),
			attribute.String("messaging.operation.type", "send"),
			attribute.String("messaging.destination.name", record.Topic),
		),
	)
	t.propagator.Inject(ctx, recordCarrier{record})
	return ctx, span
}

func markSpanError(span trace.Span, err error) {
	if err != nil {
		span.SetStatus(codes.Error, "")
		span.SetAttributes(attribute.String("error.type", fmt.Sprintf("%T", err)))
	}
}

type recordCarrier struct {
	record *kgo.Record
}

func (c recordCarrier) Get(key string) string {
	for i := len(c.record.Headers) - 1; i >= 0; i-- {
		if header := c.record.Headers[i]; header.Key == key {
			return string(header.Value)
		}
	}
	return ""
}

func (c recordCarrier) Set(key, value string) {
	c.record.Headers = slices.DeleteFunc(c.record.Headers, func(header kgo.RecordHeader) bool {
		return header.Key == key
	})
	c.record.Headers = append(c.record.Headers, kgo.RecordHeader{Key: key, Value: []byte(value)})
}

func (c recordCarrier) Keys() []string {
	keys := make([]string, 0, len(c.record.Headers))
	for _, header := range c.record.Headers {
		keys = append(keys, header.Key)
	}
	return keys
}
