package tasks

import (
	"encoding/json"
	"net/http"
	"time"

	"task/internal/database"
	"task/internal/middleware"
	"task/internal/models"
)

// UpdateTaskHandler ...
// @Summary Update a task
// @Description Update an existing task by ID
// @Tags tasks
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Param task body models.UpdateTaskRequest true "Updated task details"
// @Success 200 {object} map[string]string
// @Failure 401 {string} string "Unauthorized"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Error updating task or task not found"
// @Router /tasks/{id} [put]
func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	var req models.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	result, err := database.DB.Exec("UPDATE tasks SET title = $1, description = $2, status = $3, updated_at = $4 WHERE id = $5 AND user_id = $6",
		req.Title, req.Description, req.Status, time.Now(), taskID, userID)
	if err != nil {
		http.Error(w, "Error updating task", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Task updated successfully"})
}
