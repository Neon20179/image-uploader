package server

import (
	"net/http"

	"image-uploader/internal/http/middleware"
)

type Route struct {
	Method     string
	Path       string
	Handler    http.HandlerFunc
	Middleware []middleware.Middleware
	NoPrefix   bool
}

func (r *Route) WithMiddleware() http.Handler {
	return middleware.ChainMiddleware(r.Handler, r.Middleware...)
}
