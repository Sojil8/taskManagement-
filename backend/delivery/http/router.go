package http

import (
	"time"

	"taskmanager/delivery/http/handler"
	"taskmanager/delivery/http/middleware"
	"taskmanager/infrastructure/logger"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
)

func SetupRouter(authHandler *handler.AuthHandler, taskHandler *handler.TaskHandler, categoryHandler *handler.CategoryHandler, jwtSecret string) *gin.Engine {
	r := gin.New()

	// Add Zap Logger middleware
	r.Use(ginzap.Ginzap(logger.Log, time.RFC3339, true))

	// Add Logs Recovery middleware to handle panics safely
	r.Use(ginzap.RecoveryWithZap(logger.Log, true))

	// Add Custom Middleware
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.SecurityHeadersMiddleware())

	api := r.Group("/api/v1")

	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/send-otp", authHandler.SendOTP)
		auth.POST("/verify-otp", authHandler.VerifyOTP)
		auth.POST("/login", authHandler.Login)
	}

	tasks := api.Group("/tasks")
	tasks.Use(middleware.JWTAuthMiddleware(jwtSecret))
	{
		tasks.GET("", taskHandler.GetTasks)
		tasks.POST("", taskHandler.CreateTask)
		tasks.PUT("/:id", taskHandler.UpdateTask)
		tasks.PATCH("/:id/checkpoints/:cpId", taskHandler.CompleteCheckpoint)
		tasks.DELETE("/:id", taskHandler.DeleteTask)
	}

	categories := api.Group("/categories")
	categories.Use(middleware.JWTAuthMiddleware(jwtSecret))
	{
		categories.GET("", categoryHandler.GetCategories)
		categories.POST("", categoryHandler.CreateCategory)
	}

	return r
}
