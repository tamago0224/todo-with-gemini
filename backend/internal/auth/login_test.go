package auth

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
	"golang.org/x/crypto/bcrypt"
)

func TestLogin(t *testing.T) {
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
	router.POST("/login", Login(db))

	// Hash a sample password for testing
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	// Test case 1: Successful login
	user := map[string]string{
		"username": "testuser",
		"password": "password123",
	}
	jsonValue, _ := json.Marshal(user)

	mock.ExpectQuery(`^SELECT id, password_hash FROM users WHERE username = \$1$`).
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"id", "password_hash"}).AddRow(1, hashedPassword))

	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "token")

	// Test case 2: Invalid credentials (wrong password)
	user = map[string]string{
		"username": "testuser",
		"password": "wrongpassword",
	}
	jsonValue, _ = json.Marshal(user)

	mock.ExpectQuery(`^SELECT id, password_hash FROM users WHERE username = \$1$`).
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"id", "password_hash"}).AddRow(1, hashedPassword))

	req, _ = http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Test case 3: User not found
	user = map[string]string{
		"username": "nonexistentuser",
		"password": "password123",
	}
	jsonValue, _ = json.Marshal(user)

	mock.ExpectQuery(`^SELECT id, password_hash FROM users WHERE username = \$1$`).
		WithArgs("nonexistentuser").
		WillReturnError(sql.ErrNoRows)

	req, _ = http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Test case 4: Database error during user retrieval
	user = map[string]string{
		"username": "dbuser",
		"password": "password123",
	}
	jsonValue, _ = json.Marshal(user)

	mock.ExpectQuery(`^SELECT id, password_hash FROM users WHERE username = \$1$`).
		WithArgs("dbuser").
		WillReturnError(sql.ErrConnDone)

	req, _ = http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}