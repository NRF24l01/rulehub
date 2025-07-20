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

	resp := schemas.ArticleCreateResponse{
		ID:             article.ID.String(),
		Title:          article.Title,
		Content:        article.Content,
		MediaPresignedUrl: mediaResponses,
		AuthorUsername: user.Username,
	}
	return c.JSON(http.StatusCreated, resp)
}