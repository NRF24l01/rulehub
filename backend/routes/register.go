package routes

import (
	"rulehub/handlers"

	"log"
	"os"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, h *handlers.Handler) {
	RegisterAuthRoutes(e, h)
	RegisterArticleRoutes(e, h)

	if os.Getenv("RUNTIME_PRODUCTION") != "true" || os.Getenv("TEST_ENV") == "true" {
		log.Println("Registering debug endpoints (development mode)")
		RegisterDevRoutes(e, h)
	}
}