package handlers

import (
	"net/http"
	"log"
	"os"
	"rulehub/models"
	"rulehub/schemas"
	"rulehub/utils"
	"github.com/labstack/echo/v4"
)

func (h* Handler) ArticleCreateHandler(c echo.Context) error {
	article_data := c.Get("validatedBody").(*schemas.ArticleCreateRequest)
	log.Printf("Creating article: %+v", article_data)

	// Create a new article in the database
	article := models.Article{
		Title:   article_data.Title,
		Content: article_data.Content,
		UserID:  c.Get("userID").(string),
	}

	// Create media entries if any
	for _, mediaFileName := range article_data.Media {
		media := models.Media{
			FileName: mediaFileName,
			S3Key:    utils.GenerateS3Key(mediaFileName),
			ArticleID: article.ID.String(),
			Article:   article,
		}

	if err := h.DB.Create(&article).Error; err != nil {
		log.Printf("Error creating article: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Internal server error"})
	}

	return c.JSON(http.StatusCreated, schemas.ArticleCreateResponse{
		ID:    article.ID.String(),
		Title: article.Title,
	})
}