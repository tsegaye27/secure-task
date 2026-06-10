package tasks

import (
	"encoding/json"
	"net/http"

	"task/internal/database"
	"task/internal/middleware"
	"task/internal/models"
)

// CreateTaskHandler ...
// @Summary Create a task
// @Description Create a new task for the logged-in user
// @Tags tasks
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param task body models.CreateTaskRequest true "Task details"
// @Success 201 {object} models.Task
// @Failure 401 {string} string "Unauthorized"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Error creating task"
// @Router /tasks [post]
func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	var task models.Task
	err := database.DB.QueryRow(
		"INSERT INTO tasks (title, description, status, user_id) VALUES ($1, $2, $3, $4) RETURNING id, title, description, status, created_at, updated_at",
		req.Title, req.Description, req.Status, userID).Scan(
		&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		http.Error(w, "Error creating task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}
