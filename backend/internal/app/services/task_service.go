package services

import (
	"context"
	"errors"

	"github.com/tamago/todo-with-gemini/backend/internal/app/models"
	"github.com/tamago/todo-with-gemini/backend/internal/app/repositories"
	"github.com/tamago/todo-with-gemini/backend/internal/platform/utils"
	"go.opentelemetry.io/otel"
)

type TaskServiceInterface interface {
	GetTasks(ctx context.Context, userID uint) ([]models.Task, error)
	CreateTask(ctx context.Context, task *models.Task, userID uint) (*models.Task, error)
	UpdateTask(ctx context.Context, task *models.Task, taskID uint, userID uint) error
	DeleteTask(ctx context.Context, taskID uint, userID uint) error
}

type TaskService struct {
	repo repositories.TaskRepository
}

func NewTaskService(repo repositories.TaskRepository) TaskServiceInterface {
	return &TaskService{repo: repo}
}

func (s *TaskService) GetTasks(ctx context.Context, userID uint) ([]models.Task, error) {
	_, span := otel.Tracer("").Start(ctx, "TaskService.GetTasks")
	defer span.End()

	utils.RandomSleep()
	return s.repo.GetTasks(ctx, userID)
}

func (s *TaskService) CreateTask(ctx context.Context, task *models.Task, userID uint) (*models.Task, error) {
	_, span := otel.Tracer("").Start(ctx, "TaskService.CreateTask")
	defer span.End()

	utils.RandomSleep()
	task.UserID = int(userID)
	task.Completed = false

	if err := s.repo.CreateTask(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, task *models.Task, taskID uint, userID uint) error {
	_, span := otel.Tracer("").Start(ctx, "TaskService.UpdateTask")
	defer span.End()

	utils.RandomSleep()
	task.ID = int(taskID)
	task.UserID = int(userID)

	// You might want to add logic here to check if the user is authorized to update the task

	err := s.repo.UpdateTask(ctx, task)
	if errors.Is(err, repositories.ErrTaskNotFound) {
		return err
	}
	return nil
}

func (s *TaskService) DeleteTask(ctx context.Context, taskID uint, userID uint) error {
	_, span := otel.Tracer("").Start(ctx, "TaskService.DeleteTask")
	defer span.End()

	// You might want to add logic here to check if the user is authorized to delete the task

	utils.RandomSleep()
	err := s.repo.DeleteTask(ctx, taskID, userID)
	if errors.Is(err, repositories.ErrTaskNotFound) {
		return err
	}
	return nil
}
