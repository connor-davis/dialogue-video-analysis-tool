package http

import (
	"fmt"
	"regexp"

	"github.com/connor-davis/dialogue-video-analysis-tool/cmd/api/http/middleware"
	"github.com/connor-davis/dialogue-video-analysis-tool/cmd/api/http/routes/authentication"
	"github.com/connor-davis/dialogue-video-analysis-tool/cmd/api/http/routes/roles"
	"github.com/connor-davis/dialogue-video-analysis-tool/cmd/api/http/routes/users"
	"github.com/connor-davis/dialogue-video-analysis-tool/common"
	"github.com/connor-davis/dialogue-video-analysis-tool/internal/routing"
	"github.com/connor-davis/dialogue-video-analysis-tool/internal/routing/bodies"
	"github.com/connor-davis/dialogue-video-analysis-tool/internal/routing/parameters"
	"github.com/connor-davis/dialogue-video-analysis-tool/internal/routing/schemas"
	"github.com/connor-davis/dialogue-video-analysis-tool/internal/storage"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v3"
	"github.com/openai/openai-go/v3"
)

type HttpRouter interface {
	InitializeRoutes(router fiber.Router)
	InitializeOpenAPI() *openapi3.T
}

type httpRouter struct {
	storage    storage.Storage
	middleware middleware.Middleware
	openai     openai.Client
	routes     []routing.Route
}

func New(storage storage.Storage, middleware middleware.Middleware, openai openai.Client) HttpRouter {
	authenticationRouter := authentication.New(storage, middleware)
	authenticationRoutes := authenticationRouter.LoadRoutes()

	usersRouter := users.New(storage, middleware)
	usersRoutes := usersRouter.LoadRoutes()

	rolesRouter := roles.New(storage, middleware)
	rolesRoutes := rolesRouter.LoadRoutes()

	routes := []routing.Route{}

	routes = append(routes, authenticationRoutes...)
	routes = append(routes, usersRoutes...)
	routes = append(routes, rolesRoutes...)

	return &httpRouter{
		storage:    storage,
		middleware: middleware,
		openai:     openai,
		routes:     routes,
	}
}

func (h *httpRouter) InitializeRoutes(router fiber.Router) {
	for _, route := range h.routes {
		path := regexp.MustCompile(`\{([^}]+)\}`).ReplaceAllString(route.Path, ":$1")

		routes := route.Middlewares
		routes = append(routes, route.Handler)

		switch route.Method {
		case routing.GET:
			router.Get(path, routes[0], routes[1:]...)
		case routing.POST:
			router.Post(path, routes[0], routes[1:]...)
		case routing.PUT:
			router.Put(path, routes[0], routes[1:]...)
		case routing.PATCH:
			router.Patch(path, routes[0], routes[1:]...)
		case routing.DELETE:
			router.Delete(path, routes[0], routes[1:]...)
		}
	}
}

func (h *httpRouter) InitializeOpenAPI() *openapi3.T {
	paths := openapi3.NewPaths()

	parameters := openapi3.ParametersMap{
		"Id":                parameters.IdParameter,
		"DisablePagination": parameters.DisablePaginationParameter,
		"Page":              parameters.PageParameter,
		"PageSize":          parameters.PageSizeParameter,
		"SearchTerm":        parameters.SearchTermParameter,
		"SearchColumn":      parameters.SearchColumnParameter,
		"Preload":           parameters.PreloadParameter,
		"Code":              parameters.CodeParameter,
		"State":             parameters.StateParameter,
		"To":                parameters.ToParameter,
	}

	bodies := openapi3.RequestBodies{
		"CreateUserPayload": bodies.CreateUserSchema,
		"UpdateUserPayload": bodies.UpdateUserSchema,
		"CreateRolePayload": bodies.CreateRoleSchema,
		"UpdateRolePayload": bodies.UpdateRoleSchema,
	}

	schemas := openapi3.Schemas{
		"SuccessResponse": schemas.SuccessSchema,
		"ErrorResponse":   schemas.ErrorSchema,
		"Pagination":      schemas.PaginationSchema,
		"User":            schemas.UserSchema,
		"Users":           schemas.UsersSchema,
		"Role":            schemas.RoleSchema,
		"Roles":           schemas.RolesSchema,
	}

	for _, route := range h.routes {
		pathItem := &openapi3.PathItem{}

		switch route.Method {
		case routing.GET:
			pathItem.Get = &openapi3.Operation{
				Summary:     route.Summary,
				Description: route.Description,
				Tags:        route.Tags,
				Parameters:  route.Parameters,
				Responses:   route.Responses,
			}
		case routing.POST:
			pathItem.Post = &openapi3.Operation{
				Summary:     route.Summary,
				Description: route.Description,
				Tags:        route.Tags,
				Parameters:  route.Parameters,
				RequestBody: route.RequestBody,
				Responses:   route.Responses,
			}
		case routing.PUT:
			pathItem.Put = &openapi3.Operation{
				Summary:     route.Summary,
				Description: route.Description,
				Tags:        route.Tags,
				Parameters:  route.Parameters,
				RequestBody: route.RequestBody,
				Responses:   route.Responses,
			}
		case routing.PATCH:
			pathItem.Patch = &openapi3.Operation{
				Summary:     route.Summary,
				Description: route.Description,
				Tags:        route.Tags,
				Parameters:  route.Parameters,
				RequestBody: route.RequestBody,
				Responses:   route.Responses,
			}
		case routing.DELETE:
			pathItem.Delete = &openapi3.Operation{
				Summary:     route.Summary,
				Description: route.Description,
				Tags:        route.Tags,
				Parameters:  route.Parameters,
				Responses:   route.Responses,
			}
		}

		path := fmt.Sprintf("/api/v1%s", route.Path)

		existingPathItem := paths.Find(path)

		if existingPathItem != nil {
			switch route.Method {
			case routing.GET:
				existingPathItem.Get = pathItem.Get
			case routing.POST:
				existingPathItem.Post = pathItem.Post
			case routing.PUT:
				existingPathItem.Put = pathItem.Put
			case routing.PATCH:
				existingPathItem.Patch = pathItem.Patch
			case routing.DELETE:
				existingPathItem.Delete = pathItem.Delete
			}
		} else {
			paths.Set(path, pathItem)
		}
	}

	return &openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:   common.EnvString("API_NAME", "One REST API"),
			Version: common.EnvString("API_VERSION", "1.0.0"),
		},
		Servers: openapi3.Servers{
			{
				URL:         fmt.Sprintf("http://localhost:%s", common.EnvString("API_PORT", "6173")),
				Description: "Development",
			},
			{
				URL:         common.EnvString("API_BASE_URL", "https://example.com"),
				Description: "Production",
			},
		},
		Tags:  openapi3.Tags{},
		Paths: paths,
		Components: &openapi3.Components{
			Schemas:       schemas,
			RequestBodies: bodies,
			Parameters:    parameters,
		},
	}
}
