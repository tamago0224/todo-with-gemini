package tasks

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/tamago/todo-with-gemini/backend/internal/middleware"
	"github.com/tamago/todo-with-gemini/backend/internal/utils"
)

func TestUpdateTask(t *testing.T) {
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
	router.PUT("/tasks/:id", UpdateTask(db))

	// Generate a valid token for testing
	token, err := utils.GenerateToken(1) // User ID 1
	assert.NoError(t, err)

	// Test case 1: Successful task update
	task := map[string]interface{}{
		"title":     "Updated Task",
		"completed": true,
	}
	jsonValue, _ := json.Marshal(task)

	mock.ExpectQuery(`SELECT user_id FROM tasks WHERE id = \$1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(1))
	mock.ExpectExec(`UPDATE tasks SET title = \$1, completed = \$2 WHERE id = \$3 AND user_id = \$4`).
		WithArgs("Updated Task", true, 1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	req, _ := http.NewRequest(http.MethodPut, "/tasks/1", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token) // Set Authorization header

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Task updated successfully")

	// Test case 2: Task not found
	task = map[string]interface{}{
		"title":     "Nonexistent Task",
		"completed": false,
	}
	jsonValue, _ = json.Marshal(task)

	mock.ExpectQuery(`SELECT user_id FROM tasks WHERE id = \$1`).
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	req, _ = http.NewRequest(http.MethodPut, "/tasks/999", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token) // Set Authorization header
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	// Test case 3: Unauthorized update
	task = map[string]interface{}{
		"title":     "Unauthorized Task",
		"completed": false,
	}
	jsonValue, _ = json.Marshal(task)

	mock.ExpectQuery(`SELECT user_id FROM tasks WHERE id = \$1`).
		WithArgs(2).
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(2))

	req, _ = http.NewRequest(http.MethodPut, "/tasks/2", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token) // Set Authorization header
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}