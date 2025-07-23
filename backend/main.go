package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/tamago/todo-with-gemini/backend/internal/auth"
	"github.com/tamago/todo-with-gemini/backend/internal/db"
	"github.com/tamago/todo-with-gemini/backend/internal/logging"
	"github.com/tamago/todo-with-gemini/backend/internal/middleware"
	"github.com/tamago/todo-with-gemini/backend/internal/tasks"
	"github.com/tamago/todo-with-gemini/backend/internal/telemetry"
)

func main() {
	// Initialize structured logger
	logging.InitLogger()

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		logging.ContextLogger(context.Background()).Info("No .env file found, using environment variables")
	}

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

	router := gin.Default()

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

	// Public routes
	router.POST("/signup", auth.Signup(dbConn))
	router.POST("/login", auth.Login(dbConn))

	// Protected routes
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// Task routes
		protected.GET("/tasks", tasks.GetTasks(dbConn))
		protected.POST("/tasks", tasks.CreateTask(dbConn))
		protected.PUT("/tasks/:id", tasks.UpdateTask(dbConn))
		protected.DELETE("/tasks/:id", tasks.DeleteTask(dbConn))
	}

	logging.ContextLogger(context.Background()).Info("Server starting on port 8080")
	log.Fatal(router.Run(":8080"))
}
