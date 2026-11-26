package telemetry

import (
	"context"
	"maps"

	"image-uploader/internal/domain"

	"go.opentelemetry.io/otel/propagation"
)

type outboxWriter interface {
	Add(context.Context, domain.OutboxEvent) (domain.OutboxEvent, error)
}

type tracedOutboxWriter struct {
	writer     outboxWriter
	propagator propagation.TextMapPropagator
}

func (t Tracing) TraceOutboxWrites(writer outboxWriter) tracedOutboxWriter {
	return tracedOutboxWriter{writer: writer, propagator: t.Propagator}
}

func (w tracedOutboxWriter) Add(ctx context.Context, event domain.OutboxEvent) (domain.OutboxEvent, error) {
	headers := make(map[string]string, len(event.Headers))
	maps.Copy(headers, event.Headers)
	w.propagator.Inject(ctx, propagation.MapCarrier(headers))
	event.Headers = headers
	return w.writer.Add(ctx, event)
}
