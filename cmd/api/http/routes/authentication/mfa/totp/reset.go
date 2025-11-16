package totp

import (
	"github.com/connor-davis/dialogue-video-analysis-tool/internal/models"
	"github.com/connor-davis/dialogue-video-analysis-tool/internal/routing"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

func (r *TotpRouter) ResetRoute() routing.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription("The MFA code has been reset successfully.").
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
			Summary:     "Reset TOTP MFA",
			Description: "Resets the Time-based One-Time Password (TOTP) Multi-Factor Authentication for the user.",
			Tags:        []string{"Authentication"},
			Parameters: []*openapi3.ParameterRef{
				{
					Ref: "#/components/parameters/Id",
				},
			},
			RequestBody: nil,
			Responses:   responses,
		},
		Method: routing.POST,
		Path:   "/authentication/mfa/totp/reset/{id}",
		Middlewares: []fiber.Handler{
			r.middleware.Authenticated(),
		},
		Handler: func(ctx fiber.Ctx) error {
			userId := ctx.Params("id")

			var currentUser models.User

			if err := r.storage.Database().
				Where("id = ?", userId).
				First(&currentUser).Error; err != nil {
				log.Errorf("🔥 Error fetching user: %s", err.Error())

				return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error":   "Not Found",
					"message": "User not found.",
				})
			}

			currentUser.MfaSecret = nil
			currentUser.MfaEnabled = false
			currentUser.MfaVerified = false

			if err := r.storage.Database().
				Where("id = ?", currentUser.Id).
				Model(&currentUser).
				Updates(map[string]any{
					"mfa_enabled":  currentUser.MfaEnabled,
					"mfa_verified": currentUser.MfaVerified,
				}).Error; err != nil {
				log.Errorf("🔥 Error updating user: %s", err.Error())

				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error":   "Internal Server Error",
					"message": "An error occurred while processing your request.",
				})
			}

			return ctx.SendStatus(fiber.StatusOK)
		},
	}
}
