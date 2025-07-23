package tasks

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/tamago/todo-with-gemini/backend/internal/models"
)

// UpdateTask handles updating an existing task.
func UpdateTask(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
			return
		}

		taskID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
			return
		}

		var task models.Task
		if err := c.ShouldBindJSON(&task); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Ensure the task belongs to the authenticated user
		var ownerID int
		err = db.QueryRow("SELECT user_id FROM tasks WHERE id = $1", taskID).Scan(&ownerID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve task owner"})
			return
		}

		if ownerID != userID.(int) {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update this task"})
			return
		}

		query := "UPDATE tasks SET title = $1, completed = $2 WHERE id = $3 AND user_id = $4"
		_, err = db.Exec(query, task.Title, task.Completed, taskID, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Task updated successfully"})
	}
}
