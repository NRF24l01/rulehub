package routes

import (
	"rulehub/handlers"
	"rulehub/middleware"

	"github.com/labstack/echo/v4"
)

func RegisterMediaRoutes(e *echo.Echo, h* handlers.Handler) {
	group := e.Group("/media")

	group.POST("/upload-temp", h.MediaUploadTempHandler, middleware.JWTMiddleware())
}