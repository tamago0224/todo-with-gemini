package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tamago/todo-with-gemini/backend/internal/app/models"
	"github.com/tamago/todo-with-gemini/backend/internal/app/services"
	"github.com/tamago/todo-with-gemini/backend/internal/platform/utils"
	"go.opentelemetry.io/otel"
)

type TaskController struct {
	service services.TaskServiceInterface
}

func NewTaskController(service services.TaskServiceInterface) *TaskController {
	return &TaskController{service: service}
}

func (tc *TaskController) GetTasks(c *gin.Context) {
	_, span := otel.Tracer("TaskController").Start(c.Request.Context(), "TaskController.GetTasks")
	defer span.End()

	utils.RandomSleep()
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	tasks, err := tc.service.GetTasks(c.Request.Context(), uint(userID.(int)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (tc *TaskController) CreateTask(c *gin.Context) {
	_, span := otel.Tracer("").Start(c.Request.Context(), "TaskController.CreateTask")
	defer span.End()

	utils.RandomSleep()
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdTask, err := tc.service.CreateTask(c.Request.Context(), &task, uint(userID.(int)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	c.JSON(http.StatusCreated, createdTask)
}

func (tc *TaskController) UpdateTask(c *gin.Context) {
	_, span := otel.Tracer("").Start(c.Request.Context(), "TaskController.UpdateTask")
	defer span.End()

	utils.RandomSleep()
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := tc.service.UpdateTask(c.Request.Context(), &task, uint(taskID), uint(userID.(int))); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task updated successfully"})
}

func (tc *TaskController) DeleteTask(c *gin.Context) {
	_, span := otel.Tracer("").Start(c.Request.Context(), "TaskController.DeleteTask")
	defer span.End()

	utils.RandomSleep()
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	if err := tc.service.DeleteTask(c.Request.Context(), uint(taskID), uint(userID.(int))); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}
