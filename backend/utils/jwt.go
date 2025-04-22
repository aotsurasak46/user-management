package utils

import (
	"fmt"
	"os"
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func GenerateJWT(userID uint) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	return token.SignedString([]byte(jwtSecretKey))
}

func ParseJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrInvalidKey
	}
	return claims, nil
}

func GetUserIDFromJWT(c *fiber.Ctx) (uint, error) {
    tokenString := c.Cookies("jwt")
    if tokenString == "" {
        return 0, fiber.ErrUnauthorized
    }

    claims,err := ParseJWT(tokenString)
	if err != nil {
		return 0 ,err
	}
	userIDValue, ok := claims["user_id"]
	if !ok {
		return 0, fmt.Errorf("user_id not found in claims")
	}
	userIDFloat, ok := userIDValue.(float64)
	if !ok {
		return 0, fmt.Errorf("user_id is not a valid number")
	}
	return uint(userIDFloat), nil
}
