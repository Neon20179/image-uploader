package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	_ "image-uploader/docs"
	imagehandler "image-uploader/internal/handler/image"
	uploadhandler "image-uploader/internal/handler/upload"
	"image-uploader/internal/http/middleware"
	"image-uploader/internal/http/server"
	"image-uploader/internal/logger"
	"image-uploader/internal/postgres"
	"image-uploader/internal/postgres/pgx"
	imagerepository "image-uploader/internal/repository/image"
	normalizationrepository "image-uploader/internal/repository/normalization"
	outboxrepository "image-uploader/internal/repository/outbox"
	uploadrepository "image-uploader/internal/repository/upload"
	"image-uploader/internal/s3/minio"
	imageservice "image-uploader/internal/service/image"
	uploadservice "image-uploader/internal/service/upload"
	"image-uploader/internal/telemetry"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// @title Image Uploader API
// @version 1.0
// @BasePath /api/v1
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

	tracing, stopTracing := telemetry.StartTracing(ctx, log.Logger, "image-uploader-api")
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
	imageCfg := imagehandler.MustNewConfig()
	urlBuilder := minio.NewURLBuilder(imageCfg.PublicBaseURL)

	outboxRepo := tracing.TraceOutboxWrites(outboxrepository.NewOutboxRepository(pool))

	log.Debug("initializing image api")
	imageRepo := imagerepository.NewImageRepository(pool)
	uploadCfg := uploadservice.MustNewConfig()
	imageSvc := imageservice.NewImageService(imageRepo, s3Client, uploadCfg.Bucket)
	imageHTTPHandler := imagehandler.NewImageHandler(imageSvc, urlBuilder)

	log.Debug("initializing upload api")
	uploadRepo := uploadrepository.NewUploadRepository(pool)
	normalizationRepo := normalizationrepository.NewNormalizationRepository(pool)
	uploadSvc := uploadservice.NewUploadService(uploadRepo, outboxRepo, normalizationRepo, txManager, s3Client, uploadCfg)
	uploadHTTPHandler := uploadhandler.NewUploadHandler(uploadSvc, urlBuilder)

	log.Debug("initializing HTTP server")
	corsCfg := middleware.MustNewConfig()
	middlewares := []middleware.Middleware{
		middleware.CORS(corsCfg),
		middleware.RequestID(),
		middleware.Logger(log.Logger),
		middleware.Trace(),
		tracing.HTTPMiddleware(func(r *http.Request) bool {
			return strings.HasPrefix(r.URL.Path, "/api/")
		}),
		middleware.Panic(),
		middleware.RequestLimits(middleware.MustNewRequestLimitsConfig(), "/api/v1/"),
	}
	serverCfg := server.MustNewConfig()
	httpServer := server.NewHTTPServer(
		serverCfg,
		log.Logger,
		middlewares,
	)
	apiVersionRouter := server.NewAPIVersionRouter(server.APIVersion1)
	apiVersionRouter.RegisterRoutes(
		server.Route{
			Method:   "GET",
			Path:     "/healthz",
			NoPrefix: true,
			Handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}},
		server.Route{
			Method:   "GET",
			Path:     "/readyz",
			NoPrefix: true,
			Handler: func(w http.ResponseWriter, r *http.Request) {
				probeCtx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
				defer cancel()
				if err := pool.Ping(probeCtx); err != nil {
					http.Error(w, "database unavailable", http.StatusServiceUnavailable)
					return
				}
				w.WriteHeader(http.StatusOK)
			}},
	)
	apiVersionRouter.RegisterRoutes(server.Route{
		Method:   "GET",
		Path:     "/swagger/",
		Handler:  httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json")),
		NoPrefix: true,
	})

	apiVersionRouter.RegisterRoutes(imageHTTPHandler.Routes()...)
	apiVersionRouter.RegisterRoutes(uploadHTTPHandler.Routes()...)

	httpServer.RegisterAPIRouters(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		log.Error("HTTP server run error", slog.Any("error", err))
	}
}
