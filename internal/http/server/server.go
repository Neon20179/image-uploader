package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"image-uploader/internal/http/middleware"
)

type HTTPServer struct {
	mux        *http.ServeMux
	cfg        Config
	log        *slog.Logger
	middleware []middleware.Middleware
}

func NewHTTPServer(cfg Config, log *slog.Logger, middleware []middleware.Middleware) *HTTPServer {
	return &HTTPServer{
		mux:        http.NewServeMux(),
		cfg:        cfg,
		log:        log,
		middleware: middleware,
	}
}

func (s *HTTPServer) RegisterAPIRouters(routers ...*APIVersionRouter) {
	for _, router := range routers {
		s.mux.Handle(router.prefix()+"/", router.WithMiddleware())

		for _, route := range router.NoPrefixRoutes() {
			pattern := fmt.Sprintf("%s %s", route.Method, route.Path)
			s.mux.Handle(pattern, route.WithMiddleware())
		}
	}
}

func (s *HTTPServer) Run(ctx context.Context) error {
	mux := middleware.ChainMiddleware(s.mux, s.middleware...)

	server := http.Server{
		Addr:              s.cfg.Addr,
		Handler:           mux,
		ReadHeaderTimeout: s.cfg.ReadHeaderTimeout,
		ReadTimeout:       s.cfg.ReadTimeout,
		WriteTimeout:      s.cfg.WriteTimeout,
		IdleTimeout:       s.cfg.IdleTimeout,
	}

	errCh := make(chan error, 1)

	go func() {
		defer close(errCh)

		s.log.Warn("HTTP server started", slog.String("addr", s.cfg.Addr))

		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("listen and serve HTTP server: %w", err)
		}
	case <-ctx.Done():
		s.log.Warn("shutdown HTTP server...")

		shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), s.cfg.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			if err := server.Close(); err != nil {
				return fmt.Errorf("close HTTP server: %w", err)
			}

			return fmt.Errorf("shutdown HTTP server: %w", err)
		}

		s.log.Warn("HTTP server stopped")
	}

	return nil
}
