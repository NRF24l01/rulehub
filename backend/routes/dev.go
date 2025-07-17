package routes

import (
	"rulehub/handlers"

	"github.com/labstack/echo/v4"
)

func RegisterDevRoutes(e *echo.Echo, h *handlers.Handler) {
	group := e.Group("/dev")

	group.POST("/reset-db", h.DropDB)
}