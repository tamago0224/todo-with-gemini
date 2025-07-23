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

func TestDeleteTask(t *testing.T) {
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
	router.DELETE("/tasks/:id", DeleteTask(db))

	// Generate a valid token for testing
	token, err := utils.GenerateToken(1) // User ID 1
	assert.NoError(t, err)

	// Test case 1: Successful task deletion
	mock.ExpectQuery(`SELECT user_id FROM tasks WHERE id = \$1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(1))
	mock.ExpectExec(`DELETE FROM tasks WHERE id = \$1 AND user_id = \$2`).
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	req, _ := http.NewRequest(http.MethodDelete, "/tasks/1", nil)
	req.Header.Set("Authorization", "Bearer "+token) // Set Authorization header
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Task deleted successfully")

	// Test case 2: Task not found
	mock.ExpectQuery(`SELECT user_id FROM tasks WHERE id = \$1`).
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	req, _ = http.NewRequest(http.MethodDelete, "/tasks/999", nil)
	req.Header.Set("Authorization", "Bearer "+token) // Set Authorization header
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	// Test case 3: Unauthorized deletion
	mock.ExpectQuery(`SELECT user_id FROM tasks WHERE id = \$1`).
		WithArgs(2).
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(2))

	req, _ = http.NewRequest(http.MethodDelete, "/tasks/2", nil)
	req.Header.Set("Authorization", "Bearer "+token) // Set Authorization header
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}