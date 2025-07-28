
package repositories

import (
	"context"

	"github.com/tamago/todo-with-gemini/backend/internal/models"
)

type TaskRepository interface {
	GetTasks(ctx context.Context, userID uint) ([]models.Task, error)
	CreateTask(ctx context.Context, task *models.Task) error
	UpdateTask(ctx context.Context, task *models.Task) error
	DeleteTask(ctx context.Context, taskID uint, userID uint) error
}
