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

func (h* Handler) UserRegistrationHandler(c echo.Context) error {
	user_data := c.Get("validatedBody").(*schemas.SignUpRequest)
	log.Printf("Registering user: %+v", user_data)

	// Check if the username already exists
	var existingUser models.User
	if err := h.DB.Where("username = ?", user_data.Username).First(&existingUser).Error; err == nil {
		log.Printf("Username already exists: %s", user_data.Username)
		return c.JSON(http.StatusConflict, echo.Map{"message": "Username already exists"})
	}

	// Create a new user
	newUser := models.User{
		Username: user_data.Username,
		Password: utils.HashPassword(user_data.Password),
	}

	if err := h.DB.Create(&newUser).Error; err != nil {
		log.Printf("Error creating user: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Internal server error"})
	}

	return c.JSON(http.StatusCreated,  schemas.SignUpResponse{
		ID:       newUser.ID.String(),
		Username: newUser.Username,
	})
}

func (h* Handler) UserJwtRefreshHandler(c echo.Context) error {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		log.Printf("No refresh token found: %v", err)
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Refresh token not found"})
	}

	jwtClaims, err := utils.ValidateToken(refreshToken.Value, []byte(os.Getenv("PASSWORD_JWT_REFRESH_SECRET")))
	if err != nil {
		log.Printf("Invalid refresh token: %v", err)
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Invalid refresh token"})
	}

	userID, ok := jwtClaims["user_id"].(string)
	if !ok {
		log.Println("Invalid user ID in token claims")
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Invalid token claims"})
	}

	// Check if the user exists in the database
	var user models.User
	if err := h.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		log.Printf("User not found: %v", err)
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "User not found"})
	}
	// Generate a new access token
	accessToken, err := utils.GenerateAccessToken(user.ID.String(), user.Username)
	if err != nil {
		log.Printf("Error generating access token: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Internal server error"})
	}
	
	return c.JSON(http.StatusOK, schemas.SignInResponse{
		AccessJWT: accessToken,
	})
}