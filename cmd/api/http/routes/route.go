package routes

import "github.com/connor-davis/dialogue-video-analysis-tool/internal/routing"

type Router interface {
	LoadRoutes() []routing.Route
}
