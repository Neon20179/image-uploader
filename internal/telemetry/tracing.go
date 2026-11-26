package telemetry

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.40.0"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

type Tracing struct {
	Provider   trace.TracerProvider
	Propagator propagation.TraceContext
}

func StartTracing(ctx context.Context, log *slog.Logger, serviceName string) (Tracing, func()) {
	tracing := Tracing{Provider: noop.NewTracerProvider()}
	stop := func() {}
	if os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") == "" && os.Getenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT") == "" {
		return tracing, stop
	}

	exporter, err := otlptracegrpc.New(ctx)
	if err != nil {
		log.Error("initialize trace exporter", slog.Any("error", err))
		return tracing, stop
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(semconv.ServiceName(serviceName)),
		resource.WithTelemetrySDK(),
		resource.WithFromEnv(),
	)
	if err != nil {
		log.Warn("read OpenTelemetry resource attributes")
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	tracing.Provider = provider
	return tracing, func() {
		shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
		defer cancel()
		if err := provider.Shutdown(shutdownCtx); err != nil {
			log.Error("shut down trace provider", slog.Any("error", err))
		}
	}
}

func (t Tracing) HTTPMiddleware(filter func(*http.Request) bool) func(http.Handler) http.Handler {
	instrument := otelhttp.NewMiddleware(
		"",
		otelhttp.WithTracerProvider(t.Provider),
		otelhttp.WithPropagators(t.Propagator),
		otelhttp.WithFilter(filter),
		otelhttp.WithSpanNameFormatter(func(_ string, r *http.Request) string {
			if route := routePath(r.Pattern); route != "" {
				return r.Method + " " + route
			}
			return r.Method
		}),
	)
	return func(next http.Handler) http.Handler {
		return instrument(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
			if route := routePath(r.Pattern); route != "" {
				trace.SpanFromContext(r.Context()).SetAttributes(attribute.String("http.route", route))
			}
		}))
	}
}

func routePath(pattern string) string {
	_, path, _ := strings.Cut(pattern, " ")
	return path
}
