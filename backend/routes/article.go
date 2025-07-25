package routes

import (
	"rulehub/handlers"
	"rulehub/middleware"
	"rulehub/schemas"

	"github.com/labstack/echo/v4"
)

func RegisterArticleRoutes(e *echo.Echo, h* handlers.Handler) {
	group := e.Group("/articles")

	group.POST("/", h.ArticleCreateHandler, middleware.ValidationMiddleware(func() interface{} {
		return &schemas.ArticleCreateRequest{}
	}), middleware.JWTMiddleware())

	group.GET("/:uuid", h.ArticleGetHandler)
	group.PUT("/:uuid", h.ArticleChangeHandler, middleware.ValidationMiddleware(func() interface{} {
		return &schemas.ArticleUpdateRequest{}
	}), middleware.JWTMiddleware())
}