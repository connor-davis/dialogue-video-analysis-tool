package assignApi

import (
	"github.com/connor-davis/dialogue-video-analysis-tool/internal/routing"
	"github.com/connor-davis/dialogue-video-analysis-tool/internal/storage"
	"github.com/gofiber/fiber/v3"
)

type AssignmentApi[ParentEntity any, ChildEntity any] interface {
	AssignRoute(middleware ...fiber.Handler) routing.Route
	AssignWithPayloadRoute(requestBodyRef string, middleware ...fiber.Handler) routing.Route
	UnassignRoute(middleware ...fiber.Handler) routing.Route
	ListRoute(middleware ...fiber.Handler) routing.Route
}

type assignmentApi[ParentEntity any, ChildEntity any] struct {
	storage    storage.Storage
	baseUrl    string
	parentName string
	childName  string
}

func New[ParentEntity any, ChildEntity any](storage storage.Storage, baseUrl string, parentName string, childName string) AssignmentApi[ParentEntity, ChildEntity] {
	return &assignmentApi[ParentEntity, ChildEntity]{
		storage:    storage,
		baseUrl:    baseUrl,
		parentName: parentName,
		childName:  childName,
	}
}
