package handlers

import (
	"net/http"
	"rulehub/schemas"

	"log"

	"github.com/labstack/echo/v4"
)

func (h* Handler) DropDB(c echo.Context) error {
	log.Println("Dropping the database...")

	// Delete all data from all tables
	// Получаем все таблицы из базы данных динамически
	var tables []string
	if err := h.DB.Raw("SELECT tablename FROM pg_tables WHERE schemaname = 'public'").Scan(&tables).Error; err != nil {
		log.Printf("Error fetching table names: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Internal server error"})
	}

	for _, table := range tables {
		if err := h.DB.Exec("DELETE FROM \"" + table + "\"").Error; err != nil {
			log.Printf("Error deleting data from %s table: %v", table, err)
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Internal server error"})
		}
	}
	log.Println("Database reset successfully")
	return c.JSON(http.StatusOK, schemas.Message{
		Status: "Database reset successfully",
	})
}