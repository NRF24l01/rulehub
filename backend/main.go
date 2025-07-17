package main

import (
	"rulehub/handlers"
	"rulehub/middleware"
	"rulehub/models"
	"rulehub/routes"
	"rulehub/schemas"

	"github.com/go-playground/validator"

	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/labstack/echo/v4"
	echoMw "github.com/labstack/echo/v4/middleware"
)
func main() {
	if os.Getenv("RUNTIME_PRODUCTION") != "true" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("failed to load .env: %v", err)
		}
	}

	db := models.RegisterPostgres()

	validater := validator.New()
	schemas.RegisterCustomValidations(validater)

	e := echo.New()

	e.Validator = &middleware.CustomValidator{Validator: validater}

	e.Use(echoMw.Logger())
    e.Use(echoMw.Recover())

	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(200, schemas.Message{Status: "RuleHUB backend is ok"})
	})

	handler := &handlers.Handler{DB: db}
	routes.RegisterRoutes(e, handler)
	
	e.Logger.Fatal(e.Start(":1324"))
}