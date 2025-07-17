package handlers

import (
	"net/http"

	"log"

	"rulehub/models"
	"rulehub/schemas"
	"rulehub/utils"

	"github.com/labstack/echo/v4"
)

func (h* Handler) UserLoginHandler(c echo.Context) error {
	user_data := c.Get("validatedBody").(*schemas.SignInRequest)
	log.Printf("Logining user: %+v", user_data)

	// Check if the user exists in the database and validate the password
	var user models.User
	if err := h.DB.Where("username = ?", user_data.Username).First(&user).Error; err != nil {
		log.Printf("User not found: %v", err)
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Invalid username or password"})
	}

	if utils.CheckPassword(user_data.Password, user.Password) == false {
		log.Printf("Invalid password for user: %s", user_data.Username)
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Invalid username or password"})
	}

	// Generate access and refresh tokens
	accessToken, err := utils.GenerateAccessToken(user.ID.String(), user.Username)
	if err != nil {
		log.Printf("Error generating access token: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Internal server error"})
	}
	refreshToken, err := utils.GenerateRefreshToken(user.ID.String())
	if err != nil {
		log.Printf("Error generating refresh token: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Internal server error"})
	}

	// Set the refresh token in an HttpOnly cookie
	cookie := new(http.Cookie)
	cookie.Name = "refresh_token"
	cookie.Value = refreshToken
	cookie.HttpOnly = true
	cookie.Path = "/"
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, schemas.SignInResponse{
		AccessJWT: accessToken,
	})
}