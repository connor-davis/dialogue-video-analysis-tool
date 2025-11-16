package middleware

import (
	"slices"
	"strings"

	"github.com/connor-davis/dialogue-video-analysis-tool/internal/models"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/session"
	"gorm.io/gorm"
)

func (m *middleware) Authorized(permissions ...string) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		session := session.FromContext(ctx)

		currentUserId := session.Get("user_id")

		if currentUserId == nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Unauthorized",
				"message": "You must be logged in to access this resource.",
			})
		}

		var currentUser *models.User

		if err := m.storage.Database().Where("id = ?", currentUserId).Preload("Roles").First(&currentUser).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error":   "Unauthorized",
					"message": "You must be logged in to access this resource.",
				})
			}

			log.Errorf("🔥 Failed to retrieve user from database: %s", err.Error())

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Internal Server Error",
				"message": err.Error(),
			})
		}

		combinedPermissions := []string{}

		for _, role := range currentUser.Roles {
			combinedPermissions = append(combinedPermissions, role.Permissions...)
		}

		if slices.Contains(combinedPermissions, "*") {
			return ctx.Next()
		}

		for _, permission := range combinedPermissions {
			strippedPermission := strings.TrimSuffix(permission, ".*")

			for _, requiredPermission := range permissions {
				if strings.HasPrefix(requiredPermission, strippedPermission) {
					return ctx.Next()
				}
			}
		}

		return ctx.Status(fiber.StatusForbidden).
			JSON(&fiber.Map{
				"error":   "Forbidden",
				"message": "You do not have permission to access this resource.",
			})
	}
}
