package sessions

import (
	"time"

	"github.com/connor-davis/dialogue-video-analysis-tool/common"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/gofiber/storage/postgres/v3"
)

func New() *session.Store {
	return session.NewStore(session.Config{
		Storage: postgres.New(postgres.Config{
			Table:         "sessions",
			ConnectionURI: common.EnvString("DATABASE_DSN", "host=localhost user=<user> password=<password> dbname=<database> port=5432 sslmode=disable TimeZone=Africa/Johannesburg"),
		}),
		CookieDomain:      common.EnvString("API_COOKIE_DOMAIN", "localhost"),
		CookiePath:        "/",
		CookieSameSite:    "Strict",
		CookieSecure:      true,
		CookieSessionOnly: false,
		CookieHTTPOnly:    false,
		Extractor:         extractors.FromCookie(common.EnvString("API_SESSION_COOKIE", "one_session")),
		IdleTimeout:       1 * time.Hour,
		AbsoluteTimeout:   1 * time.Hour,
	})
}
