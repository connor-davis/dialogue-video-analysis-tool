package mfa

import (
	"github.com/connor-davis/dialogue-video-analysis-tool/cmd/api/http/middleware"
	"github.com/connor-davis/dialogue-video-analysis-tool/cmd/api/http/routes"
	"github.com/connor-davis/dialogue-video-analysis-tool/cmd/api/http/routes/authentication/mfa/totp"
	"github.com/connor-davis/dialogue-video-analysis-tool/internal/routing"
	"github.com/connor-davis/dialogue-video-analysis-tool/internal/storage"
)

type MfaRouter struct {
	storage    storage.Storage
	middleware middleware.Middleware
	totp       routes.Router
}

func NewMfaRouter(storage storage.Storage, middleware middleware.Middleware) routes.Router {
	totp := totp.NewTotpRouter(storage, middleware)

	return &MfaRouter{
		storage:    storage,
		middleware: middleware,
		totp:       totp,
	}
}

func (r *MfaRouter) LoadRoutes() []routing.Route {
	routes := []routing.Route{}

	routes = append(routes, r.totp.LoadRoutes()...)

	return routes
}
