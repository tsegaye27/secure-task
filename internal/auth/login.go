package auth

import (
	"encoding/json"
	"net/http"

	"task/internal/database"
	"task/internal/models"
)

// LoginHandler ...
// @Summary Login user
// @Description Authenticate user and return JWT token and userID
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.User true "User login credentials"
// @Success 200 {object} map[string]string
// @Failure 401 {string} string "Invalid email or password"
// @Router /auth/login [post]
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials models.User
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var userID string
	var hashedPassword string
	err := database.DB.QueryRow("SELECT id, password FROM users WHERE email = $1", credentials.Email).Scan(&userID, &hashedPassword)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	if !CheckPasswordHash(credentials.Password, hashedPassword) {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	token, err := GenerateJWT(userID)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token, "user_id": userID})
}
