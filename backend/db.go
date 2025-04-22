package main

import (
	"os"
	"fmt"
    "gorm.io/gorm"
	"gorm.io/driver/postgres"
	"github.com/aotsurasak46/user-management/models"
)

var DB *gorm.DB
 
func ConnectDB() error {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

    dsn := fmt.Sprintf("host=%s port=%s user=%s "+
	"password=%s dbname=%s sslmode=disable",
	dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	err = db.AutoMigrate(&models.User{} ,&models.Message{})
	if err != nil { 
		return fmt.Errorf("failed to auto migrate database: %w", err)
	}
	fmt.Println("Database migration completed!")
	
    DB = db
	return nil
}


func CloseDB() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}
	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}
	fmt.Println("Database connection closed successfully!")
	return nil
}