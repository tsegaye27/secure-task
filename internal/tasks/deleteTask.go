// Package tasks
package tasks

import (
	"net/http"

	"task/internal/database"
	"task/internal/middleware"
)

// DeleteTaskHandler ...
// @Summary Delete a task
// @Description Remove a task by ID
// @Tags tasks
// @Security BearerAuth
// @Param id path string true "Task ID"
// @Success 204 "No Content"
// @Failure 401 {string} string "Unauthorized"
// @Router /tasks/{id} [delete]
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	_, err := database.DB.Exec("DELETE FROM tasks WHERE id = $1 AND user_id = $2", taskID, userID)
	if err != nil {
		http.Error(w, "Error deleting task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
