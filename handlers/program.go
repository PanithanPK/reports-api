package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"reports-api/db"
	"reports-api/models"

	"github.com/gorilla/mux"
)

// AddProgramHandler handles POST requests to add new programs
func AddProgramHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req models.ProgramRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Name == "" {
		response := models.ProgramResponse{
			Success: false,
			Message: "Missing required field: name",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Insert into database
	query := `INSERT INTO program (name) VALUES (?)`
	_, err := db.DB.Exec(query, req.Name)
	if err != nil {
		log.Printf("Error inserting program: %v", err)
		response := models.ProgramResponse{
			Success: false,
			Message: "Failed to add program",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := models.ProgramResponse{
		Success: true,
		Message: "Program added successfully",
	}
	json.NewEncoder(w).Encode(response)
}

// GetProgramsHandler handles GET requests to retrieve all programs
func GetProgramsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := `SELECT id, name FROM program ORDER BY name`

	rows, err := db.DB.Query(query)
	if err != nil {
		log.Printf("Error querying programs: %v", err)
		response := models.GetProgramsResponse{
			Success: false,
			Message: "Failed to retrieve programs",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	defer rows.Close()

	var programs []models.Program
	for rows.Next() {
		var program models.Program
		err := rows.Scan(&program.ID, &program.Name)
		if err != nil {
			log.Printf("Error scanning program: %v", err)
			continue
		}
		programs = append(programs, program)
	}

	response := models.GetProgramsResponse{
		Success:  true,
		Message:  "Programs retrieved successfully",
		Programs: programs,
	}
	json.NewEncoder(w).Encode(response)
}

func UpdateProgramHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	var req models.ProgramRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Name == "" {
		response := models.ProgramResponse{
			Success: false,
			Message: "Missing required field: name",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Update database
	query := `UPDATE program SET name = ? WHERE id = ?`
	_, err := db.DB.Exec(query, req.Name, id)
	if err != nil {
		log.Printf("Error updating program: %v", err)
		response := models.ProgramResponse{
			Success: false,
			Message: "Failed to update program",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := models.ProgramResponse{
		Success: true,
		Message: "Program updated successfully",
	}
	json.NewEncoder(w).Encode(response)
}

func DeleteProgramHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	// Delete from database
	query := `DELETE FROM program WHERE id = ?`
	_, err := db.DB.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting program: %v", err)
		response := models.ProgramResponse{
			Success: false,
			Message: "Failed to delete program",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := models.ProgramResponse{
		Success: true,
		Message: "Program deleted successfully",
	}
	json.NewEncoder(w).Encode(response)
}

// DeleteAllPrograms ลบข้อมูลโปรแกรมทั้งหมด
func DeleteAllProgramsHandler(w http.ResponseWriter, r *http.Request) {
	database := db.DB

	// ลบข้อมูลทั้งหมดจากตาราง program และรีเซ็ต auto-increment
	// ใช้ DELETE FROM แทน TRUNCATE TABLE เพื่อหลีกเลี่ยง foreign key constraints
	_, err := database.Exec("DELETE FROM program")
	if err != nil {
		log.Printf("Error deleting all programs: %v", err)
		response := models.DeleteAllProgramsResponse{
			Success: false,
			Message: "Failed to delete all programs: " + err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// รีเซ็ต auto-increment counter
	_, err = database.Exec("ALTER TABLE program AUTO_INCREMENT = 1")
	if err != nil {
		log.Printf("Error resetting auto increment: %v", err)
		// ไม่ return error เพราะข้อมูลถูกลบแล้ว
	}

	response := models.DeleteAllProgramsResponse{
		Success: true,
		Message: "All programs deleted successfully",
	}
	json.NewEncoder(w).Encode(response)
}
