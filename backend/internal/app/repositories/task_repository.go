package repositories

import (
	"context"
	"database/sql"

	"github.com/tamago/todo-with-gemini/backend/internal/app/models"
	"go.opentelemetry.io/otel"
)

type TaskRepository interface {
	GetTasks(ctx context.Context, userID uint) ([]models.Task, error)
	CreateTask(ctx context.Context, task *models.Task) error
	UpdateTask(ctx context.Context, task *models.Task) error
	DeleteTask(ctx context.Context, taskID uint, userID uint) error
}

type PostgresTaskRepository struct {
	db *sql.DB
}

func NewPostgresTaskRepository(db *sql.DB) *PostgresTaskRepository {
	return &PostgresTaskRepository{db: db}
}

func (r *PostgresTaskRepository) GetTasks(ctx context.Context, userID uint) ([]models.Task, error) {
	_, span := otel.Tracer("").Start(ctx, "TaskRepository.GetTasks")
	defer span.End()

	rows, err := r.db.QueryContext(ctx, "SELECT id, user_id, title, completed FROM tasks WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []models.Task{}
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.UserID, &task.Title, &task.Completed); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, rows.Err()
}

func (r *PostgresTaskRepository) CreateTask(ctx context.Context, task *models.Task) error {
	_, span := otel.Tracer("").Start(ctx, "TaskRepository.CreateTask")
	defer span.End()

	query := "INSERT INTO tasks (user_id, title, completed) VALUES ($1, $2, $3) RETURNING id"
	var id int
	err := r.db.QueryRowContext(ctx, query, task.UserID, task.Title, task.Completed).Scan(&id)
	if err != nil {
		return err
	}
	task.ID = id

	return nil
}

func (r *PostgresTaskRepository) UpdateTask(ctx context.Context, task *models.Task) error {
	_, span := otel.Tracer("").Start(ctx, "TaskRepository.UpdateTask")
	defer span.End()

	query := "UPDATE tasks SET title = $1, completed = $2 WHERE id = $3 AND user_id = $4"
	_, err := r.db.ExecContext(ctx, query, task.Title, task.Completed, task.ID, task.UserID)
	return err
}

func (r *PostgresTaskRepository) DeleteTask(ctx context.Context, taskID uint, userID uint) error {
	utils.RandomSleep()
	_, span := otel.Tracer("").Start(ctx, "TaskRepository.DeleteTask")
	defer span.End()

	query := "DELETE FROM tasks WHERE id = $1 AND user_id = $2"
	_, err := r.db.ExecContext(ctx, query, taskID, userID)
	return err
}