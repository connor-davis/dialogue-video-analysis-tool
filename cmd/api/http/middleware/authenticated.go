package middleware

import (
	"time"

	"github.com/connor-davis/dialogue-video-analysis-tool/internal/models"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/session"
	"gorm.io/gorm"
)

func (m *middleware) Authenticated() fiber.Handler {
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

		session.Session.SetIdleTimeout(1 * time.Hour)

		if err := session.Session.Save(); err != nil {
			log.Errorf("🔥 Failed to save session: %s", err.Error())

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Internal Server Error",
				"message": "Failed to save session.",
			})
		}

		ctx.Locals("user_id", currentUser.Id)
		ctx.Locals("user", currentUser)

		return ctx.Next()
	}
}
