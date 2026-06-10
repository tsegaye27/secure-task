// Package auth
package auth

import (
	"encoding/json"
	"net/http"

	"task/internal/database"
	"task/internal/models"
)

// RegisterHandler ...
// @Summary Register a new user
// @Description Create a new user account with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.User true "User registration info"
// @Success 201 {object} map[string]string
// @Failure 400 {string} string "Invalid request payload"
// @Failure 409 {string} string "User already exists"
// @Router /auth/register [post]
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if user.Email == "" || user.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	var userID string
	err = database.DB.QueryRow("INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id", user.Email, hashedPassword).Scan(&userID)
	if err != nil {
		http.Error(w, "User already exists or database error", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully", "user_id": userID})
}
