package handlers

import (
	"log"
	"net/http"
	"os"
	"rulehub/schemas"
	"rulehub/utils"

	"github.com/labstack/echo/v4"
)

// MediaUploadTempHandler generates a presigned URL for temporary file uploads
func (h *Handler) MediaUploadTempHandler(c echo.Context) error {
	// Get bucket name from environment variable
	bucketName := os.Getenv("MINIO_BUCKET")
	if bucketName == "" {
		log.Printf("MINIO_BUCKET environment variable not set")
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Server configuration error"})
	}

	// Get the presigned URL expiration time
	expires := utils.GetPresignedLifetime()

	// Generate a presigned PUT URL for temporary upload
	fileID, presignedURL, err := utils.GeneratePresignedPutURL(h.MinIOClient, bucketName, expires)
	if err != nil {
		log.Printf("Error generating presigned URL: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Error generating upload URL"})
	}

	// Return the presigned URL and file ID to the client
	resp := schemas.MediaUploadResponse{
		TempURL: presignedURL,
		FileID:  fileID,
	}

	return c.JSON(http.StatusOK, resp)
}