package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tamago/todo-with-gemini/backend/internal/app/models"
	"github.com/tamago/todo-with-gemini/backend/internal/app/repositories"
)

// MockTaskRepository is a mock implementation of the TaskRepository interface

type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) GetTasks(ctx context.Context, userID uint) ([]models.Task, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Task), args.Error(1)
}

func (m *MockTaskRepository) CreateTask(ctx context.Context, task *models.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskRepository) UpdateTask(ctx context.Context, task *models.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskRepository) DeleteTask(ctx context.Context, taskID uint, userID uint) error {
	args := m.Called(ctx, taskID, userID)
	return args.Error(0)
}

func TestTaskService_GetTasks(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	taskService := NewTaskService(mockRepo)

	ctx := context.Background()
	userID := uint(1)

	tasks := []models.Task{{ID: 1, UserID: int(userID), Title: "Test Task"}}
	mockRepo.On("GetTasks", ctx, userID).Return(tasks, nil)

	result, err := taskService.GetTasks(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, tasks, result)
	mockRepo.AssertExpectations(t)
}

func TestTaskService_CreateTask(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	taskService := NewTaskService(mockRepo)

	ctx := context.Background()
	userID := uint(1)
	task := &models.Task{Title: "New Task"}

	mockRepo.On("CreateTask", ctx, mock.Anything).Return(nil)

	createdTask, err := taskService.CreateTask(ctx, task, userID)

	assert.NoError(t, err)
	assert.NotNil(t, createdTask)
	assert.Equal(t, int(userID), createdTask.UserID)
	assert.Equal(t, "New Task", createdTask.Title)
	mockRepo.AssertExpectations(t)
}

func TestTaskService_UpdateTask(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	taskService := NewTaskService(mockRepo)

	ctx := context.Background()
	userID := uint(1)
	taskID := uint(1)
	task := &models.Task{Title: "Updated Task"}

	mockRepo.On("UpdateTask", ctx, mock.Anything).Return(nil)

	err := taskService.UpdateTask(ctx, task, taskID, userID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTaskService_UpdateTask_NotFound(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	taskService := NewTaskService(mockRepo)

	ctx := context.Background()
	userID := uint(1)
	taskID := uint(1)
	task := &models.Task{Title: "Updated Task"}

	mockRepo.On("UpdateTask", ctx, mock.Anything).Return(repositories.ErrTaskNotFound)

	err := taskService.UpdateTask(ctx, task, taskID, userID)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, repositories.ErrTaskNotFound))
	mockRepo.AssertExpectations(t)
}

func TestTaskService_DeleteTask(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	taskService := NewTaskService(mockRepo)

	ctx := context.Background()
	userID := uint(1)
	taskID := uint(1)

	mockRepo.On("DeleteTask", ctx, taskID, userID).Return(nil)

	err := taskService.DeleteTask(ctx, taskID, userID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTaskService_DeleteTask_NotFound(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	taskService := NewTaskService(mockRepo)

	ctx := context.Background()
	userID := uint(1)
	taskID := uint(1)

	mockRepo.On("DeleteTask", ctx, taskID, userID).Return(repositories.ErrTaskNotFound)

	err := taskService.DeleteTask(ctx, taskID, userID)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, repositories.ErrTaskNotFound))
	mockRepo.AssertExpectations(t)
}

func TestTaskService_GetTasks_Error(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	taskService := NewTaskService(mockRepo)

	ctx := context.Background()
	userID := uint(1)

	mockRepo.On("GetTasks", ctx, userID).Return([]models.Task{}, errors.New("some error"))

	_, err := taskService.GetTasks(ctx, userID)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}
