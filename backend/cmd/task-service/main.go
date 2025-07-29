package main

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/tamago/todo-with-gemini/backend/internal/app/controllers"
	"github.com/tamago/todo-with-gemini/backend/internal/app/repositories"
	"github.com/tamago/todo-with-gemini/backend/internal/app/services"
	"github.com/tamago/todo-with-gemini/backend/internal/platform/db"
	"github.com/tamago/todo-with-gemini/backend/internal/platform/logging"
	"github.com/tamago/todo-with-gemini/backend/internal/platform/middleware"
	"github.com/tamago/todo-with-gemini/backend/internal/platform/telemetry"
)

func main() {
	// Initialize structured logger
	logging.InitLogger()

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		logging.ContextLogger(context.Background()).Info("No .env file found, using environment variables")
	}

	// Set SERVICE_NAME for OpenTelemetry
	os.Setenv("SERVICE_NAME", "task-service")

	// Initialize OpenTelemetry
	shutdownTracer := telemetry.InitTracer()
	defer func() {
		if err := shutdownTracer(context.Background()); err != nil {
			logging.ContextLogger(context.Background()).Error("failed to shutdown TracerProvider", "error", err)
		}
	}()

	// Initialize database connection
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		logging.ContextLogger(context.Background()).Error("DATABASE_URL environment variable not set")
		os.Exit(1)
	}
	dbConn, err := db.InitDB(dbURL)
	if err != nil {
		logging.ContextLogger(context.Background()).Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer func(dbConn *sql.DB) {
		err := dbConn.Close()
		if err != nil {
			logging.ContextLogger(context.Background()).Error("Error closing database connection", "error", err)
		}
	}(dbConn)

	router := gin.New()
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Apply CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Apply OpenTelemetry Gin middleware
	router.Use(telemetry.GinMiddleware())

	// Initialize the layers
	taskRepo := repositories.NewPostgresTaskRepository(dbConn)
	taskService := services.NewTaskService(taskRepo)
	taskController := controllers.NewTaskController(taskService)

	// Protected routes
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// Task routes
		protected.GET("/tasks", taskController.GetTasks)
		protected.POST("/tasks", taskController.CreateTask)
		protected.PUT("/tasks/:id", taskController.UpdateTask)
		protected.DELETE("/tasks/:id", taskController.DeleteTask)
	}

	logging.ContextLogger(context.Background()).Info("Task Service starting on port 8080")
	if err := router.Run(":8080"); err != nil {
		slog.Error("Failed to run task router", "error", err)
		os.Exit(1)
	}
}
