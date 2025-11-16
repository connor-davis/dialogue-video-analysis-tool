package assignApi

import (
	"fmt"
	"math"
	"strings"

	"github.com/connor-davis/dialogue-video-analysis-tool/internal/routing"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-openapi/inflect"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ListParams struct {
	ParentId uuid.UUID `param:"parentId"`
	ChildId  uuid.UUID `param:"childId"`
}

type ListQueryParams struct {
	DisablePagination bool           `query:"disablePagination"`
	Page              int            `query:"page"`
	PageSize          int            `query:"pageSize"`
	Preloads          pq.StringArray `query:"preload"`
	SearchTerm        string         `query:"searchTerm"`
	SearchColumns     pq.StringArray `query:"searchColumn"`
}

func (a *assignmentApi[ParentEntity, ChildEntity]) ListRoute(middleware ...fiber.Handler) routing.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription(fmt.Sprintf("%s %s retrieved successfully.", a.parentName, inflect.Pluralize(a.childName))).
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
				"List %s",
				inflect.Pluralize(a.childName),
			),
			Description: fmt.Sprintf(
				"This endpoint retrieves a list of %s assigned to a %s",
				strings.ToLower(a.childName),
				strings.ToLower(a.parentName),
			),
			Tags: []string{fmt.Sprintf(
				"%s",
				inflect.Pluralize(a.parentName),
			)},
			Parameters: []*openapi3.ParameterRef{
				{
					Value: openapi3.NewPathParameter(fmt.Sprintf("%sId", inflect.Parameterize(a.parentName))).
						WithRequired(true).
						WithSchema(openapi3.NewUUIDSchema()),
				},
				{
					Ref: "#/components/parameters/DisablePagination",
				},
				{
					Ref: "#/components/parameters/Page",
				},
				{
					Ref: "#/components/parameters/PageSize",
				},
				{
					Ref: "#/components/parameters/Preload",
				},
				{
					Ref: "#/components/parameters/SearchTerm",
				},
				{
					Ref: "#/components/parameters/SearchColumn",
				},
			},
			RequestBody: nil,
			Responses:   responses,
		},
		Method: routing.GET,
		Path: fmt.Sprintf(
			"%s/{%sId}/list-%s",
			a.baseUrl,
			inflect.Parameterize(a.parentName),
			strings.ToLower(inflect.Dasherize(inflect.Pluralize(a.childName))),
		),
		Middlewares: middleware,
		Handler: func(ctx fiber.Ctx) error {
			parentId := ctx.Params(fmt.Sprintf(
				"%sId",
				inflect.Parameterize(a.parentName),
			))

			var queryParams ListQueryParams

			if err := ctx.Bind().Query(&queryParams); err != nil {
				return ctx.Status(fiber.StatusBadRequest).
					JSON(fiber.Map{
						"error":   "Bad Request",
						"message": err.Error(),
					})
			}

			var parentEntity ParentEntity

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

				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error":   "Internal Server Error",
					"message": err.Error(),
				})
			}

			clauses := []clause.Expression{}

			if queryParams.SearchTerm != "" && len(queryParams.SearchColumns) > 0 {
				for _, column := range queryParams.SearchColumns {
					clauses = append(clauses, clause.Expr{
						SQL:  fmt.Sprintf("%s ILIKE ?", column),
						Vars: []any{fmt.Sprintf("%%%s%%", queryParams.SearchTerm)},
					})
				}
			}

			countQuery := a.storage.Database().
				Model(&parentEntity)

			if len(clauses) > 0 {
				countQuery = countQuery.Clauses(clause.Or(clauses...))
			}

			totalEntities := countQuery.
				Association(fmt.Sprintf(
					"%ss",
					a.childName,
				)).
				Count()

			if queryParams.Page < 1 {
				queryParams.Page = 1
			}

			if queryParams.PageSize < 1 {
				queryParams.PageSize = 10
			}

			limit := queryParams.PageSize
			offset := (queryParams.Page - 1) * queryParams.PageSize
			totalPages := int64(math.Ceil(float64(totalEntities) / float64(limit)))

			if totalPages == 0 {
				totalPages = 1
			}

			nextPage := queryParams.Page + 1
			previousPage := queryParams.Page - 1

			if nextPage > int(totalPages) {
				nextPage = int(totalPages)
			}

			if previousPage < 1 {
				previousPage = 1
			}

			var existingAssociations []ChildEntity

			query := a.storage.Database().
				Model(&parentEntity)

			if !queryParams.DisablePagination {
				query = query.Limit(limit).Offset(offset)
			}

			for _, preload := range queryParams.Preloads {
				parts := strings.Split(preload, ".")

				for i, part := range parts {
					parts[i] = inflect.Camelize(strings.ToLower(part))
				}

				preload := strings.Join(parts, ".")

				if preload == "" {
					continue
				}

				query = query.Preload(preload)
			}

			if len(clauses) > 0 {
				query = query.Clauses(clause.Or(clauses...))
			}

			if err := query.
				Association(fmt.Sprintf(
					"%ss",
					a.childName,
				)).
				Find(&existingAssociations); err != nil {
				return ctx.Status(fiber.StatusInternalServerError).
					JSON(fiber.Map{
						"error":   "Internal Server Error",
						"message": err.Error(),
					})
			}

			return ctx.Status(fiber.StatusOK).
				JSON(fiber.Map{
					"items": existingAssociations,
					"pagination": fiber.Map{
						"count":        totalEntities,
						"pages":        totalPages,
						"pageSize":     queryParams.PageSize,
						"currentPage":  queryParams.Page,
						"nextPage":     nextPage,
						"previousPage": previousPage,
					},
				})
		},
	}
}
