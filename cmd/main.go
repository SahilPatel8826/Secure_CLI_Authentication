package main

import (
	"fmt"
	"log"
	"os"

	"cli/internal/cli"
	"cli/internal/database"
	"cli/internal/models"
	"cli/internal/repository"
	service "cli/internal/services"

	"github.com/joho/godotenv"
)

// ==========================================================
// Main Entry Point
//
// Initializes the application by:
//   - Loading environment variables
//   - Establishing database connection
//   - Running database migrations
//   - Creating repositories and services
//   - Starting the interactive CLI
// ==========================================================

func main() {

	// Load environment variables from .env file.
	// Falls back to system environment variables if not found.
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, using environment variables")
	}

	// Build PostgreSQL Data Source Name (DSN).
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	// Establish database connection.
	db, err := database.Connect(dsn)
	if err != nil {
		log.Fatal(err)
	}

	// Automatically create/update database tables.
	err = db.AutoMigrate(
		&models.User{},
		&models.Session{},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(`
==========================================
        AUTH CLI SYSTEM v1.0
==========================================

Secure Authentication Platform

Type 'help' to view available commands.
`)

	// Initialize repositories.
	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)

	// Initialize authentication service.
	authService := service.NewAuthService(
		userRepo,
		sessionRepo,
	)

	// Start the interactive CLI.
	cli := cli.NewCLI(authService)
	cli.Run()
}
