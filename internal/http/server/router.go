package server

import (
	"fmt"
	"net/http"

	"image-uploader/internal/http/middleware"
)

type APIVersion string

const (
	APIVersion1 = APIVersion("v1")
)

type APIVersionRouter struct {
	*http.ServeMux
	apiVersion     APIVersion
	middleware     []middleware.Middleware
	noPrefixRoutes []Route
}

func NewAPIVersionRouter(apiVersion APIVersion, middleware ...middleware.Middleware) *APIVersionRouter {
	return &APIVersionRouter{
		ServeMux:   http.NewServeMux(),
		apiVersion: apiVersion,
		middleware: middleware,
	}
}

func (r *APIVersionRouter) prefix() string {
	return "/api/" + string(r.apiVersion)
}

func (r *APIVersionRouter) RegisterRoutes(routes ...Route) {
	prefix := r.prefix()
	for _, route := range routes {
		if route.NoPrefix {
			r.noPrefixRoutes = append(r.noPrefixRoutes, route)
			continue
		}

		pattern := fmt.Sprintf("%s %s%s", route.Method, prefix, route.Path)
		r.Handle(pattern, http.StripPrefix(prefix, route.WithMiddleware()))
	}
}

func (r *APIVersionRouter) NoPrefixRoutes() []Route {
	routes := make([]Route, 0, len(r.noPrefixRoutes))

	for _, route := range r.noPrefixRoutes {
		combinedMiddleware := make([]middleware.Middleware, 0, len(route.Middleware)+len(r.middleware))
		combinedMiddleware = append(combinedMiddleware, route.Middleware...)
		combinedMiddleware = append(combinedMiddleware, r.middleware...)

		routeWithRouterMiddleware := Route{
			Method:     route.Method,
			Path:       route.Path,
			Handler:    route.Handler,
			Middleware: combinedMiddleware,
			NoPrefix:   route.NoPrefix,
		}
		routes = append(routes, routeWithRouterMiddleware)
	}

	return routes
}

func (r *APIVersionRouter) WithMiddleware() http.Handler {
	return middleware.ChainMiddleware(r, r.middleware...)
}
