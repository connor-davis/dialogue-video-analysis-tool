package totp

import (
	"github.com/connor-davis/dialogue-video-analysis-tool/internal/models"
	"github.com/connor-davis/dialogue-video-analysis-tool/internal/routing"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/pquerna/otp/totp"
)

func (r *TotpRouter) VerifyRoute() routing.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription("The MFA code has been verified successfully.").
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
			Summary:     "Verify TOTP MFA",
			Description: "Verifies the Time-based One-Time Password (TOTP) Multi-Factor Authentication code for the user.",
			Tags:        []string{"Authentication"},
			Parameters: []*openapi3.ParameterRef{
				{
					Ref: "#/components/parameters/Code",
				},
			},
			RequestBody: nil,
			Responses:   responses,
		},
		Method: routing.POST,
		Path:   "/authentication/mfa/totp/verify",
		Middlewares: []fiber.Handler{
			r.middleware.Authenticated(),
		},
		Handler: func(ctx fiber.Ctx) error {
			currentUser := ctx.Locals("user").(*models.User)

			var queryParams struct {
				Code string `query:"code"`
			}

			if err := ctx.Bind().Query(&queryParams); err != nil {
				log.Infof("🔥 Error parsing query parameters: %s", err.Error())

				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":   "Bad Request",
					"message": "Invalid query parameters.",
				})
			}

			if queryParams.Code == "" || len(queryParams.Code) < 6 || len(queryParams.Code) > 6 {
				log.Warn("🚫 Unauthorized access attempt: No MFA code provided")

				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":   "Bad Request",
					"message": "Unable to verify Multi-Factor Authentication (MFA) status. Please provide a valid MFA code.",
				})
			}

			if currentUser == nil || currentUser.MfaSecret == nil {
				log.Warn("🚫 Unauthorized access attempt: User not found or MFA not enabled")

				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":   "Bad Request",
					"message": "Unable to verify Multi-Factor Authentication (MFA) status. Please ensure MFA is enabled for your account.",
				})
			}

			if !totp.Validate(queryParams.Code, string(currentUser.MfaSecret)) {
				return ctx.Status(fiber.StatusUnauthorized).
					JSON(fiber.Map{
						"error":   "Unauthorized",
						"message": "Invalid Multi-Factor Authentication code. Please try again.",
					})
			}

			currentUser.MfaEnabled = true
			currentUser.MfaVerified = true

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
