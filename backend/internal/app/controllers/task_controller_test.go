package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tamago/todo-with-gemini/backend/internal/app/models"
	"github.com/tamago/todo-with-gemini/backend/internal/app/services"
)

// MockTaskService is a mock that implements the TaskServiceInterface
type MockTaskService struct {
	mock.Mock
}

// Statically assert that MockTaskService implements the interface.
var _ services.TaskServiceInterface = (*MockTaskService)(nil)

func (m *MockTaskService) GetTasks(ctx context.Context, userID uint) ([]models.Task, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Task), args.Error(1)
}

func (m *MockTaskService) CreateTask(ctx context.Context, task *models.Task, userID uint) (*models.Task, error) {
	args := m.Called(ctx, task, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Task), args.Error(1)
}

func (m *MockTaskService) UpdateTask(ctx context.Context, task *models.Task, taskID uint, userID uint) error {
	args := m.Called(ctx, task, taskID, userID)
	return args.Error(0)
}

func (m *MockTaskService) DeleteTask(ctx context.Context, taskID uint, userID uint) error {
	args := m.Called(ctx, taskID, userID)
	return args.Error(0)
}

func TestTaskController_GetTasks(t *testing.T) {
	mockService := new(MockTaskService)
	taskController := NewTaskController(mockService)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/tasks", nil)
	c.Set("userID", 1)

	tasks := []models.Task{{ID: 1, UserID: 1, Title: "Test Task"}}
	mockService.On("GetTasks", mock.Anything, uint(1)).Return(tasks, nil)

	taskController.GetTasks(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestTaskController_CreateTask(t *testing.T) {
	mockService := new(MockTaskService)
	taskController := NewTaskController(mockService)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", 1)

	task := &models.Task{Title: "New Task"}
	jsonValue, _ := json.Marshal(task)
	c.Request, _ = http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(jsonValue))
	c.Request.Header.Set("Content-Type", "application/json")

	createdTask := &models.Task{ID: 1, UserID: 1, Title: "New Task", Completed: false}
	mockService.On("CreateTask", mock.Anything, mock.AnythingOfType("*models.Task"), uint(1)).Return(createdTask, nil)

	taskController.CreateTask(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestTaskController_UpdateTask(t *testing.T) {
	mockService := new(MockTaskService)
	taskController := NewTaskController(mockService)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", 1)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	task := &models.Task{Title: "Updated Task"}
	jsonValue, _ := json.Marshal(task)
	c.Request, _ = http.NewRequest(http.MethodPut, "/tasks/1", bytes.NewBuffer(jsonValue))
	c.Request.Header.Set("Content-Type", "application/json")

	mockService.On("UpdateTask", mock.Anything, mock.AnythingOfType("*models.Task"), uint(1), uint(1)).Return(nil)

	taskController.UpdateTask(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestTaskController_DeleteTask(t *testing.T) {
	mockService := new(MockTaskService)
	taskController := NewTaskController(mockService)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodDelete, "/tasks/1", nil)
	c.Set("userID", 1)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

			mockService.On("DeleteTask", mock.Anything, uint(1), uint(1)).Return(nil)

	taskController.DeleteTask(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}