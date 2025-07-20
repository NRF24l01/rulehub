package handlers

import (
	"net/http"
	"rulehub/schemas"

	"log"

	"github.com/labstack/echo/v4"
)

func (h* Handler) DropDB(c echo.Context) error {
	log.Println("Dropping the database...")

	// Disable triggers (including foreign key checks) for all tables
	if err := h.DB.Exec("SET session_replication_role = 'replica';").Error; err != nil {
		log.Printf("Error disabling triggers: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Internal server error"})
	}

	// Delete all data from all tables
	// Получаем все таблицы из базы данных динамически
	var tables []string
	if err := h.DB.Raw("SELECT tablename FROM pg_tables WHERE schemaname = 'public'").Scan(&tables).Error; err != nil {
		log.Printf("Error fetching table names: %v", err)
		// Re-enable triggers before returning
		h.DB.Exec("SET session_replication_role = 'origin';")
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Internal server error"})
	}

	for _, table := range tables {
		if err := h.DB.Exec("DELETE FROM \"" + table + "\"").Error; err != nil {
			log.Printf("Error deleting data from %s table: %v", table, err)
			// Re-enable triggers before returning
			h.DB.Exec("SET session_replication_role = 'origin';")
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Internal server error"})
		}
	}

	// Re-enable triggers
	if err := h.DB.Exec("SET session_replication_role = 'origin';").Error; err != nil {
		log.Printf("Error re-enabling triggers: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Internal server error"})
	}

	log.Println("Database reset successfully")
	return c.JSON(http.StatusOK, schemas.Message{
		Status: "Database reset successfully",
	})
}