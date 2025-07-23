package auth

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/tamago/todo-with-gemini/backend/internal/models"
	"github.com/tamago/todo-with-gemini/backend/internal/utils"
)

// Login handles user login and returns a JWT token.
func Login(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Retrieve the user from the database
		var storedPasswordHash string
		var userID int
		query := "SELECT id, password_hash FROM users WHERE username = $1"
		err := db.QueryRow(query, user.Username).Scan(&userID, &storedPasswordHash)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		// Compare the provided password with the stored hashed password
		if err := bcrypt.CompareHashAndPassword([]byte(storedPasswordHash), []byte(user.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		// Generate JWT token
		token, err := utils.GenerateToken(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}
