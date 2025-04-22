package controllers

import (
	"errors"
	"log"

	_ "github.com/aotsurasak46/user-management/docs"
	"github.com/aotsurasak46/user-management/dto"
	"github.com/aotsurasak46/user-management/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with the provided information (Admin only)
// @Tags users
// @Accept json
// @Produce json
// @Param user body dto.UserCreateRequest true "User Information"
// @Success 201 {object} dto.UserResponse
// @Failure 400 {object} object{error=string} "Invalid request body, invalid role or email already exists"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /api/v1/users [post]
func CreateUser(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		inputUser := new(dto.UserCreateRequest)
		if err := c.BodyParser(&inputUser); err != nil {
			log.Printf("Error parsing request body: %v", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}
		if inputUser.Name == "" || inputUser.Email == "" || inputUser.Password == "" || inputUser.Role == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Input can't be empty"})
		}

		validRoles := map[string]bool{
			"user":  true,
			"admin": true,
		}
		if !validRoles[inputUser.Role] {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid role"})
		}

		var existingUser models.User
		if err := db.Unscoped().Where("email = ?", inputUser.Email).First(&existingUser).Error; err == nil {
			log.Printf("Duplicate email detected: %v", inputUser.Email)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email already exists"})
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(inputUser.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		user := &models.User{
			Name:     inputUser.Name,
			Email:    inputUser.Email,
			Password: string(hashedPassword),
			Role:     inputUser.Role,
		}

		if err := db.Create(user).Error; err != nil {
			log.Printf("Error creating user: %v", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.Status(fiber.StatusCreated).JSON(user)
	}
}

// GetUsers godoc
// @Summary Get all users
// @Description Retrieve a list of all users from the database
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} dto.UserResponse
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /api/v1/users [get]
func GetUsers(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		users := new([]models.User)
		result := db.Find(&users)
		if result.Error != nil {
			log.Printf("Error getting users from database: %v", result.Error)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.JSON(users)
	}
}

// GetUserById godoc
// @Summary Get user by id
// @Description Retrieve a user information from the database by using id
// @Tags users
// @Accept json
// @Produce json
// @param id path int true "User id"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} object{error=string} "Bad request or User ID is missing"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /api/v1/users/:id [get]
func GetUserById(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := c.Params("id")
		if userId == "" {
			log.Printf("User ID is missing in the request")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User ID is required"})
		}
		user := new(models.User)
		result := db.Where("id = ?", userId).First(&user)
		if result.Error != nil {
			log.Printf("Error finding user in database: %v", result.Error)
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
			}
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.JSON(user)
	}
}

// UpdateUser godoc
// @Summary Update User by id
// @Description Update a user information by using id (Admin only)
// @Tags users
// @Accept json
// @Produce json
// @param id path int true "User id"
// @Param user body  dto.UserUpdateRequest true "User Information"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} object{error=string} "Bad request, invalid request body, invalid role, email is existed or invalid role"
// @Failure 404 {object} object{error=string} "User not found"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /api/v1/users/:id [put]
func UpdateUser(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := c.Params("id")
		if userId == "" {
			log.Printf("User ID is missing in the request")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User ID is required"})
		}

		input := new(dto.UserUpdateRequest)
		if err := c.BodyParser(&input); err != nil {
			log.Printf("Error parsing request body: %v", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		var user models.User
		if err := db.First(&user, userId).Error; err != nil {
			log.Printf("Error finding user in database: %v", err)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
			}
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if input.Email != "" && input.Email != user.Email {
			var existing models.User
			if err := db.Unscoped().Where("email = ?", input.Email).First(&existing).Error; err == nil {
				log.Printf("Duplicate email detected: %v", input.Email)
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email already exists"})
			}
			user.Email = input.Email
		}

		if input.Name != "" {
			user.Name = input.Name
		}

		if input.Role != "" {
			validRoles := map[string]bool{
				"user":  true,
				"admin": true,
			}
			if !validRoles[input.Role] {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid role"})
			}
			user.Role = input.Role
		}

		if err := db.Save(&user).Error; err != nil {
			log.Printf("Error updating user in database: %v", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.JSON(user)
	}
}

// DeleteUser godoc
// @Summary Delete User by id (Admin only)
// @Description Delete a user from the database by using their id
// @Tags users
// @Accept json
// @Produce json
// @param id path int true "User id"
// @Success 200 {object} object{message=string}
// @Failure 400 {object} object{error=string} "Bad request or User ID is missing"
// @Failure 404 {object} object{error=string} "User not found"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /api/v1/users/:id [delete]
func DeleteUser(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := c.Params("id")
		if userId == "" {
			log.Printf("User ID is missing in the request")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User ID is required"})
		}
		if err := db.Delete(&models.User{}, userId).Error; err != nil {
			log.Printf("Error deleting user from database: %v", err)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
			}
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.JSON(fiber.Map{"message": "User deleted successfully"})
	}
}
