package baseApi

import (
	"github.com/connor-davis/dialogue-video-analysis-tool/internal/routing"
	"github.com/connor-davis/dialogue-video-analysis-tool/internal/storage"
	"github.com/gofiber/fiber/v3"
)

type BaseApi[Entity any] interface {
	CreateRoute(requestBodyRef string, middleware ...fiber.Handler) routing.Route
	UpdateRoute(requestBodyRef string, middleware ...fiber.Handler) routing.Route
	DeleteRoute(middleware ...fiber.Handler) routing.Route
	GetOneRoute(middleware ...fiber.Handler) routing.Route
	GetAllRoute(middleware ...fiber.Handler) routing.Route
	GetAllByFieldRoute(fieldName string, middleware ...fiber.Handler) routing.Route
}

type baseApi[Entity any] struct {
	storage storage.Storage
	baseUrl string
	name    string
}

func New[Entity any](storage storage.Storage, baseUrl string, name string) BaseApi[Entity] {
	return &baseApi[Entity]{
		storage: storage,
		baseUrl: baseUrl,
		name:    name,
	}
}
