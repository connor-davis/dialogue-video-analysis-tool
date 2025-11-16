package authentication

import (
	"github.com/connor-davis/dialogue-video-analysis-tool/cmd/api/http/middleware"
	"github.com/connor-davis/dialogue-video-analysis-tool/cmd/api/http/routes"
	"github.com/connor-davis/dialogue-video-analysis-tool/cmd/api/http/routes/authentication/mfa"
	"github.com/connor-davis/dialogue-video-analysis-tool/internal/routing"
	"github.com/connor-davis/dialogue-video-analysis-tool/internal/storage"
)

type AuthenticationRouter struct {
	storage    storage.Storage
	middleware middleware.Middleware
	mfa        routes.Router
}

func New(storage storage.Storage, middleware middleware.Middleware) routes.Router {
	mfa := mfa.NewMfaRouter(storage, middleware)

	return &AuthenticationRouter{
		storage:    storage,
		middleware: middleware,
		mfa:        mfa,
	}
}

func (r *AuthenticationRouter) LoadRoutes() []routing.Route {
	routes := []routing.Route{
		r.CheckRoute(),
		r.MeRoute(),
		r.LogoutRoute(),
	}

	routes = append(routes, r.mfa.LoadRoutes()...)

	return routes
}
