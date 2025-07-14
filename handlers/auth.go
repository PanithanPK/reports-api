package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"reports-api/db"
	"reports-api/models"
	"time"
)

// LoginHandler handles user login requests
// It validates username/password against the database and returns success/failure response
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	LogRequest(UserLogger, r.Method, r.URL.Path, r.RemoteAddr)
	start := time.Now()

	w.Header().Set("Content-Type", "application/json")

	// Parse the login request from JSON body
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		LogValidationError(UserLogger, "request body", "invalid JSON", r.RemoteAddr)
		http.Error(w, "❌ Invalid request", http.StatusBadRequest)
		return
	}

	LogInfo(UserLogger, "Login attempt", "Username: "+req.Username)

	// Query database for user password and role
	var password string
	var role string
	err := db.DB.QueryRow("SELECT password, role FROM users WHERE username = ?", req.Username).Scan(&password, &role)
	if err == sql.ErrNoRows {
		// User not found in database
		LogSecurityEvent(UserLogger, "Login failed", "User not found: "+req.Username, r.RemoteAddr)
		json.NewEncoder(w).Encode(models.LoginResponse{Success: false, Message: "❌ ไม่พบชื่อผู้ใช้นี้ในระบบ"})
		return
	} else if err != nil {
		// Database error occurred
		LogError(UserLogger, "Database query", err, "Login for user: "+req.Username)
		http.Error(w, "❌ Database error", http.StatusInternalServerError)
		return
	}

	// Compare provided password with stored password
	if req.Password != password {
		// Password doesn't match
		LogSecurityEvent(UserLogger, "Login failed", "Invalid password for user: "+req.Username, r.RemoteAddr)
		json.NewEncoder(w).Encode(models.LoginResponse{Success: false, Message: "❌ รหัสผ่านไม่ถูกต้อง"})
		return
	}

	// Login successful
	LogSuccess(UserLogger, "User login", "User: "+req.Username+", Role: "+role)
	LogPerformance(UserLogger, "Login", time.Since(start), "User: "+req.Username)

	response := models.LoginResponse{
		Success: true,
		Message: "✅ เข้าสู่ระบบสำเร็จ",
		Role:    role,
	}
	json.NewEncoder(w).Encode(response)
}

// RegisterHandler creates a new user registration handler with specified role
// This is a closure that returns an http.HandlerFunc configured for a specific role
func RegisterHandler(role string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		LogRequest(UserLogger, r.Method, r.URL.Path, r.RemoteAddr)
		start := time.Now()

		w.Header().Set("Content-Type", "application/json")

		// Parse the registration request from JSON body
		var req models.RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			LogValidationError(UserLogger, "request body", "invalid JSON", r.RemoteAddr)
			http.Error(w, "❌ Invalid request", http.StatusBadRequest)
			return
		}

		LogInfo(UserLogger, "User registration", "Username: "+req.Username+", Role: "+role)

		// Insert new user into database with the specified role
		result, err := db.DB.Exec("INSERT INTO users (username, password, role) VALUES (?, ?, ?)", req.Username, req.Password, role)
		if err != nil {
			// Database error occurred during insertion
			LogError(UserLogger, "User registration", err, "Username: "+req.Username+", Role: "+role)
			http.Error(w, "❌ Database error", http.StatusInternalServerError)
			return
		}

		// Get the inserted ID
		id, _ := result.LastInsertId()
		LogSuccess(UserLogger, "User registration", "User ID: "+string(rune(id))+", Username: "+req.Username+", Role: "+role)
		LogPerformance(UserLogger, "Registration", time.Since(start), "User: "+req.Username)

		// Registration successful
		json.NewEncoder(w).Encode(models.RegisterResponse{Success: true, Message: "✅ ลงทะเบียนสำเร็จ"})
	}
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	LogRequest(UserLogger, r.Method, r.URL.Path, r.RemoteAddr)
	start := time.Now()

	w.Header().Set("Content-Type", "application/json")

	var req models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		LogValidationError(UserLogger, "request body", "invalid JSON", r.RemoteAddr)
		http.Error(w, "❌ Invalid request", http.StatusBadRequest)
		return
	}

	LogInfo(UserLogger, "User update", "User ID: "+string(rune(req.ID))+", Username: "+req.Username)

	result, err := db.DB.Exec("UPDATE users SET username = ?, password = ?, role = ? WHERE id = ?", req.Username, req.Password, req.Role, req.ID)
	if err != nil {
		LogError(UserLogger, "User update", err, "User ID: "+string(rune(req.ID)))
		http.Error(w, "❌ Database error", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	LogSuccess(UserLogger, "User update", "User ID: "+string(rune(req.ID))+", Rows affected: "+string(rune(rowsAffected)))
	LogPerformance(UserLogger, "User update", time.Since(start), "User ID: "+string(rune(req.ID)))

	json.NewEncoder(w).Encode(models.UpdateUserResponse{Success: true, Message: "✅ อัพเดตผู้ใช้สำเร็จ"})
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	LogRequest(UserLogger, r.Method, r.URL.Path, r.RemoteAddr)
	start := time.Now()

	w.Header().Set("Content-Type", "application/json")

	var req models.DeleteUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		LogValidationError(UserLogger, "request body", "invalid JSON", r.RemoteAddr)
		http.Error(w, "❌ Invalid request", http.StatusBadRequest)
		return
	}

	LogInfo(UserLogger, "User deletion", "User ID: "+string(rune(req.ID)))

	result, err := db.DB.Exec("DELETE FROM users WHERE id = ?", req.ID)
	if err != nil {
		LogError(UserLogger, "User deletion", err, "User ID: "+string(rune(req.ID)))
		http.Error(w, "❌ Database error", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	LogSuccess(UserLogger, "User deletion", "User ID: "+string(rune(req.ID))+", Rows affected: "+string(rune(rowsAffected)))
	LogPerformance(UserLogger, "User deletion", time.Since(start), "User ID: "+string(rune(req.ID)))

	json.NewEncoder(w).Encode(models.DeleteUserResponse{Success: true, Message: "✅ ลบผู้ใช้สำเร็จ"})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	LogRequest(UserLogger, r.Method, r.URL.Path, r.RemoteAddr)

	w.Header().Set("Content-Type", "application/json")

	cookie := http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour * 24),
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)

	LogSuccess(UserLogger, "User logout", "User logged out successfully")
	json.NewEncoder(w).Encode(models.LogoutResponse{Success: true, Message: "✅ ออกจากระบบสำเร็จ"})
}
