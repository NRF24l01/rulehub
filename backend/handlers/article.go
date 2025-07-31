package handlers

import (
	"log"
	"net/http"
	"os"
	"rulehub/models"
	"rulehub/schemas"
	"rulehub/utils"

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
	for _, mediaFileName := range article_data.Media {
		// Генерируем presigned PUT URL и UUID ключ для загрузки медиа
		name, uploadURL, err := utils.GeneratePresignedPutURL(h.MinIOClient, os.Getenv("MINIO_BUCKET"), utils.GetPresignedLifetime())
		if err != nil {
			log.Printf("Error generating presigned PUT URL: %v", err)
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

		mediaResponses = append(mediaResponses, schemas.MediaCreateResponse{
			FileName: mediaFileName,
			S3Key:    uploadURL, // ссылка для загрузки
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
	bucket := os.Getenv("MINIO_BUCKET")
	endpoint := os.Getenv("MINIO_PUBLIC_ENDPOINT") // например, https://minio.example.com

	for _, media := range article.Media {
		// Формируем постоянную публичную ссылку на файл
		publicURL := endpoint + "/" + bucket + "/" + media.S3Key
		mediaResponses = append(mediaResponses, schemas.MediaCreateResponse{
			FileName: media.FileName,
			S3Key:    publicURL, // постоянная ссылка для скачивания
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
		// Удаляем старые медиа
		if err := h.DB.Where("article_id = ?", article.ID).Delete(&models.Media{}).Error; err != nil {
			log.Printf("Error deleting old media: %v", err)
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Internal server error"})
		}
		article.Media = nil

		// Добавляем новые медиа с presigned PUT URL для загрузки
		for _, mediaFileName := range *articleData.Media {
			name, uploadURL, err := utils.GeneratePresignedPutURL(h.MinIOClient, os.Getenv("MINIO_BUCKET"), utils.GetPresignedLifetime())
			if err != nil {
				log.Printf("Error generating presigned PUT URL: %v", err)
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
				S3Key:    uploadURL, // ссылка для загрузки
			})
		}
	} else {
		// Если media не обновлялась, вернуть presigned GET URL для существующих
		for _, media := range article.Media {
			presignedURL, err := utils.GeneratePresignedGetURL(h.MinIOClient, os.Getenv("MINIO_BUCKET"), media.S3Key, utils.GetPresignedLifetime())
			if err != nil {
				log.Printf("Error generating presigned GET URL for media, ignore: %v", err)
				presignedURL = "" 
			}
			mediaResponses = append(mediaResponses, schemas.MediaCreateResponse{
				FileName: media.FileName,
				S3Key:    presignedURL, // ссылка для скачивания
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
