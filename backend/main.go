package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aotsurasak46/user-management/controllers"
	_ "github.com/aotsurasak46/user-management/docs"
	"github.com/aotsurasak46/user-management/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
)

// @title User Management API
// @description This is a server for managing users, authentication, and chat functionality.
// @version 1.0
// @host localhost:8080
// @schemes http
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// This project is a backend server for a user management system. It provides the following features:
// - User authentication and authorization
// - CRUD operations for user management
// - Chat functionality with WebSocket support
// - API documentation using Swagger
// - Middleware for CORS, authentication, and admin-only access
// - Graceful shutdown handling
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	err = ConnectDB()
	if err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	}

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173,http://localhost",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,PUT,DELETE",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
	}))

	app.Get("/swagger/*", swagger.HandlerDefault)

	app.Use("/api/v1/users", middleware.Authen(DB))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World")
	})

	app.Get("/ws/chat", middleware.WebSocketUpgradeAuth(DB), controllers.ChatSocketHandler(DB))
	app.Get("/api/v1/messages/:userId", middleware.Authen(DB), controllers.GetChatHistory(DB))
	app.Get("/api/v1/conversations", middleware.Authen(DB), controllers.GetConversations((DB)))

	app.Post("/api/v1/login", controllers.LoginUser(DB))
	app.Post("/api/v1/logout", controllers.LogoutUser())
	app.Post("/api/v1/register", controllers.RegisterUser(DB))
	app.Get("/api/v1/check-auth", middleware.Authen(DB), controllers.CheckAuth(DB))

	app.Get("/api/v1/users", controllers.GetUsers(DB))
	app.Get("/api/v1/users/:id", controllers.GetUserById(DB))
	app.Post("/api/v1/users", middleware.AdminOnly(DB), controllers.CreateUser(DB))
	app.Put("/api/v1/users/:id", middleware.AdminOnly(DB), controllers.UpdateUser(DB))
	app.Delete("/api/v1/users/:id", middleware.AdminOnly(DB), controllers.DeleteUser(DB))

	idleConnsClosed := make(chan struct{})
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		log.Println("Shutting down gracefully...")

		_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := app.Shutdown(); err != nil {
			log.Printf("Error shutting down server: %v", err)
		}
		if err := CloseDB(); err != nil {
			log.Printf("Error closing db connection: %v", err)
		}

		log.Println("Server shutdown complete")
		close(idleConnsClosed)
	}()

	if err := app.Listen(":8080"); err != nil {
		log.Printf("Server error: %v", err)
		os.Exit(1)
	}

	<-idleConnsClosed
}
