package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/tamago/todo-with-gemini/backend/internal/app/models"
	"github.com/tamago/todo-with-gemini/backend/internal/app/services"
)

// MockAuthService is a mock implementation of the AuthServiceInterface
type MockAuthService struct {
	mock.Mock
}

// Statically assert that MockAuthService implements the interface.
var _ services.AuthServiceInterface = (*MockAuthService)(nil)

func (m *MockAuthService) Login(ctx context.Context, username, password string) (string, error) {
	args := m.Called(ctx, username, password)
	return args.Get(0).(string), args.Error(1)
}

func (m *MockAuthService) Signup(ctx context.Context, username, password string) error {
	args := m.Called(ctx, username, password)
	return args.Error(0)
}

func TestAuthController_Signup(t *testing.T) {
	mockService := new(MockAuthService)
	authController := NewAuthController(mockService)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	user := models.User{Username: "testuser", Password: "password123"}
	jsonValue, _ := json.Marshal(user)
	c.Request, _ = http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(jsonValue))
	c.Request.Header.Set("Content-Type", "application/json")

	mockService.On("Signup", mock.Anything, user.Username, user.Password).Return(nil)

	authController.Signup(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestAuthController_Signup_Error(t *testing.T) {
	mockService := new(MockAuthService)
	authController := NewAuthController(mockService)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	user := models.User{Username: "testuser", Password: "password123"}
	jsonValue, _ := json.Marshal(user)
	c.Request, _ = http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(jsonValue))
	c.Request.Header.Set("Content-Type", "application/json")

	mockService.On("Signup", mock.Anything, user.Username, user.Password).Return(errors.New("service error"))

	authController.Signup(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestAuthController_Login(t *testing.T) {
	mockService := new(MockAuthService)
	authController := NewAuthController(mockService)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	username := "testuser"
	password := "password123"

	user := models.User{Username: username, Password: password}
	jsonValue, _ := json.Marshal(user)
	c.Request, _ = http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonValue))
	c.Request.Header.Set("Content-Type", "application/json")

	mockService.On("Login", mock.Anything, username, password).Return("dummy_token", nil)

	authController.Login(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestAuthController_Login_InvalidCredentials(t *testing.T) {
	mockService := new(MockAuthService)
	authController := NewAuthController(mockService)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	username := "testuser"
	password := "wrongpassword"

	user := models.User{Username: username, Password: password}
	jsonValue, _ := json.Marshal(user)
	c.Request, _ = http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonValue))
	c.Request.Header.Set("Content-Type", "application/json")

	mockService.On("Login", mock.Anything, username, password).Return("", errors.New("Invalid credentials"))

	authController.Login(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockService.AssertExpectations(t)
}