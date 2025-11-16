package baseApi

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/connor-davis/dialogue-video-analysis-tool/internal/routing"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-openapi/inflect"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func (b *baseApi[Entity]) CreateRoute(requestBodyRef string, middleware ...fiber.Handler) routing.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription(fmt.Sprintf("%s created successfully.", b.name)).
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
				"Create %s",
				b.name,
			),
			Description: fmt.Sprintf(
				"This endpoint creates a new %s.",
				strings.ToLower(b.name),
			),
			Tags: []string{fmt.Sprintf(
				"%s",
				inflect.Pluralize(b.name),
			)},
			Parameters: nil,
			RequestBody: &openapi3.RequestBodyRef{
				Ref: requestBodyRef,
			},
			Responses: responses,
		},
		Method: routing.POST,
		Path: fmt.Sprintf(
			"%s",
			b.baseUrl,
		),
		Middlewares: middleware,
		Handler: func(ctx fiber.Ctx) error {
			var entity *Entity

			if err := ctx.Bind().
				Body(&entity); err != nil {
				return ctx.Status(fiber.StatusBadRequest).
					JSON(fiber.Map{
						"error":   "Bad Request",
						"message": "Invalid request body.",
					})
			}

			id, err := uuid.NewUUID()

			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).
					JSON(fiber.Map{
						"error":   "Internal Server Error",
						"message": "Could not generate UUID.",
					})
			}

			entityValue := reflect.ValueOf(entity).Elem()
			entityValueId := entityValue.FieldByName("Id")

			entityValueId.Set(reflect.ValueOf(id))

			if err := b.storage.Database().
				Create(&entity).Error; err != nil {
				return ctx.Status(fiber.StatusInternalServerError).
					JSON(fiber.Map{
						"error":   "Internal Server Error",
						"message": "Could not create entity.",
					})
			}

			return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
				"item": entity,
			})
		},
	}
}
