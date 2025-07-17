package routes

import (
	"rulehub/handlers"
	"rulehub/middleware"
	"rulehub/schemas"

	"github.com/labstack/echo/v4"
)

func RegisterAuthRoutes(e *echo.Echo, h* handlers.Handler) {
	group := e.Group("/auth")

	group.POST("/login", h.UserLoginHandler, middleware.ValidationMiddleware(func() interface{} {
		return &schemas.SignInRequest{}
	}))

	group.POST("/register", h.UserRegistrationHandler, middleware.ValidationMiddleware(func() interface{} {
		return &schemas.SignUpRequest{}
	}))

	group.POST("/refresh", h.UserJwtRefreshHandler)
}