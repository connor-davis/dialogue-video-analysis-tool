package middleware

import (
	"github.com/connor-davis/dialogue-video-analysis-tool/internal/storage"
	"github.com/gofiber/fiber/v3"
)

type Middleware interface {
	Authenticated() fiber.Handler
	Authorized(permissions ...string) fiber.Handler
	// Policies(policies ...models.PolicyType) fiber.Handler
}

type middleware struct {
	storage storage.Storage
}

func New(storage storage.Storage) Middleware {
	return &middleware{
		storage: storage,
	}
}
