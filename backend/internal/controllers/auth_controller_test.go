package controllers

import (
	"bytes"
	
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	"github.com/tamago/todo-with-gemini/backend/internal/models"
)

func TestAuthController_Signup(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	authController := NewAuthController(db)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	user := models.User{Username: "testuser", Password: "password123"}
	jsonValue, _ := json.Marshal(user)
	c.Request, _ = http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(jsonValue))
	c.Request.Header.Set("Content-Type", "application/json")

	mock.ExpectQuery("INSERT INTO users").WithArgs(user.Username, sqlmock.AnyArg()).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	authController.Signup(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthController_Login(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	authController := NewAuthController(db)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	username := "testuser"
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := models.User{Username: username, Password: password}
	jsonValue, _ := json.Marshal(user)
	c.Request, _ = http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonValue))
	c.Request.Header.Set("Content-Type", "application/json")

	mock.ExpectQuery(`SELECT id, password_hash FROM users WHERE username = $1`).WithArgs(username).WillReturnRows(sqlmock.NewRows([]string{"id", "password_hash"}).AddRow(1, hashedPassword))

	authController.Login(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}
