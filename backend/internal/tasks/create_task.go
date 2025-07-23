package tasks

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/tamago/todo-with-gemini/backend/internal/models"
)

// CreateTask handles the creation of a new task for the authenticated user.
func CreateTask(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		query := "INSERT INTO tasks (user_id, title, completed) VALUES ($1, $2, $3) RETURNING id"
		var id int
		err := db.QueryRow(query, userID, task.Title, false).Scan(&id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
			return
		}

		task.ID = id
		task.UserID = userID.(int)
		task.Completed = false

		c.JSON(http.StatusCreated, task)
	}
}
