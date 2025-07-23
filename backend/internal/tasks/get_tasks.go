package tasks

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/tamago/todo-with-gemini/backend/internal/models"
)

// GetTasks retrieves all tasks for the authenticated user.
func GetTasks(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
			return
		}

		rows, err := db.Query("SELECT id, user_id, title, completed FROM tasks WHERE user_id = $1", userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
			return
		}
		defer func(rows *sql.Rows) {
			err := rows.Close()
			if err != nil {
				// Log the error, but don't return it as it's not critical for the response
				return
			}
		}(rows)

		tasks := []models.Task{}
		for rows.Next() {
			var task models.Task
			if err := rows.Scan(&task.ID, &task.UserID, &task.Title, &task.Completed); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan task"})
				return
			}
			tasks = append(tasks, task)
		}

		if err = rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating tasks"})
			return
		}

		c.JSON(http.StatusOK, tasks)
	}
}
