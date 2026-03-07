package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"taskmanager/delivery/http"
	"taskmanager/delivery/http/handler"
	"taskmanager/infrastructure/email"
	"taskmanager/infrastructure/postgres"
	"taskmanager/usecase"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL environment variable is required")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Database unreachable: %v", err)
	}
	log.Println("Connected to PostgreSQL successfully")

	// Email Config
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default_secret_key"
		log.Println("Warning: JWT_SECRET not set, using default")
	}

	// Init Repositories
	userRepo := postgres.NewUserRepository(db)
	otpRepo := postgres.NewOTPRepository(db)
	taskRepo := postgres.NewTaskRepository(db)
	categoryRepo := postgres.NewCategoryRepository(db)

	// Init Services
	emailSvc := email.NewSMTPService(smtpHost, smtpPort, smtpUser, smtpPass)

	// Init UseCases
	authUseCase := usecase.NewAuthUseCase(userRepo, otpRepo, emailSvc, jwtSecret)
	taskUseCase := usecase.NewTaskUseCase(taskRepo)
	categoryUseCase := usecase.NewCategoryUseCase(categoryRepo)

	// Init Handlers
	authHandler := handler.NewAuthHandler(authUseCase)
	taskHandler := handler.NewTaskHandler(taskUseCase)
	categoryHandler := handler.NewCategoryHandler(categoryUseCase)

	// Setup Router
	r := http.SetupRouter(authHandler, taskHandler, categoryHandler, jwtSecret)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
