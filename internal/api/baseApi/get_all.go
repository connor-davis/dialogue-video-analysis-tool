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

type GetAllQueryParams struct {
	DisablePagination bool           `query:"disablePagination"`
	Page              int            `query:"page"`
	PageSize          int            `query:"pageSize"`
	Preloads          pq.StringArray `query:"preload"`
	SearchTerm        string         `query:"searchTerm"`
	SearchColumns     pq.StringArray `query:"searchColumn"`
}

func (b *baseApi[Entity]) GetAllRoute(middleware ...fiber.Handler) routing.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription(fmt.Sprintf("%s retrieved successfully.", inflect.Pluralize(b.name))).
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
				inflect.Pluralize(b.name),
			),
			Description: fmt.Sprintf(
				"This endpoint retrieves all %s.",
				strings.ToLower(inflect.Pluralize(b.name)),
			),
			Tags: []string{fmt.Sprintf(
				"%s",
				inflect.Pluralize(b.name),
			)},
			Parameters: []*openapi3.ParameterRef{
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
			"%s",
			b.baseUrl,
		),
		Middlewares: middleware,
		Handler: func(ctx fiber.Ctx) error {
			var query GetAllQueryParams

			if err := ctx.Bind().
				Query(&query); err != nil {
				return ctx.Status(fiber.StatusBadRequest).
					JSON(fiber.Map{
						"error":   "Bad Request",
						"message": err.Error(),
					})
			}

			var existingEntities []Entity

			var baseQuery = b.storage.Database().Model(&existingEntities)

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

			if query.SearchTerm != "" && len(query.SearchColumns) > 0 {
				var searchConditions []string
				var searchValues []interface{}

				for _, column := range query.SearchColumns {
					searchConditions = append(searchConditions, fmt.Sprintf("%s ILIKE ?", column))
					searchValues = append(searchValues, fmt.Sprintf("%%%s%%", query.SearchTerm))
				}

				baseQuery = baseQuery.Where(strings.Join(searchConditions, " OR "), searchValues...)
			}

			totalEntities := int64(0)

			if err := baseQuery.Count(&totalEntities).Error; err != nil {
				return ctx.Status(fiber.StatusInternalServerError).
					JSON(fiber.Map{
						"error":   "Internal Server Error",
						"message": err.Error(),
					})
			}

			if query.Page < 1 {
				query.Page = 1
			}

			if query.PageSize < 1 {
				query.PageSize = 10
			}

			if query.Page > 0 && query.PageSize > 0 && !query.DisablePagination {
				baseQuery = baseQuery.Offset((query.Page - 1) * query.PageSize).Limit(query.PageSize)
			}

			totalPages := 1

			if totalEntities > 0 {
				totalPages = int((totalEntities + int64(query.PageSize) - 1) / int64(query.PageSize))
			}

			if query.Page > totalPages {
				query.Page = totalPages
			}

			previousPage := max(query.Page-1, 1)
			nextPage := min(query.Page+1, totalPages)

			if err := baseQuery.
				Find(&existingEntities).Error; err != nil {
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
				"items": existingEntities,
				"pagination": fiber.Map{
					"count":        totalEntities,
					"pages":        totalPages,
					"pageSize":     query.PageSize,
					"currentPage":  query.Page,
					"nextPage":     nextPage,
					"previousPage": previousPage,
				},
			})
		},
	}
}
