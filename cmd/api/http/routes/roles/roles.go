package roles

import (
	"github.com/connor-davis/dialogue-video-analysis-tool/cmd/api/http/middleware"
	"github.com/connor-davis/dialogue-video-analysis-tool/cmd/api/http/routes"
	"github.com/connor-davis/dialogue-video-analysis-tool/internal/api/baseApi"
	"github.com/connor-davis/dialogue-video-analysis-tool/internal/models"
	"github.com/connor-davis/dialogue-video-analysis-tool/internal/routing"
	"github.com/connor-davis/dialogue-video-analysis-tool/internal/storage"
)

type RolesRouter struct {
	storage    storage.Storage
	middleware middleware.Middleware
}

func New(storage storage.Storage, middleware middleware.Middleware) routes.Router {
	return &RolesRouter{
		storage:    storage,
		middleware: middleware,
	}
}

func (r *RolesRouter) LoadRoutes() []routing.Route {
	rolesApi := baseApi.New[models.Role](r.storage, "/roles", "Role")

	return []routing.Route{
		rolesApi.GetAllRoute(
			r.middleware.Authenticated(),
			r.middleware.Authorized("roles.list"),
		),
		rolesApi.GetOneRoute(
			r.middleware.Authenticated(),
			r.middleware.Authorized("roles.view"),
		),
		rolesApi.CreateRoute(
			"#/components/requestBodies/CreateRolePayload",
			r.middleware.Authenticated(),
			r.middleware.Authorized("roles.create"),
		),
		rolesApi.UpdateRoute(
			"#/components/requestBodies/UpdateRolePayload",
			r.middleware.Authenticated(),
			r.middleware.Authorized("roles.update"),
		),
		rolesApi.DeleteRoute(
			r.middleware.Authenticated(),
			r.middleware.Authorized("roles.delete"),
		),
	}
}
