package handlers

import (
	"log"
	"net/http"
	"os"
	"rulehub/models"
	"rulehub/schemas"
	"rulehub/utils"
	"time"

	"github.com/labstack/echo/v4"

	googleUUID "github.com/google/uuid"
)

func (h* Handler) ArticleCreateHandler(c echo.Context) error {
	article_data := c.Get("validatedBody").(*schemas.ArticleCreateRequest)
	log.Printf("Creating article: %+v", article_data)

	var user models.User
	if err := h.DB.Where("id = ?", c.Get("userID").(string)).First(&user).Error; err != nil {
		log.Printf("User not found: %v", err)
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Invalid username or password"})
	}

	// Create a new article in the database
	article := models.Article{
		Title:   article_data.Title,
		Content: article_data.Content,
		UserID:  c.Get("userID").(string),
	}

	// Create media entries if any
	var mediaResponses []schemas.MediaCreateResponse
	for _, mediaFileName := range article_data.Media {
		name, URL, err := utils.FullGeneratePresignedURL(h.MinIOClient, os.Getenv("MINIO_BUCKET"), 1*time.Hour)
		if err != nil {
			log.Printf("Error generating presigned URL: %v", err)
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Internal server error"})
		}

		media := models.Media{
			FileName: mediaFileName,
			S3Key:    name,
			ArticleID: article.ID.String(),
			Article:   article,
		}
		mediaResponses = append(mediaResponses, schemas.MediaCreateResponse{
			FileName: mediaFileName,
			S3Key:    URL,
		})
		if err := h.DB.Create(&media).Error; err != nil {
			log.Printf("Error creating media: %v", err)
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Internal server error"})
		}
	}

	if err := h.DB.Create(&article).Error; err != nil {
		log.Printf("Error creating article: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Internal server error"})
	}

	resp := schemas.ArticleResponse{
		ID:             article.ID.String(),
		Title:          article.Title,
		Content:        article.Content,
		MediaPresignedUrl: mediaResponses,
		AuthorUsername: user.Username,
	}
	return c.JSON(http.StatusCreated, resp)
}

func (h* Handler) ArticleGetHandler(c echo.Context) error {
	uuid := c.Param("uuid")

	if err := googleUUID.Validate(uuid); err != nil {
		log.Printf("Error getting article, bad uuid: %v", uuid)
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "UUID is required"})
	}

	var article models.Article
	// Preload User and Media relations
	if err := h.DB.Preload("User").Preload("Media").Where("id = ?", uuid).First(&article).Error; err != nil {
		log.Printf("Error getting article with id: %v, 404", uuid)
		return c.JSON(http.StatusNotFound, echo.Map{"message": "No such article"})
	}

	// Prepare media responses with presigned URLs
	var mediaResponses []schemas.MediaCreateResponse
	for _, media := range article.Media {
		// Generate presigned URL for each media file
		_, presignedURL, err := utils.FullGeneratePresignedURL(h.MinIOClient, os.Getenv("MINIO_BUCKET"), 1*time.Hour)
		if err != nil {
			log.Printf("Error generating presigned URL for media: %v", err)
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Internal server error"})
		}
		mediaResponses = append(mediaResponses, schemas.MediaCreateResponse{
			FileName: media.FileName,
			S3Key:    presignedURL,
		})
	}

	resp := schemas.ArticleResponse{
		ID: article.ID.String(),
		Title:          article.Title,
		Content:        article.Content,
		MediaPresignedUrl: mediaResponses,
		AuthorUsername: article.User.Username,
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

	// Handle media update
	if articleData.Media != nil {
		// Delete old media records
		if err := h.DB.Where("article_id = ?", article.ID).Delete(&models.Media{}).Error; err != nil {
			log.Printf("Error deleting old media: %v", err)
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Internal server error"})
		}
		article.Media = nil

		// Add new media records and generate presigned upload URLs
		for _, mediaFileName := range *articleData.Media {
			name, uploadURL, err := utils.FullGeneratePresignedURL(h.MinIOClient, os.Getenv("MINIO_BUCKET"), 1*time.Hour)
			if err != nil {
				log.Printf("Error generating presigned URL: %v", err)
				return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Internal server error"})
			}
			media := models.Media{
				FileName:  mediaFileName,
				S3Key:     name,
				ArticleID: article.ID.String(),
			}
			if err := h.DB.Create(&media).Error; err != nil {
				log.Printf("Error creating media: %v", err)
				return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Internal server error"})
			}
			article.Media = append(article.Media, media)
			mediaResponses = append(mediaResponses, schemas.MediaCreateResponse{
				FileName: mediaFileName,
				S3Key:    uploadURL, // This is the presigned URL for upload
			})
		}
	} else {
		// If media is not updated, return presigned URLs for existing media
		for _, media := range article.Media {
			_, presignedURL, err := utils.FullGeneratePresignedURL(h.MinIOClient, os.Getenv("MINIO_BUCKET"), 1*time.Hour)
			if err != nil {
				log.Printf("Error generating presigned URL for media: %v", err)
				return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Internal server error"})
			}
			mediaResponses = append(mediaResponses, schemas.MediaCreateResponse{
				FileName: media.FileName,
				S3Key:    presignedURL,
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