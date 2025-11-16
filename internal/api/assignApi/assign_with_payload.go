package assignApi

import (
	"fmt"
	"strings"

	"github.com/connor-davis/dialogue-video-analysis-tool/internal/routing"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-openapi/inflect"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func (a *assignmentApi[ParentEntity, ChildEntity]) AssignWithPayloadRoute(requestBodyRef string, middleware ...fiber.Handler) routing.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription(fmt.Sprintf("%s assigned to %s successfully.", a.childName, a.parentName)).
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

	return routing.Route{OpenAPIMetadata: routing.OpenAPIMetadata{
		Summary: fmt.Sprintf(
			"Assign %s With Payload",
			a.childName,
		),
		Description: fmt.Sprintf(
			"This endpoint assigns a %s to a %s with the %s payload.",
			strings.ToLower(a.childName),
			strings.ToLower(a.parentName),
			strings.ToLower(a.childName),
		),
		Tags: []string{fmt.Sprintf(
			"%s",
			inflect.Pluralize(a.parentName),
		)},
		Parameters: []*openapi3.ParameterRef{
			{
				Value: openapi3.NewPathParameter(fmt.Sprintf(
					"%sId",
					inflect.Parameterize(a.parentName),
				)).
					WithRequired(true).
					WithSchema(openapi3.NewUUIDSchema()),
			},
		},
		RequestBody: &openapi3.RequestBodyRef{
			Ref: requestBodyRef,
		},
		Responses: responses,
	}, Method: routing.POST, Path: fmt.Sprintf(
		"%s/assign-%s/{%sId}",
		a.baseUrl,
		strings.ToLower(inflect.Dasherize(a.childName)),
		inflect.Parameterize(a.parentName),
	), Middlewares: middleware, Handler: func(ctx fiber.Ctx) error {
		parentId := ctx.Params(fmt.Sprintf(
			"%sId",
			inflect.Parameterize(a.parentName),
		))

		var parentEntity ParentEntity
		var childEntity ChildEntity

		if err := a.storage.Database().
			Model(&parentEntity).
			Where("id = ?", parentId).
			First(&parentEntity).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return ctx.Status(fiber.StatusNotFound).
					JSON(fiber.Map{
						"error": "Not Found",
						"message": fmt.Sprintf(
							"The %s was not found.",
							strings.ToLower(a.parentName),
						),
					})
			}

			return ctx.Status(fiber.StatusInternalServerError).
				JSON(fiber.Map{
					"error":   "Internal Server Error",
					"message": err.Error(),
				})
		}

		if err := ctx.Bind().
			Body(&childEntity); err != nil {
			return ctx.Status(fiber.StatusBadRequest).
				JSON(fiber.Map{
					"error":   "Bad Request",
					"message": "Invalid request body.",
				})
		}

		if err := a.storage.Database().
			Session(&gorm.Session{
				FullSaveAssociations: true,
			}).
			Model(&parentEntity).
			Association(fmt.Sprintf(
				"%s",
				inflect.Pluralize(a.childName),
			)).
			Append(&childEntity); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).
				JSON(fiber.Map{
					"error":   "Internal Server Error",
					"message": err.Error(),
				})
		}

		return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
			"item": childEntity,
		})
	}}
}
