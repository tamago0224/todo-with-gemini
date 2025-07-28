package controllers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"golang.org/x/crypto/bcrypt"

	"github.com/tamago/todo-with-gemini/backend/internal/models"
	"github.com/tamago/todo-with-gemini/backend/internal/utils"
	"fmt"
)

type AuthController struct {
	db *sql.DB
}

func NewAuthController(db *sql.DB) *AuthController {
	return &AuthController{db: db}
}

// Login handles user login and returns a JWT token.
func (ac *AuthController) Login(c *gin.Context) {
	_, span := otel.Tracer("").Start(c.Request.Context(), "AuthController.Login")
	defer span.End()

	fmt.Println("Login function called")

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Debugging: Print received username and password
	fmt.Printf("Received username: %s, password: %s\n", user.Username, user.Password)

	// Retrieve the user from the database
	var storedPasswordHash string
	var userID int
	query := "SELECT id, password_hash FROM users WHERE username = $1"
	err := ac.db.QueryRow(query, user.Username).Scan(&userID, &storedPasswordHash)
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

// Signup handles user registration.
func (ac *AuthController) Signup(c *gin.Context) {
	_, span := otel.Tracer("").Start(c.Request.Context(), "AuthController.Signup")
	defer span.End()

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		query := "INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id"
		var id int
		err = ac.db.QueryRow(query, user.Username, string(hashedPassword)).Scan(&id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "id": id})
}
