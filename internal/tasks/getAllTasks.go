package tasks

import (
	"encoding/json"
	"net/http"

	"task/internal/database"
	"task/internal/middleware"
	"task/internal/models"
)

// GetAllTasksHandler ...
// @Summary Get all tasks
// @Description List all tasks for the logged-in user
// @Tags tasks
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.Task
// @Failure 401 {string} string "Unauthorized"
// @Router /tasks [get]
func GetAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	rows, err := database.DB.Query("SELECT id, title, description, status, user_id, created_at, updated_at FROM tasks WHERE user_id = $1", userID)
	if err != nil {
		http.Error(w, "Error fetching tasks", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.UserID, &task.CreatedAt, &task.UpdatedAt); err != nil {
			http.Error(w, "Error scanning tasks", http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, task)
	}

	json.NewEncoder(w).Encode(tasks)
}
