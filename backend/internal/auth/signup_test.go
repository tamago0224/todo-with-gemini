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
)

func TestSignup(t *testing.T) {
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
	router.POST("/signup", Signup(db))

	// Test case 1: Successful signup
	user := map[string]string{
		"username": "testuser",
		"password": "password123",
	}
	jsonValue, _ := json.Marshal(user)

	mock.ExpectQuery(`^INSERT INTO users`).
		WithArgs("testuser", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	req, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "User created successfully")

	// Test case 2: Missing username
	user = map[string]string{
		"password": "password123",
	}
	jsonValue, _ = json.Marshal(user)
	req, _ = http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Test case 3: Database error
	user = map[string]string{
		"username": "dbuser",
		"password": "password123",
	}
	jsonValue, _ = json.Marshal(user)

	mock.ExpectQuery(`^INSERT INTO users`).
		WithArgs("dbuser", sqlmock.AnyArg()).
		WillReturnError(sql.ErrConnDone)

	req, _ = http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// AnyString is a custom argument matcher for bcrypt hashed passwords.
type AnyString struct{}

func (a AnyString) Match(v interface{}) bool {
	_, ok := v.(string)
	return ok
}

func (a AnyString) String() string {
	return "AnyString()"
}
