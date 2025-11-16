package baseApi

import (
	"fmt"
	"strings"

	"github.com/connor-davis/dialogue-video-analysis-tool/internal/routing"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-openapi/inflect"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type UpdateParams struct {
	Id string `json:"id" msg:"id"`
}

func (b *baseApi[Entity]) UpdateRoute(requestBodyRef string, middleware ...fiber.Handler) routing.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription(fmt.Sprintf("%s updated successfully.", b.name)).
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
				"Update %s",
				b.name,
			),
			Description: fmt.Sprintf(
				"This endpoint updates an existing %s by their id.",
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
			},
			RequestBody: &openapi3.RequestBodyRef{
				Ref: requestBodyRef,
			},
			Responses: responses,
		},
		Method: routing.PUT,
		Path: fmt.Sprintf(
			"%s/{id}",
			b.baseUrl,
		),
		Middlewares: middleware,
		Handler: func(ctx fiber.Ctx) error {
			var params UpdateParams

			if err := ctx.Bind().
				URI(&params); err != nil {
				return ctx.Status(fiber.StatusBadRequest).
					JSON(fiber.Map{
						"error":   "Bad Request",
						"message": err.Error(),
					})
			}

			var entity map[string]any

			if err := ctx.Bind().
				Body(&entity); err != nil {
				return ctx.Status(fiber.StatusBadRequest).
					JSON(fiber.Map{
						"error":   "Bad Request",
						"message": err.Error(),
					})
			}

			for key, value := range entity {
				entity[inflect.Underscore(key)] = value

				if key != inflect.Underscore(key) {
					delete(entity, key)
				}
			}

			var existingEntity Entity

			if err := b.storage.Database().
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

			if err := b.storage.Database().
				Model(&existingEntity).
				Updates(&entity).Error; err != nil {
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
				"item": entity,
			})
		},
	}
}
