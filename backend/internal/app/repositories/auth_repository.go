package repositories

import (
	"context"
	"database/sql"

	"github.com/tamago/todo-with-gemini/backend/internal/app/models"
	"go.opentelemetry.io/otel"
)

type AuthRepository interface {
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
}

type PostgresAuthRepository struct {
	db *sql.DB
}

func NewPostgresAuthRepository(db *sql.DB) *PostgresAuthRepository {
	return &PostgresAuthRepository{db: db}
}

func (r *PostgresAuthRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	_, span := otel.Tracer("").Start(ctx, "AuthRepository.GetUserByUsername")
	defer span.End()

	var user models.User
	var storedPasswordHash string
	query := "SELECT id, username, password_hash FROM users WHERE username = $1"
	err := r.db.QueryRowContext(ctx, query, username).Scan(&user.ID, &user.Username, &storedPasswordHash)
	if err != nil {
		return nil, err
	}
	user.Password = storedPasswordHash // Temporarily store hash in Password field

	return &user, nil
}

func (r *PostgresAuthRepository) CreateUser(ctx context.Context, user *models.User) error {
	_, span := otel.Tracer("").Start(ctx, "AuthRepository.CreateUser")
	defer span.End()

	query := "INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id"
	var id int
	err := r.db.QueryRowContext(ctx, query, user.Username, user.Password).Scan(&id)
	if err != nil {
		return err
	}
	user.ID = id

	return nil
}
