package main

import (
	"fmt"
	"time"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/connor-davis/dialogue-video-analysis-tool/cmd/api/http"
	"github.com/connor-davis/dialogue-video-analysis-tool/cmd/api/http/middleware"
	"github.com/connor-davis/dialogue-video-analysis-tool/common"
	"github.com/connor-davis/dialogue-video-analysis-tool/internal/storage"
	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/gofiber/storage/postgres/v3"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func main() {
	storage := storage.New()

	if err := storage.Migrate(); err != nil {
		log.Errorf("🔥 Failed to migrate models: %s", err.Error())
	}

	openai := openai.NewClient(
		option.WithAPIKey(common.EnvString("OPENAI_API_KEY", "sk...")), // or set OPENAI_API_KEY in your env
	)

	middleware := middleware.New(storage)

	app := fiber.New(fiber.Config{
		AppName:       "One REST API",
		ServerHeader:  "One REST API",
		JSONEncoder:   json.Marshal,
		JSONDecoder:   json.Unmarshal,
		StrictRouting: true,
		CaseSensitive: true,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			common.EnvString(
				"APP_BASE_URL",
				"http://localhost:3000",
			),
			"http://localhost:3000",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowCredentials: true,
	}))

	app.Use(logger.New(logger.Config{
		TimeFormat: "2006-01-01 00:00:00",
		TimeZone:   "Africa/Johannesburg",
	}))

	app.Use(session.New(session.Config{
		Storage: postgres.New(postgres.Config{
			Table:         "sessions",
			ConnectionURI: common.EnvString("DATABASE_DSN", "host=localhost user=<user> password=<password> dbname=<database> port=5432 sslmode=disable TimeZone=Africa/Johannesburg"),
		}),
		CookieDomain:      common.EnvString("API_COOKIE_DOMAIN", "localhost"),
		CookiePath:        "/",
		CookieSameSite:    "Lax",
		CookieSecure:      true,
		CookieSessionOnly: false,
		CookieHTTPOnly:    false,
		Extractor:         extractors.FromCookie(common.EnvString("API_SESSION_COOKIE", "one_session")),
		IdleTimeout:       1 * time.Hour,
		AbsoluteTimeout:   1 * time.Hour,
	}))

	apiv1 := app.Group("/api/v1")

	apiv1.Use("/ws", func(ctx fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(ctx) {
			ctx.Locals("ws-allowed", true)

			return ctx.Next()
		}

		return fiber.ErrUpgradeRequired
	})

	httpRouter := http.New(storage, middleware, openai)
	httpRouter.InitializeRoutes(apiv1)

	openapi := httpRouter.InitializeOpenAPI()

	apiv1.Get("/api-spec", func(ctx fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).JSON(openapi)
	})

	apiv1.Get("/api-docs", func(ctx fiber.Ctx) error {
		html, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecURL: func() string {
				if common.EnvString(
					"API_MODE",
					"development",
				) == "production" {
					return fmt.Sprintf(
						"%s/api/v1/api-spec",
						common.EnvString("API_BASE_URL", "https://example.com"),
					)
				}

				return fmt.Sprintf(
					"http://localhost:%s/api/v1/api-spec",
					common.EnvString("API_PORT", "6173"),
				)
			}(),
			Theme:  scalar.ThemeDefault,
			Layout: scalar.LayoutModern,
			BaseServerURL: func() string {
				if common.EnvString(
					"API_MODE",
					"development",
				) == "production" {
					return common.EnvString(
						"API_BASE_URL",
						"https://example.com",
					)
				}

				return fmt.Sprintf(
					"http://localhost:%s",
					common.EnvString("API_PORT", "6173"),
				)
			}(),
			DarkMode: true,
		})

		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).
				SendString(err.Error())
		}

		return ctx.Status(fiber.StatusOK).
			Type("html").
			SendString(html)
	})

	log.Infof(
		"✅ Starting API on port %s...",
		common.EnvString("API_PORT", "6173"),
	)

	if err := app.Listen(
		fmt.Sprintf(
			":%s",
			common.EnvString(
				"API_PORT",
				"6173",
			),
		),
		fiber.ListenConfig{
			EnablePrintRoutes: true,
		},
	); err != nil {
		log.Errorf("🔥 Failed to start server: %v", err)
	}
}
