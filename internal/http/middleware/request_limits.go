package middleware

import (
	"net/http"
	"strings"

	"golang.org/x/time/rate"

	"image-uploader/internal/http/response"
	"image-uploader/internal/logger"
)

func RequestLimits(cfg RequestLimitsConfig, pathPrefix string) Middleware {
	limiter := rate.NewLimiter(rate.Limit(cfg.RequestsPerSecond), cfg.Burst)
	inFlight := make(chan struct{}, cfg.MaxInFlight)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !strings.HasPrefix(r.URL.Path, pathPrefix) {
				next.ServeHTTP(w, r)
				return
			}

			select {
			case inFlight <- struct{}{}:
				defer func() { <-inFlight }()
			default:
				writeLimitResponse(w, r, http.StatusServiceUnavailable)
				return
			}

			if !limiter.Allow() {
				writeLimitResponse(w, r, http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func writeLimitResponse(w http.ResponseWriter, r *http.Request, status int) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	handler := response.NewHTTPResponseHandler(log, w)
	handler.JSONResponse(map[string]string{"message": http.StatusText(status)}, status)
}
