package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"reports-api/db"
	"reports-api/models"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// AddUserHandler handles POST requests to add new users
func AddUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req models.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if err := validate.Struct(req); err != nil {
		response := models.UserResponse{
			Success: false,
			Message: err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Validate required fields
	if req.Username == "" || req.Password == "" {
		response := models.UserResponse{
			Success: false,
			Message: "Missing required fields: username, password",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Set default role if not provided
	if req.Role == "" {
		req.Role = "user"
	}

	// Validate role
	if req.Role != "admin" && req.Role != "user" {
		response := models.UserResponse{
			Success: false,
			Message: "Invalid role. Must be 'admin' or 'user'",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Insert into database
	query := `INSERT INTO users (username, password, role) VALUES (?, ?, ?)`
	_, err := db.DB.Exec(query, req.Username, req.Password, req.Role)
	if err != nil {
		log.Printf("Error inserting user: %v", err)

		// Check if error is due to duplicate username
		if strings.Contains(err.Error(), "Duplicate entry") && strings.Contains(err.Error(), "users.username") {
			response := models.UserResponse{
				Success: false,
				Message: "ชื่อผู้ใช้นี้มีอยู่ในระบบแล้ว กรุณาใช้ชื่อผู้ใช้อื่น",
			}
			json.NewEncoder(w).Encode(response)
			return
		}

		response := models.UserResponse{
			Success: false,
			Message: "Failed to add user",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := models.UserResponse{
		Success: true,
		Message: "User added successfully",
	}
	json.NewEncoder(w).Encode(response)
}

// GetUsersHandler handles GET requests to retrieve all users
func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := `SELECT id, username, role FROM users ORDER BY username`

	rows, err := db.DB.Query(query)
	if err != nil {
		log.Printf("Error querying users: %v", err)
		response := models.GetUsersResponse{
			Success: false,
			Message: "Failed to retrieve users",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username, &user.Role)
		if err != nil {
			log.Printf("Error scanning user: %v", err)
			continue
		}
		users = append(users, user)
	}

	response := models.GetUsersResponse{
		Success: true,
		Message: "Users retrieved successfully",
		Users:   users,
	}
	json.NewEncoder(w).Encode(response)
}
