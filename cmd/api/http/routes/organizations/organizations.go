package organizations

import (
	"github.com/connor-davis/dialogue-video-analysis-tool/cmd/api/http/middleware"
	"github.com/connor-davis/dialogue-video-analysis-tool/cmd/api/http/routes"
	"github.com/connor-davis/dialogue-video-analysis-tool/internal/api/assignApi"
	"github.com/connor-davis/dialogue-video-analysis-tool/internal/api/baseApi"
	"github.com/connor-davis/dialogue-video-analysis-tool/internal/models"
	"github.com/connor-davis/dialogue-video-analysis-tool/internal/routing"
	"github.com/connor-davis/dialogue-video-analysis-tool/internal/storage"
)

type OrganizationsRouter struct {
	storage    storage.Storage
	middleware middleware.Middleware
}

func New(storage storage.Storage, middleware middleware.Middleware) routes.Router {
	return &OrganizationsRouter{
		storage:    storage,
		middleware: middleware,
	}
}

func (r *OrganizationsRouter) LoadRoutes() []routing.Route {
	organizationUserAssignmentApi := assignApi.New[models.Organization, models.User](
		r.storage,
		"/organizations",
		"Organization",
		"User",
	)
	organizationRoleAssignmentApi := assignApi.New[models.Organization, models.Role](
		r.storage,
		"/organizations",
		"Organization",
		"Role",
	)
	organizationsApi := baseApi.New[models.Organization](r.storage, "/organizations", "Organization")

	return []routing.Route{
		organizationUserAssignmentApi.AssignRoute(
			r.middleware.Authenticated(),
			r.middleware.Authorized("organizations.users.assign"),
		),
		organizationUserAssignmentApi.UnassignRoute(
			r.middleware.Authenticated(),
			r.middleware.Authorized("organizations.users.unassign"),
		),
		organizationUserAssignmentApi.ListRoute(
			r.middleware.Authenticated(),
			r.middleware.Authorized("organizations.users.list"),
		),

		organizationRoleAssignmentApi.AssignRoute(
			r.middleware.Authenticated(),
			r.middleware.Authorized("organizations.roles.assign"),
		),
		organizationRoleAssignmentApi.UnassignRoute(
			r.middleware.Authenticated(),
			r.middleware.Authorized("organizations.roles.unassign"),
		),
		organizationRoleAssignmentApi.ListRoute(
			r.middleware.Authenticated(),
			r.middleware.Authorized("organizations.roles.list"),
		),

		organizationsApi.GetAllRoute(
			r.middleware.Authenticated(),
			r.middleware.Authorized("organizations.list"),
		),
		organizationsApi.GetOneRoute(
			r.middleware.Authenticated(),
			r.middleware.Authorized("organizations.view"),
		),
		organizationsApi.CreateRoute(
			"#/components/requestBodies/CreateOrganizationPayload",
			r.middleware.Authenticated(),
			r.middleware.Authorized("organizations.create"),
		),
		organizationsApi.UpdateRoute(
			"#/components/requestBodies/UpdateOrganizationPayload",
			r.middleware.Authenticated(),
			r.middleware.Authorized("organizations.update"),
		),
		organizationsApi.DeleteRoute(
			r.middleware.Authenticated(),
			r.middleware.Authorized("organizations.delete"),
		),
	}
}
