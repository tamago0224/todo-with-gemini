package main

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	otelhttp "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

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

	// Create HTTP client with OpenTelemetry transport
	client := &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	// Auth Service base URL
	authServiceBaseURL := "http://auth-service:8081"

	// Task Service base URL
	taskServiceBaseURL := "http://task-service:8080"

	// Helper function to forward requests
	forwardRequest := func(c *gin.Context, targetBaseURL string) {
		// Construct the new URL
		targetURL, err := url.Parse(targetBaseURL + c.Request.URL.Path)
		if err != nil {
			slog.Error("Failed to parse target URL", "error", err, "targetBaseURL", targetBaseURL, "path", c.Request.URL.Path)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		// Create a new request with the same method, body, and headers
		reqBody, err := io.ReadAll(c.Request.Body)
		if err != nil {
			slog.Error("Failed to read request body", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		defer c.Request.Body.Close()

		proxyReq, err := http.NewRequestWithContext(c.Request.Context(), c.Request.Method, targetURL.String(), bytes.NewBuffer(reqBody))
		if err != nil {
			slog.Error("Failed to create proxy request", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		// Copy headers from original request
		for name, values := range c.Request.Header {
			for _, value := range values {
				proxyReq.Header.Add(name, value)
			}
		}

		// Send the request
		resp, err := client.Do(proxyReq)
		if err != nil {
			slog.Error("Failed to send proxy request", "error", err, "targetURL", targetURL.String())
			c.JSON(http.StatusBadGateway, gin.H{"error": "Bad Gateway"})
			return
		}
		defer resp.Body.Close()

		// Copy response headers and status code
		for name, values := range resp.Header {
			for _, value := range values {
				c.Writer.Header().Add(name, value)
			}
		}
		c.Writer.WriteHeader(resp.StatusCode)

		// Copy response body
		if _, err := io.Copy(c.Writer, resp.Body); err != nil {
			slog.Error("Failed to copy response body", "error", err)
			return
		}
	}

	// Auth Service routes
	router.POST("/signup", func(c *gin.Context) { forwardRequest(c, authServiceBaseURL) })
	router.POST("/login", func(c *gin.Context) { forwardRequest(c, authServiceBaseURL) })

	// Task Service routes
	router.Any("/api/tasks/*any", func(c *gin.Context) {
		// Rewrite path for task service
		c.Request.URL.Path = "/api/tasks" + c.Request.URL.Path[len("/api/tasks"):]
		forwardRequest(c, taskServiceBaseURL)
	})

	logging.ContextLogger(context.Background()).Info("API Gateway starting on port 8080")
	if err := router.Run("0.0.0.0:8080"); err != nil {
		slog.Error("Failed to run API Gateway", "error", err)
		os.Exit(1)
	}
}