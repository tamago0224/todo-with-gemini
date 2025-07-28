package services

import (
	"context"

	"github.com/tamago/todo-with-gemini/backend/internal/models"
	"github.com/tamago/todo-with-gemini/backend/internal/repositories"
	"go.opentelemetry.io/otel"
)

type TaskService struct {
	repo repositories.TaskRepository
}

func NewTaskService(repo repositories.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) GetTasks(ctx context.Context, userID uint) ([]models.Task, error) {
	_, span := otel.Tracer("").Start(ctx, "TaskService.GetTasks")
	defer span.End()

	return s.repo.GetTasks(ctx, userID)
}

func (s *TaskService) CreateTask(ctx context.Context, task *models.Task, userID uint) (*models.Task, error) {
	_, span := otel.Tracer("").Start(ctx, "TaskService.CreateTask")
	defer span.End()

	task.UserID = userID
	task.Completed = false

	if err := s.repo.CreateTask(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, task *models.Task, taskID uint, userID uint) error {
	_, span := otel.Tracer("").Start(ctx, "TaskService.UpdateTask")
	defer span.End()

	task.ID = taskID
	task.UserID = userID

	// You might want to add logic here to check if the user is authorized to update the task

	return s.repo.UpdateTask(ctx, task)
}

func (s *TaskService) DeleteTask(ctx context.Context, taskID uint, userID uint) error {
	_, span := otel.Tracer("").Start(ctx, "TaskService.DeleteTask")
	defer span.End()

	// You might want to add logic here to check if the user is authorized to delete the task

	return s.repo.DeleteTask(ctx, taskID, userID)
}
