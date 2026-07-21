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

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, using environment variables")
	}
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := database.Connect(dsn)

	if err != nil {
		log.Fatal(err)
	}
	er := db.AutoMigrate(&models.User{}, &models.Session{})
	if er != nil {
		panic(er)

	}
	fmt.Println("✅ Connected successfully!")
	userRepo := repository.NewUserRepository(db)

	sessionRepo := repository.NewSessionRepository(db)

	authService := service.NewAuthService(
		userRepo,
		sessionRepo,
	)

	cli := cli.NewCLI(authService)

	cli.Run()

}
