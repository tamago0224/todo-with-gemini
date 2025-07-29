package main

import (
	"context"
	"log/slog"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/tamago/todo-with-gemini/backend/internal/platform/logging"
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
	os.Setenv("SERVICE_NAME", "api-gateway")

	// Initialize OpenTelemetry
	shutdownTracer := telemetry.InitTracer()
	defer func() {
		if err := shutdownTracer(context.Background()); err != nil {
			logging.ContextLogger(context.Background()).Error("failed to shutdown TracerProvider", "error", err)
		}
	}()

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

	// Define target URLs for services
	authServiceURL, _ := url.Parse("http://auth-service:8081") // Auth Service
	taskServiceURL, _ := url.Parse("http://task-service:8080") // Task Service

	// Create reverse proxies
	authProxy := httputil.NewSingleHostReverseProxy(authServiceURL)
	taskProxy := httputil.NewSingleHostReverseProxy(taskServiceURL)

	// Auth Service routes
	router.POST("/signup", func(c *gin.Context) { authProxy.ServeHTTP(c.Writer, c.Request) })
	router.POST("/login", func(c *gin.Context) { authProxy.ServeHTTP(c.Writer, c.Request) })

	// Task Service routes
	router.Any("/api/tasks/*any", func(c *gin.Context) { taskProxy.ServeHTTP(c.Writer, c.Request) })

	logging.ContextLogger(context.Background()).Info("API Gateway starting on port 8080")
	if err := router.Run("0.0.0.0:8080"); err != nil {
		slog.Error("Failed to run API Gateway", "error", err)
		os.Exit(1)
	}
}
