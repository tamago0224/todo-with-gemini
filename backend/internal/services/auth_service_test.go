package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"github.com/tamago/todo-with-gemini/backend/internal/models"
	
)

// MockAuthRepository is a mock implementation of the AuthRepository interface
type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockAuthRepository) CreateUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func TestAuthService_Login(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	authService := NewAuthService(mockRepo)

	ctx := context.Background()
	username := "testuser"
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &models.User{ID: 1, Username: username, Password: string(hashedPassword)}
	mockRepo.On("GetUserByUsername", ctx, username).Return(user, nil)

	token, err := authService.Login(ctx, username, password)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	authService := NewAuthService(mockRepo)

	ctx := context.Background()
	username := "testuser"
	password := "wrongpassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	user := &models.User{ID: 1, Username: username, Password: string(hashedPassword)}
	mockRepo.On("GetUserByUsername", ctx, username).Return(user, nil)

	_, err := authService.Login(ctx, username, password)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid credentials")
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Signup(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	authService := NewAuthService(mockRepo)

	ctx := context.Background()
	username := "newuser"
	password := "newpassword123"

	mockRepo.On("CreateUser", ctx, mock.AnythingOfType("*models.User")).Return(nil)

	err := authService.Signup(ctx, username, password)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Signup_CreateUserError(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	authService := NewAuthService(mockRepo)

	ctx := context.Background()
	username := "newuser"
	password := "newpassword123"

	mockRepo.On("CreateUser", ctx, mock.AnythingOfType("*models.User")).Return(errors.New("db error"))

	err := authService.Signup(ctx, username, password)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
	mockRepo.AssertExpectations(t)
}
