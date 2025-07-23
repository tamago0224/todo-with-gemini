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

func TestCreateTask(t *testing.T) {
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
	router.POST("/tasks", CreateTask(db))

	// Generate a valid token for testing
	token, err := utils.GenerateToken(1) // User ID 1
	assert.NoError(t, err)

	// Test case 1: Successful task creation
	task := map[string]interface{}{
		"title":     "New Task",
		"completed": false,
	}
	jsonValue, _ := json.Marshal(task)

	mock.ExpectQuery(`INSERT INTO tasks \(user_id, title, completed\) VALUES \(\$1, \$2, \$3\) RETURNING id`).
		WithArgs(1, "New Task", false).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token) // Set Authorization header

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "New Task")

	// Test case 2: Invalid JSON
	req, _ = http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token) // Set Authorization header
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Test case 3: Database error
	task = map[string]interface{}{
		"title":     "Another Task",
		"completed": false,
	}
	jsonValue, _ = json.Marshal(task)

	mock.ExpectQuery(`INSERT INTO tasks \(user_id, title, completed\) VALUES \(\$1, \$2, \$3\) RETURNING id`).
		WithArgs(1, "Another Task", false).
		WillReturnError(sql.ErrConnDone)

	req, _ = http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token) // Set Authorization header
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}