package tasks

import (
	"encoding/json"
	"net/http"

	"task/internal/database"
	"task/internal/middleware"
	"task/internal/models"
)

// GetTaskHandler ...
// @Summary Get a task
// @Description Get details of a specific task by ID
// @Tags tasks
// @Security BearerAuth
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {object} models.Task
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Task not found"
// @Router /tasks/{id} [get]
func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	taskID := r.PathValue("id")
	if taskID == "" {
		http.Error(w, "Task ID is required", http.StatusBadRequest)
		return
	}

	var task models.Task
	err := database.DB.QueryRow("SELECT id, title, description, status, user_id, created_at, updated_at FROM tasks WHERE id = $1 AND user_id = $2", taskID, userID).
		Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.UserID, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(task)
}
