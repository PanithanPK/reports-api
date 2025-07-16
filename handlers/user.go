package handlers

import (
	"encoding/json"
	"net/http"
	"reports-api/db"
	"reports-api/models"
)

// LoginHandler handles user login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var credentials models.Credentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var id int
	var username, password string
	err := db.DB.QueryRow("SELECT id, username, password FROM users WHERE username = ? AND deleted_at IS NULL", credentials.Username).Scan(&id, &username, &password)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	if credentials.Password != password {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(models.LoginResponse{Message: "Login successful", Data: &models.Data{ID: id, Username: username}})
}

// RegisterHandler returns a handler for registering a user or admin
func RegisterHandler(role string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var req models.RegisterUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		_, err := db.DB.Exec("INSERT INTO users (username, password, role, created_by) VALUES (?, ?, ?, ?)", req.Username, req.Password, role, req.CreatedBy)
		if err != nil {
			http.Error(w, "Failed to register user", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Registered as " + role})
	}
}

// UpdateUserHandler updates user info
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	_, err := db.DB.Exec("UPDATE users SET username = ?, password = ?, updated_by = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ? AND deleted_at IS NULL", req.Username, req.Password, req.UpdatedBy, req.ID)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "User updated"})
}

// DeleteUserHandler deletes a user
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req models.DeleteUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	_, err := db.DB.Exec("UPDATE users SET deleted_at = CURRENT_TIMESTAMP, deleted_by = ? WHERE id = ? AND deleted_at IS NULL", req.DeletedBy, req.ID)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted"})
}

// LogoutHandler logs out a user
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// ลบ session cookie (ถ้ามี)
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out"})
}
