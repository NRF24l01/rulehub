package handlers

import (
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
	MinIOClient *minio.Client
}
