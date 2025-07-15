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

	// ตัวอย่างตรวจสอบ username/password แบบ hardcoded
	if credentials.Username == "admin" && credentials.Password == "password" {
		json.NewEncoder(w).Encode(models.LoginResponse{Message: "Login successful"})
	} else {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
	}
}

// RegisterHandler returns a handler for registering a user or admin
func RegisterHandler(role string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var user models.Credentials
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		// Insert user into DB with role ("user" or "admin")
		_, err := db.DB.Exec("INSERT INTO users (username, password, role) VALUES (?, ?, ?)", user.Username, user.Password, role)
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
	var req struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// Update user in DB
	_, err := db.DB.Exec("UPDATE users SET username = ?, password = ? WHERE id = ?", req.Username, req.Password, req.ID)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "User updated"})
}

// DeleteUserHandler deletes a user
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req struct {
		ID int `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// Delete user from DB
	_, err := db.DB.Exec("DELETE FROM users WHERE id = ?", req.ID)
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
