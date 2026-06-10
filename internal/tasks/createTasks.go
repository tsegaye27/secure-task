package tasks

import (
	"encoding/json"
	"net/http"

	"task/internal/database"
	"task/internal/middleware"
	"task/internal/models"

	"github.com/google/uuid"
)

// CreateTaskHandler ...
// @Summary Create a task
// @Description Create a new task for the logged-in user
// @Tags tasks
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param task body models.Task true "Task details"
// @Success 201 {object} models.Task
// @Failure 401 {string} string "Unauthorized"
// @Router /tasks [post]
func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	task.ID = uuid.New().String()
	task.UserID = userID

	_, err := database.DB.Exec("INSERT INTO tasks (id, title, description, user_id) VALUES ($1, $2, $3, $4)",
		task.ID, task.Title, task.Description, task.UserID)
	if err != nil {
		http.Error(w, "Error creating task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}
