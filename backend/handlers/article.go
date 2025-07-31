package handlers

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"rulehub/models"
	"rulehub/schemas"
	"rulehub/utils"
	"strings"

	"github.com/labstack/echo/v4"

	googleUUID "github.com/google/uuid"
)

func (h *Handler) ArticleCreateHandler(c echo.Context) error {
	article_data := c.Get("validatedBody").(*schemas.ArticleCreateRequest)
	log.Printf("Creating article: %+v", article_data)

	var user models.User
	if err := h.DB.Where("id = ?", c.Get("userID").(string)).First(&user).Error; err != nil {
		log.Printf("User not found: %v", err)
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Invalid username or password"})
	}

	article := models.Article{
		Title:   article_data.Title,
		Content: article_data.Content,
		UserID:  c.Get("userID").(string),
	}
	if err := h.DB.Create(&article).Error; err != nil {
		log.Printf("Error creating article: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Internal server error"})
	}

	var mediaResponses []schemas.MediaCreateResponse
	for _, mediaPath := range article_data.Media {
		// Extract the S3 key from the media path (which contains the temporary file location)
		s3Key := extractS3KeyFromPath(mediaPath)

		// Change file status from temporary to permanent
		if err := utils.ChangeObjectStatusToPermanent(h.MinIOClient, os.Getenv("MINIO_BUCKET"), s3Key); err != nil {
			log.Printf("Error changing file status to permanent: %v", err)
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Internal server error"})
		}

		// Get original filename
		originalFileName := getOriginalFileName(mediaPath)

		// Save the permanent file info in database
		media := models.Media{
			FileName:  originalFileName,
			S3Key:     s3Key,
			ArticleID: article.ID.String(),
		}
		if err := h.DB.Create(&media).Error; err != nil {
			log.Printf("Error creating media: %v", err)
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Internal server error"})
		}

		// Generate a permanent URL for the file
		permanentURL := utils.GetPermanentObjectURL(os.Getenv("MINIO_BUCKET"), s3Key)

		mediaResponses = append(mediaResponses, schemas.MediaCreateResponse{
			FileName: originalFileName,
			S3Key:    permanentURL, // Permanent URL for the file
		})
	}

	resp := schemas.ArticleResponse{
		ID:               article.ID.String(),
		Title:            article.Title,
		Content:          article.Content,
		MediaPresignedUrl: mediaResponses,
		AuthorUsername:   user.Username,
	}
	return c.JSON(http.StatusCreated, resp)
}

// extractS3KeyFromPath extracts the S3 key from a file path
func extractS3KeyFromPath(filePath string) string {
	// The path might contain the full URL or just the object key
	// If it's a URL, extract just the object key part
	if strings.HasPrefix(filePath, "http://") || strings.HasPrefix(filePath, "https://") {
		parsedURL, err := url.Parse(filePath)
		if err == nil {
			// Extract the object key from the URL path
			parts := strings.Split(parsedURL.Path, "/")
			// Usually the last part is the object key
			if len(parts) > 0 {
				return parts[len(parts)-1]
			}
		}
	}

	// If it's not a URL or parsing failed, treat it as a direct path
	// Just get the base filename which should be the S3 key
	return filepath.Base(filePath)
}

// getOriginalFileName extracts the original file name
func getOriginalFileName(filePath string) string {
	// Get the base filename first
	baseName := filepath.Base(filePath)

	// If the filename contains metadata like UUID, extract just the original part
	// This implementation depends on your naming convention
	// Example: if format is "UUID_originalname.ext"
	parts := strings.SplitN(baseName, "_", 2)
	if len(parts) > 1 {
		return parts[1]
	}

	return baseName
}

