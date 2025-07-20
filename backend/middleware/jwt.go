package middleware

import (
	"net/http"
	"os"

	"quoter_back/utils"

	"github.com/labstack/echo/v4"
)

// JWTMiddleware создает middleware для проверки JWT токена
func JWTMiddleware() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            // Извлекаем токен из заголовка Authorization
            authHeader := c.Request().Header.Get("Authorization")
            if authHeader == "" {
                return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing token"})
            }

            // Убираем "Bearer " из заголовка
            tokenString := ""
            if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
                tokenString = authHeader[7:]
            } else {
                return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token format"})
            }

            // Получаем секретный ключ через переданную функцию
            secret := []byte(os.Getenv("PASSWORD_JWT_ACCESS_SECRET"))

            // Проверяем токен с помощью функции ValidateToken
            claims, err := utils.ValidateToken(tokenString, secret)
            if err != nil {
                return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid or expired token"})
            }

            // Извлекаем user_id из claims
            userID, ok := claims["user_id"].(string)
            if !ok {
                return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token claims"})
            }

            // Передаем user_id в контекст
            c.Set("user_id", userID)

            // Переходим к следующему обработчику
            return next(c)
        }
    }
}