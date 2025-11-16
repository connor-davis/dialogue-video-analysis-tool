package assignApi

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/connor-davis/dialogue-video-analysis-tool/internal/routing"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-openapi/inflect"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func (a *assignmentApi[ParentEntity, ChildEntity]) UnassignRoute(middleware ...fiber.Handler) routing.Route {
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

	return routing.Route{
		OpenAPIMetadata: routing.OpenAPIMetadata{
			Summary: fmt.Sprintf(
				"Unassign %s",
				a.childName,
			),
			Description: fmt.Sprintf(
				"This endpoint unassigns a %s from a %s",
				strings.ToLower(a.childName),
				strings.ToLower(a.parentName),
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
				{
					Value: openapi3.NewPathParameter(fmt.Sprintf(
						"%sId",
						inflect.Parameterize(a.childName),
					)).
						WithRequired(true).
						WithSchema(openapi3.NewUUIDSchema()),
				},
			},
			RequestBody: nil,
			Responses:   responses,
		},
		Method: routing.POST,
		Path: fmt.Sprintf(
			"%s/unassign-%s/{%sId}/{%sId}",
			a.baseUrl,
			strings.ToLower(inflect.Dasherize(a.childName)),
			inflect.Parameterize(a.parentName),
			inflect.Parameterize(a.childName),
		),
		Middlewares: middleware,
		Handler: func(ctx fiber.Ctx) error {
			parentId := ctx.Params(fmt.Sprintf(
				"%sId",
				inflect.Parameterize(a.parentName),
			))
			childId := ctx.Params(fmt.Sprintf(
				"%sId",
				inflect.Parameterize(a.childName),
			))

			var parentEntity ParentEntity
			var childEntity ChildEntity

			if err := a.storage.Database().
				Model(&parentEntity).
				Where("id = ?", parentId).
				First(&parentEntity).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
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

			if err := a.storage.Database().
				Model(&childEntity).
				Where("id = ?", childId).
				First(&childEntity).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					return ctx.Status(fiber.StatusNotFound).
						JSON(fiber.Map{
							"error": "Not Found",
							"message": fmt.Sprintf(
								"The %s was not found.",
								strings.ToLower(a.childName),
							),
						})
				}

				return ctx.Status(fiber.StatusInternalServerError).
					JSON(fiber.Map{
						"error":   "Internal Server Error",
						"message": err.Error(),
					})
			}

			var existingAssociation ChildEntity

			if err := a.storage.Database().
				Model(&parentEntity).
				Association(fmt.Sprintf(
					"%s",
					inflect.Pluralize(a.childName),
				)).
				Find(&existingAssociation, "id = ?", childId); err != nil {
				return ctx.Status(fiber.StatusInternalServerError).
					JSON(fiber.Map{
						"error":   "Internal Server Error",
						"message": err.Error(),
					})
			}

			if reflect.ValueOf(existingAssociation).IsZero() {
				return ctx.Status(fiber.StatusBadRequest).
					JSON(fiber.Map{
						"error": "Bad Request",
						"message": fmt.Sprintf(
							"The %s is not assigned to the %s.",
							strings.ToLower(a.childName),
							strings.ToLower(a.parentName),
						),
					})
			}

			if err := a.storage.Database().
				Model(&parentEntity).
				Association(fmt.Sprintf(
					"%s",
					inflect.Pluralize(a.childName),
				)).
				Delete(&childEntity); err != nil {
				return ctx.Status(fiber.StatusInternalServerError).
					JSON(fiber.Map{
						"error":   "Internal Server Error",
						"message": err.Error(),
					})
			}

			return ctx.SendStatus(fiber.StatusOK)
		},
	}
}
