package baseApi

import (
	"fmt"
	"strings"

	"github.com/connor-davis/dialogue-video-analysis-tool/internal/routing"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-openapi/inflect"
	"github.com/gofiber/fiber/v3"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type GetOneParams struct {
	Id string `json:"id" msg:"id"`
}

type GetOneQueryParams struct {
	Preloads pq.StringArray `query:"preload"`
}

func (b *baseApi[Entity]) GetOneRoute(middleware ...fiber.Handler) routing.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription(fmt.Sprintf("%s retrieved successfully.", b.name)).
			WithContent(openapi3.Content{
				"text/plain": openapi3.NewMediaType().
					WithSchemaRef(&openapi3.SchemaRef{
						Ref: "#/components/schemas/SuccessResponse",
					}),
			}),
	})

	responses.Set("400", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchemaRef(&openapi3.SchemaRef{
				Ref: "#/components/schemas/ErrorResponse",
			}).
			WithDescription("Bad Request").
			WithContent(openapi3.Content{
				"application/json": openapi3.NewMediaType().
					WithSchemaRef(&openapi3.SchemaRef{
						Ref: "#/components/schemas/ErrorResponse",
					}),
			}),
	})

	responses.Set("401", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchemaRef(&openapi3.SchemaRef{
				Ref: "#/components/schemas/ErrorResponse",
			}).
			WithDescription("Unauthorized").
			WithContent(openapi3.Content{
				"application/json": openapi3.NewMediaType().
					WithSchemaRef(&openapi3.SchemaRef{
						Ref: "#/components/schemas/ErrorResponse",
					}),
			}),
	})

	responses.Set("403", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchemaRef(&openapi3.SchemaRef{
				Ref: "#/components/schemas/ErrorResponse",
			}).
			WithDescription("Forbidden").
			WithContent(openapi3.Content{
				"application/json": openapi3.NewMediaType().
					WithSchemaRef(&openapi3.SchemaRef{
						Ref: "#/components/schemas/ErrorResponse",
					}),
			}),
	})

	responses.Set("404", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchemaRef(&openapi3.SchemaRef{
				Ref: "#/components/schemas/ErrorResponse",
			}).
			WithDescription("Not Found").
			WithContent(openapi3.Content{
				"application/json": openapi3.NewMediaType().
					WithSchemaRef(&openapi3.SchemaRef{
						Ref: "#/components/schemas/ErrorResponse",
					}),
			}),
	})

	responses.Set("500", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchemaRef(&openapi3.SchemaRef{
				Ref: "#/components/schemas/ErrorResponse",
			}).
			WithDescription("Internal Server Error").
			WithContent(openapi3.Content{
				"application/json": openapi3.NewMediaType().
					WithSchemaRef(&openapi3.SchemaRef{
						Ref: "#/components/schemas/ErrorResponse",
					}),
			}),
	})

	return routing.Route{
		OpenAPIMetadata: routing.OpenAPIMetadata{
			Summary: fmt.Sprintf(
				"Get %s",
				b.name,
			),
			Description: fmt.Sprintf(
				"This endpoint retrieves an existing %s by their id.",
				strings.ToLower(b.name),
			),
			Tags: []string{fmt.Sprintf(
				"%s",
				inflect.Pluralize(b.name),
			)},
			Parameters: []*openapi3.ParameterRef{
				{
					Ref: "#/components/parameters/Id",
				},
				{
					Ref: "#/components/parameters/Preload",
				},
			},
			RequestBody: nil,
			Responses:   responses,
		},
		Method: routing.GET,
		Path: fmt.Sprintf(
			"%s/{id}",
			b.baseUrl,
		),
		Middlewares: middleware,
		Handler: func(ctx fiber.Ctx) error {
			var params GetOneParams

			if err := ctx.Bind().
				URI(&params); err != nil {
				return ctx.Status(fiber.StatusBadRequest).
					JSON(fiber.Map{
						"error":   "Bad Request",
						"message": err.Error(),
					})
			}

			var query GetOneQueryParams

			if err := ctx.Bind().
				Query(&query); err != nil {
				return ctx.Status(fiber.StatusBadRequest).
					JSON(fiber.Map{
						"error":   "Bad Request",
						"message": err.Error(),
					})
			}

			var existingEntity Entity

			var baseQuery = b.storage.Database().Model(&existingEntity)

			for _, preload := range query.Preloads {
				parts := strings.Split(preload, ".")

				for i, part := range parts {
					parts[i] = inflect.Camelize(strings.ToLower(part))
				}

				preload := strings.Join(parts, ".")

				if preload == "" {
					continue
				}

				baseQuery = baseQuery.Preload(preload)
			}

			if err := baseQuery.
				Where("id = ?", params.Id).
				First(&existingEntity).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					return ctx.Status(fiber.StatusNotFound).
						JSON(fiber.Map{
							"error": "Not Found",
							"message": fmt.Sprintf(
								"The %s was not found.",
								strings.ToLower(b.name),
							),
						})
				}

				return ctx.Status(fiber.StatusInternalServerError).
					JSON(fiber.Map{
						"error":   "Internal Server Error",
						"message": err.Error(),
					})
			}

			return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
				"item": existingEntity,
			})
		},
	}
}
