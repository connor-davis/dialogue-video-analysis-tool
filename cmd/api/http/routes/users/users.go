package users

import (
	"github.com/connor-davis/dialogue-video-analysis-tool/cmd/api/http/middleware"
	"github.com/connor-davis/dialogue-video-analysis-tool/cmd/api/http/routes"
	"github.com/connor-davis/dialogue-video-analysis-tool/internal/api/assignApi"
	"github.com/connor-davis/dialogue-video-analysis-tool/internal/api/baseApi"
	"github.com/connor-davis/dialogue-video-analysis-tool/internal/models"
	"github.com/connor-davis/dialogue-video-analysis-tool/internal/routing"
	"github.com/connor-davis/dialogue-video-analysis-tool/internal/storage"
)

type UsersRouter struct {
	storage    storage.Storage
	middleware middleware.Middleware
}

func New(storage storage.Storage, middleware middleware.Middleware) routes.Router {
	return &UsersRouter{
		storage:    storage,
		middleware: middleware,
	}
}

func (r *UsersRouter) LoadRoutes() []routing.Route {
	userOrganizationAssignmentApi := assignApi.New[models.User, models.Organization](
		r.storage,
		"/users",
		"User",
		"Organization",
	)
	userRoleAssignmentApi := assignApi.New[models.User, models.Role](
		r.storage,
		"/users",
		"User",
		"Role",
	)
	usersApi := baseApi.New[models.User](
		r.storage,
		"/users",
		"User",
	)

	routes := []routing.Route{}

	routes = append(routes, []routing.Route{
		userOrganizationAssignmentApi.AssignRoute(
			r.middleware.Authenticated(),
			r.middleware.Authorized("users.organizations.assign"),
		),
		userOrganizationAssignmentApi.UnassignRoute(
			r.middleware.Authenticated(),
			r.middleware.Authorized("users.organizations.unassign"),
		),
		userOrganizationAssignmentApi.ListRoute(
			r.middleware.Authenticated(),
			r.middleware.Authorized("users.organizations.list"),
		),

		userRoleAssignmentApi.AssignRoute(
			r.middleware.Authenticated(),
			r.middleware.Authorized("users.roles.assign"),
		),
		userRoleAssignmentApi.UnassignRoute(
			r.middleware.Authenticated(),
			r.middleware.Authorized("users.roles.unassign"),
		),
		userRoleAssignmentApi.ListRoute(
			r.middleware.Authenticated(),
			r.middleware.Authorized("users.roles.list"),
		),

		usersApi.GetAllRoute(
			r.middleware.Authenticated(),
			r.middleware.Authorized("users.list"),
		),
		usersApi.GetOneRoute(
			r.middleware.Authenticated(),
			r.middleware.Authorized("users.view"),
		),
		usersApi.CreateRoute(
			"#/components/requestBodies/CreateUserPayload",
			r.middleware.Authenticated(),
			r.middleware.Authorized("users.create"),
		),
		usersApi.UpdateRoute(
			"#/components/requestBodies/UpdateUserPayload",
			r.middleware.Authenticated(),
			r.middleware.Authorized("users.update"),
		),
		usersApi.DeleteRoute(
			r.middleware.Authenticated(),
			r.middleware.Authorized("users.delete"),
		),
	}...)

	return routes
}
