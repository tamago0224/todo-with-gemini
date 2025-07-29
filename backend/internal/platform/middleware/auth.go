package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/tamago/todo-with-gemini/backend/internal/platform/utils"
)

// AuthMiddleware validates the JWT token from the request header.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Remove "Bearer " prefix
		tokenString = tokenString[len("Bearer "):]

		userID, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		slog.Info("Authenticated user", "userID", userID)
		c.Set("userID", userID)
		c.Next()
	}
}
