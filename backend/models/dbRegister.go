package models

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func RegisterPostgres() (*gorm.DB){
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", 
	os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_SSLMODE"), os.Getenv("POSTGRES_TIMEZONE"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to get db instance: %v", err)
	}

	db.Exec(`CREATE EXTENSION IF NOT EXISTS "pgcrypto";`)

	db.AutoMigrate(&User{}, &Article{}, &Media{})

	return db
}