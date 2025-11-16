package totp

import (
	"github.com/connor-davis/dialogue-video-analysis-tool/cmd/api/http/middleware"
	"github.com/connor-davis/dialogue-video-analysis-tool/cmd/api/http/routes"
	"github.com/connor-davis/dialogue-video-analysis-tool/internal/routing"
	"github.com/connor-davis/dialogue-video-analysis-tool/internal/storage"
)

type TotpRouter struct {
	storage    storage.Storage
	middleware middleware.Middleware
}

func NewTotpRouter(storage storage.Storage, middleware middleware.Middleware) routes.Router {
	return &TotpRouter{
		storage:    storage,
		middleware: middleware,
	}
}

func (r *TotpRouter) LoadRoutes() []routing.Route {
	routes := []routing.Route{
		r.EnableRoute(),
		r.VerifyRoute(),
		r.ResetRoute(),
	}

	return routes
}
