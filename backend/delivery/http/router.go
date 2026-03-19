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

	r.Use(middleware.CORSMiddleware())

	r.Static("/css", "template/css")
	r.Static("/js", "template/js")

	r.LoadHTMLGlob("template/*.html")

	// Apply frontend redirection middleware
	frontendRoutes := r.Group("/")
	frontendRoutes.Use(middleware.AuthRedirectMiddleware(jwtSecret))
	{
		frontendRoutes.GET("/", func(c *gin.Context) { c.HTML(200, "index.html", gin.H{}) })
		frontendRoutes.GET("/index.html", func(c *gin.Context) { c.HTML(200, "index.html", gin.H{}) })
		frontendRoutes.GET("/login.html", func(c *gin.Context) { c.HTML(200, "login.html", gin.H{}) })
		frontendRoutes.GET("/signup.html", func(c *gin.Context) { c.HTML(200, "signup.html", gin.H{}) })
		frontendRoutes.GET("/otp.html", func(c *gin.Context) { c.HTML(200, "otp.html", gin.H{}) })
		frontendRoutes.GET("/dashboard.html", func(c *gin.Context) { c.HTML(200, "dashboard.html", gin.H{}) })
		frontendRoutes.GET("/pomodoro.html", func(c *gin.Context) { c.HTML(200, "pomodoro.html", gin.H{}) })
	}

	api := r.Group("/api")

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
