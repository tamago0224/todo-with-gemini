package tasks

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/tamago/todo-with-gemini/backend/internal/middleware"
	"github.com/tamago/todo-with-gemini/backend/internal/utils"
)

func TestGetTasks(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			return
		}
	}(db)

	router := gin.Default()
	router.Use(middleware.AuthMiddleware()) // Apply AuthMiddleware
	router.GET("/tasks", GetTasks(db))

	// Generate a valid token for testing
	token, err := utils.GenerateToken(1) // User ID 1
	assert.NoError(t, err)

	// Test case 1: Successful retrieval of tasks
	mock.ExpectQuery(`SELECT id, user_id, title, completed FROM tasks WHERE user_id = \$1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "title", "completed"}).
			AddRow(1, 1, "Task 1", false).
			AddRow(2, 1, "Task 2", true))

	req, _ := http.NewRequest(http.MethodGet, "/tasks", nil)
	req.Header.Set("Authorization", "Bearer "+token) // Set Authorization header
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Task 1")
	assert.Contains(t, w.Body.String(), "Task 2")

	// Test case 2: No tasks found
	mock.ExpectQuery(`SELECT id, user_id, title, completed FROM tasks WHERE user_id = \$1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "title", "completed"}))

	req, _ = http.NewRequest(http.MethodGet, "/tasks", nil)
	req.Header.Set("Authorization", "Bearer "+token) // Set Authorization header
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "[]", w.Body.String())

	// Test case 3: Database error
	mock.ExpectQuery(`SELECT id, user_id, title, completed FROM tasks WHERE user_id = \$1`).
		WithArgs(1).
		WillReturnError(sql.ErrConnDone)

	req, _ = http.NewRequest(http.MethodGet, "/tasks", nil)
	req.Header.Set("Authorization", "Bearer "+token) // Set Authorization header
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}