FROM golang:1.26.4-alpine AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY cmd ./cmd
COPY internal ./internal
COPY docs ./docs
ARG SERVICE=api
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    case "$SERVICE" in \
      api|ingest|outbox|normalizer|embedding|ocr|finalizer) ;; \
      *) echo "Unsupported Go service: $SERVICE" >&2; exit 1 ;; \
    esac && \
    CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/service ./cmd/${SERVICE}

FROM alpine:3.22 AS runner
RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S -g 10001 app \
    && adduser -S -D -H -u 10001 -G app app \
    && mkdir -p /app/out/logs \
    && chown -R app:app /app

WORKDIR /app
COPY --from=builder /out/service /app/service
USER app
ENTRYPOINT ["/app/service"]
