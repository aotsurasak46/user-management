package controllers

import (
	"errors"
	"log"
	"time"

	_ "github.com/aotsurasak46/user-management/docs"
	"github.com/aotsurasak46/user-management/dto"
	"github.com/aotsurasak46/user-management/models"
	"github.com/aotsurasak46/user-management/utils"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// RegisterUser godoc
// @Summary User Register
// @Description Create a new user with name, email and password
// @Tags authentication
// @Accept json
// @Produce json
// @Param user body dto.RegisterRequest true "User Information"
// @Success 201 {object} dto.UserResponse
// @Failure 400 {object} object{error=string} "Invalid request body"
// @Failure 401 {object} object{error=string} "Email already exists"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /api/v1/register [post]
func RegisterUser(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		inputUser := new(dto.RegisterRequest)
		if err := c.BodyParser(&inputUser); err != nil {
			log.Printf("Error parsing request body: %v", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(inputUser.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if inputUser.Name == "" || inputUser.Email == "" || inputUser.Password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Input can't be empty"})
		}

		user := new(models.User)
		user.Name = inputUser.Name
		user.Email = inputUser.Email
		user.Password = string(hashedPassword)

		if err := db.Create(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) || db.Where("email = ?", user.Email).First(&models.User{}).Error == nil {
				log.Printf("Duplicate email detected: %v", user.Email)
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Email already exists"})
			}
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.Status(fiber.StatusCreated).JSON(user)
	}
}

// LoginUser godoc
// @Summary User login
// @Description Authenticate a user with email and password
// @Tags authentication
// @Accept json
// @Produce json
// @Param credentials body dto.LoginRequest true "Login Credentials"
// @Success 200 {object} object{message=string,user=dto.UserResponse} "Login successful"
// @Failure 400 {object} object{error=string} "Bad request"
// @Failure 401 {object} object{error=string} "Invalid email or password"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /api/v1/login [post]
func LoginUser(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		inputUser := new(dto.LoginRequest)
		dbUser := new(models.User)
		if err := c.BodyParser(&inputUser); err != nil {
			log.Printf("Error parsing request body: %v", err)
			return c.SendStatus(fiber.StatusBadRequest)
		}
		if inputUser.Email == "" || inputUser.Password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email or password can't be empty"})
		}
		result := db.Where("email = ?", inputUser.Email).First(&dbUser)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				log.Printf("User not found: %v", inputUser.Email)
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
			}
			log.Printf("Error finding user in database: %v", result.Error)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(inputUser.Password)); err != nil {
			log.Printf("Error comparing password: %v", err)
			if err == bcrypt.ErrMismatchedHashAndPassword {
				log.Printf("Password mismatch: %v", err)
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
			}
			log.Printf("Error comparing hash password: %v", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		tokenString, err := utils.GenerateJWT(dbUser.ID)
		if err != nil {
			log.Printf("Error generate JWT: %v", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		c.Cookie(&fiber.Cookie{
			Name:     "jwt",
			Value:    tokenString,
			Expires:  time.Now().Add(time.Hour * 72),
			HTTPOnly: true,
			SameSite: "Lax",
			Path:     "/",
			Secure:   false,
		})

		return c.JSON(dto.UserResponse{
			ID:    dbUser.ID,
			Name:  dbUser.Name,
			Email: dbUser.Email,
			Role:  dbUser.Role,
		})
	}
}

// CheckAuth godoc
// @Summary Check Authentication
// @Description Verify if the user is authenticated and retrieve user details
// @Tags authentication
// @Produce json
// @Success 200 {object} object{authenticated=bool,user=object{id=uint,name=string,email=string,role=string}}
// @Failure 401 {object} object{error=string} "User not found or unauthorized"
// @Router /api/v1/check-auth [get]
func CheckAuth(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("userID").(uint)
		var user models.User
		if err := db.First(&user, userID).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not found"})
		}
		return c.JSON(fiber.Map{
			"authenticated": true,
			"user": fiber.Map{
				"id":    user.ID,
				"name":  user.Name,
				"email": user.Email,
				"role":  user.Role,
			},
		})
	}
}

// LogoutUser godoc
// @Summary User logout
// @Description Logout a user by clearing the JWT cookie
// @Tags authentication
// @Success 200 {object} object{message=string} "Logout successful"
// @Router /api/v1/logout [post]
func LogoutUser() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Cookie(&fiber.Cookie{
			Name:     "jwt",
			Value:    "",
			Expires:  time.Now().Add(-time.Hour),
			HTTPOnly: true,
		})
		return c.JSON(fiber.Map{"message": "Logout successful"})
	}
}
