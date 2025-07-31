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

// Gen static get url by uuid
func (h *Handler) MediaGetURLHandler(c echo.Context) error {
	// Get bucket name from environment variable
	bucketName := os.Getenv("MINIO_BUCKET")
	if bucketName == "" {
		log.Printf("MINIO_BUCKET environment variable not set")
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Server configuration error"})
	}

	// Get the file ID from the request parameters
	fileID := c.Param("static_url")
	if fileID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "File ID is required"})
	}

	// Generate a presigned GET URL for the file
	presignedURL := utils.GetPermanentObjectURL(bucketName, fileID)

	return c.JSON(http.StatusOK, echo.Map{"url": presignedURL})
}