package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"

	"github.com/tamago/todo-with-gemini/backend/internal/models"
	"github.com/tamago/todo-with-gemini/backend/internal/services"
)

type AuthController struct {
	service services.AuthServiceInterface
}

func NewAuthController(service services.AuthServiceInterface) *AuthController {
	return &AuthController{service: service}
}

// Login handles user login and returns a JWT token.
func (ac *AuthController) Login(c *gin.Context) {
	_, span := otel.Tracer("").Start(c.Request.Context(), "AuthController.Login")
	defer span.End()

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := ac.service.Login(c.Request.Context(), user.Username, user.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
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

	if err := ac.service.Signup(c.Request.Context(), user.Username, user.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}