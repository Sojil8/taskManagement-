package main

import (
	"database/sql"
	"fmt"
	"os"

	"taskmanager/delivery/http"
	"taskmanager/delivery/http/handler"
	"taskmanager/infrastructure/email"
	"taskmanager/infrastructure/logger"
	"taskmanager/infrastructure/postgres"
	"taskmanager/migrations"
	"taskmanager/usecase"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	// Initialize global Zap logger
	logger.InitLogger()
	defer logger.Sync()

	if err := godotenv.Load(); err != nil {
		logger.Log.Info("No .env file found, relying on environment variables")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		logger.Log.Fatal("DB_URL environment variable is required")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		logger.Log.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		logger.Log.Fatal("Database unreachable", zap.Error(err))
	}
	logger.Log.Info("Connected to PostgreSQL successfully")

	// Run Database Migrations Automatically
	logger.Log.Info("Checking and running database migrations...")
	if err := migrations.RunAutoMigrations(db, "migrations"); err != nil {
		logger.Log.Fatal("Database migrations failed", zap.Error(err))
	}

	// Email Config
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default_secret_key"
		logger.Log.Warn("JWT_SECRET not set, using default")
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

	logger.Log.Info("Server starting", zap.String("port", port))
	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		logger.Log.Fatal("Server failed", zap.Error(err))
	}
}
