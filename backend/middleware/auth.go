package middleware

import (
	"fmt"
	"os"

	"github.com/aotsurasak46/user-management/models"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

func Authen(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		cookie := c.Cookies("jwt")
		if cookie == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
		}
		jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
		token, err := jwt.ParseWithClaims(cookie, &jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(jwtSecretKey), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
		}

		claims := token.Claims.(*jwt.MapClaims)
		userID := uint((*claims)["user_id"].(float64))

		var user models.User
		if err := db.First(&user, userID).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not found"})
		}
		c.Locals("userID", uint(userID))
		return c.Next()
	}
}

func AdminOnly(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		userID := c.Locals("userID").(uint)
		var user models.User
		if err := db.First(&user, userID).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not found"})
		}

		if user.Role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
		}
		return c.Next()
	}
}

func WebSocketUpgradeAuth(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !websocket.IsWebSocketUpgrade(c) {
			return fiber.ErrUpgradeRequired
		}
		return Authen(db)(c)
	}
}