func (h *Handler) ArticleGetHandler(c echo.Context) error {
	uuid := c.Param("uuid")

	if err := googleUUID.Validate(uuid); err != nil {
		log.Printf("Error getting article, bad uuid: %v", uuid)
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "UUID is required"})
	}

	var article models.Article
	if err := h.DB.Preload("User").Preload("Media").Where("id = ?", uuid).First(&article).Error; err != nil {
		log.Printf("Error getting article with id: %v, 404", uuid)
		return c.JSON(http.StatusNotFound, echo.Map{"message": "No such article"})
	}

	var mediaResponses []schemas.MediaCreateResponse
	bucketName := os.Getenv("MINIO_BUCKET")

	for _, media := range article.Media {
		// Generate permanent URL for each media file
		permanentURL := utils.GetPermanentObjectURL(bucketName, media.S3Key)

		mediaResponses = append(mediaResponses, schemas.MediaCreateResponse{
			FileName: media.FileName,
			S3Key:    permanentURL,
		})
	}

	resp := schemas.ArticleResponse{
		ID:               article.ID.String(),
		Title:            article.Title,
		Content:          article.Content,
		MediaPresignedUrl: mediaResponses,
		AuthorUsername:   article.User.Username,
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) ArticleChangeHandler(c echo.Context) error {
	uuid := c.Param("uuid")

	if err := googleUUID.Validate(uuid); err != nil {
		log.Printf("Error updating article, bad uuid: %v", uuid)
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "UUID is required"})
	}

	var article models.Article
	if err := h.DB.Preload("Media").Where("id = ?", uuid).First(&article).Error; err != nil {
		log.Printf("Article not found: %v", uuid)
		return c.JSON(http.StatusNotFound, echo.Map{"message": "No such article"})
	}

	articleData := c.Get("validatedBody").(*schemas.ArticleUpdateRequest)
	log.Printf("Updating article: %+v", articleData)

	if articleData.Title != nil {
		article.Title = *articleData.Title
	}
	if articleData.Content != nil {
		article.Content = *articleData.Content
	}

	var mediaResponses []schemas.MediaCreateResponse

	if articleData.Media != nil {
		// Delete old media from database
		if err := h.DB.Where("article_id = ?", article.ID).Delete(&models.Media{}).Error; err != nil {
			log.Printf("Error deleting old media: %v", err)
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Internal server error"})
		}
		article.Media = nil

		// Process new media files that are already uploaded as temporary
		for _, mediaPath := range *articleData.Media {
			// Extract the S3 key from the media path
			s3Key := extractS3KeyFromPath(mediaPath)

			// Change file status from temporary to permanent
			if err := utils.ChangeObjectStatusToPermanent(h.MinIOClient, os.Getenv("MINIO_BUCKET"), s3Key); err != nil {
				log.Printf("Error changing file status to permanent: %v", err)
				return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Internal server error"})
			}

			// Get original filename
			originalFileName := getOriginalFileName(mediaPath)

			// Save the permanent file info in database
			media := models.Media{
				FileName:  originalFileName,
				S3Key:     s3Key,
				ArticleID: article.ID.String(),
			}
			if err := h.DB.Create(&media).Error; err != nil {
				log.Printf("Error creating media: %v", err)
				return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Internal server error"})
			}
			article.Media = append(article.Media, media)

			// Generate a permanent URL for the file
			permanentURL := utils.GetPermanentObjectURL(os.Getenv("MINIO_BUCKET"), s3Key)

			mediaResponses = append(mediaResponses, schemas.MediaCreateResponse{
				FileName: originalFileName,
				S3Key:    permanentURL, // Permanent URL for the file
			})
		}
	} else {
		// If media was not updated, return permanent URLs for existing files
		for _, media := range article.Media {
			permanentURL := utils.GetPermanentObjectURL(os.Getenv("MINIO_BUCKET"), media.S3Key)

			mediaResponses = append(mediaResponses, schemas.MediaCreateResponse{
				FileName: media.FileName,
				S3Key:    permanentURL, // Permanent URL for the file
			})
		}
	}

	if err := h.DB.Save(&article).Error; err != nil {
		log.Printf("Error updating article: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Internal server error"})
	}

	var user models.User
	if err := h.DB.Where("id = ?", article.UserID).First(&user).Error; err != nil {
		log.Printf("User not found for article: %v", article.UserID)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Internal server error"})
	}

	resp := schemas.ArticleResponse{
		ID:               article.ID.String(),
		Title:            article.Title,
		Content:          article.Content,
		MediaPresignedUrl: mediaResponses,
		AuthorUsername:   user.Username,
	}
	return c.JSON(http.StatusOK, resp)
}