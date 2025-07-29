package services

import (
	"context"
	"errors"

	"github.com/tamago/todo-with-gemini/backend/internal/app/models"
	"github.com/tamago/todo-with-gemini/backend/internal/app/repositories"
	"github.com/tamago/todo-with-gemini/backend/internal/platform/utils"
	"go.opentelemetry.io/otel"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceInterface interface {
	Login(ctx context.Context, username, password string) (string, error)
	Signup(ctx context.Context, username, password string) error
}

type AuthService struct {
	repo repositories.AuthRepository
}

func NewAuthService(repo repositories.AuthRepository) AuthServiceInterface {
	return &AuthService{repo: repo}
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	_, span := otel.Tracer("").Start(ctx, "AuthService.Login")
	defer span.End()

	utils.RandomSleep()
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return "", errors.New("Invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("Invalid credentials")
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return "", errors.New("Failed to generate token")
	}

	return token, nil
}

func (s *AuthService) Signup(ctx context.Context, username, password string) error {
	_, span := otel.Tracer("").Start(ctx, "AuthService.Signup")
	defer span.End()

	utils.RandomSleep()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("Failed to hash password")
	}

	user := &models.User{Username: username, Password: string(hashedPassword)}

	return s.repo.CreateUser(ctx, user)
}
