package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"reports-api/db"
	"reports-api/models"
	"time"
)

// LoginHandler handles user login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials models.Credentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var id int
	var username, password string
	var role string
	err := db.DB.QueryRow("SELECT id, username, password, role FROM users WHERE username = ? AND deleted_at IS NULL", credentials.Username).Scan(&id, &username, &password, &role)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Simple password comparison (will be replaced with bcrypt in future)
	if credentials.Password != password {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}
	w.Header().Set("role", role)           // Set user role in response header
	w.Header().Set("token", "dummy-token") // Set a dummy token header (will be replaced with JWT in future)

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    username, // Simple session value (will be replaced with secure token in future)
		Path:     "/",
		HttpOnly: true,
		MaxAge:   3600 * 24, // 24 hours
	})
	log.Printf("User %s logged in successfully", username)
	json.NewEncoder(w).Encode(models.LoginResponse{
		Message: "Login successful",
		Data:    &models.Data{ID: id, Username: username, Role: role},
	})
}

// RegisterHandler returns a handler for registering a user or admin
func RegisterHandler(role string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.RegisterUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Check if username already exists
		var count int
		err := db.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", req.Username).Scan(&count)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		if count > 0 {
			http.Error(w, "Username already exists", http.StatusConflict)
			return
		}

		_, err = db.DB.Exec(
			"INSERT INTO users (username, password, role, created_by, created_at) VALUES (?, ?, ?, ?, ?)",
			req.Username, req.Password, role, req.CreatedBy, time.Now(),
		)
		if err != nil {
			http.Error(w, "Failed to register user", http.StatusInternalServerError)
			return
		}
		log.Printf("User %s registered successfully as %s", req.Username, role)
		json.NewEncoder(w).Encode(map[string]string{"message": "Registered as " + role})
	}
}

// UpdateUserHandler updates user info
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if user exists
	var count int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM users WHERE id = ? AND deleted_at IS NULL", req.ID).Scan(&count)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if count == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	_, err = db.DB.Exec(
		"UPDATE users SET username = ?, password = ?, updated_by = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL",
		req.Username, req.Password, req.UpdatedBy, time.Now(), req.ID,
	)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}
	log.Printf("Updating user ID: %d with username: %s", req.ID, req.Username)
	json.NewEncoder(w).Encode(map[string]string{"message": "User updated"})
}

// DeleteUserHandler deletes a user
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	var req models.DeleteUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Use soft delete
	_, err := db.DB.Exec(
		"UPDATE users SET deleted_at = ?, deleted_by = ? WHERE id = ? AND deleted_at IS NULL",
		time.Now(), req.DeletedBy, req.ID,
	)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}
	log.Printf("User ID: %d deleted by user ID: %d", req.ID, req.DeletedBy)
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted"})
}

// LogoutHandler logs out a user
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Clear session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	log.Printf("User logged out successfully")
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out"})
}
