package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"reports-api/db"
	"reports-api/models"

	"github.com/gorilla/mux"
)

// ส่วนในการเพิ่มข้อมูลสาขา
func AddBranchOfficeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // ตั้งค่า Content-Type เป็น JSON

	var req models.BranchOfficeRequest                           // รับข้อมูลจาก request body ที่อยู่ในโฟลเดอร์ models โดยใช้ branchoffice
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { // เงื่อนไขถ้าไม่สามารถถอดข้อมูลจาก request body ได้
		http.Error(w, "Invalid request body", http.StatusBadRequest) // ส่งข้อความ error ถ้า request body ไม่ถูกต้อง
		return
	}

	// ตรวจสอบว่าค่า ipPhone และ branchoffice ไม่เป็นค่าว่างหรือไม่
	if req.IpPhone == "" || req.Branchoffice == "" {
		response := models.BranchOfficeResponse{ // ส่งข้อความ error
			Success: false,
			Message: "Missing required fields: ipPhone, branchoffice",
		}
		json.NewEncoder(w).Encode(response) // ส่งข้อความ error ถ้า request body ไม่ถูกต้อง
		return
	}

	// เพิ่มข้อมูลเข้า database
	query := `INSERT INTO branch_office (ip_phone, branchoffice) VALUES (?, ?)`
	_, err := db.DB.Exec(query, req.IpPhone, req.Branchoffice)
	if err != nil {
		log.Printf("Error inserting branch office: %v", err)
		response := models.BranchOfficeResponse{
			Success: false,
			Message: "Failed to add branch office",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := models.BranchOfficeResponse{
		Success: true,
		Message: "Branch office added successfully",
	}
	json.NewEncoder(w).Encode(response)
}

// UpdateBranchOfficeHandler handles PUT requests to update existing branch offices
func UpdateBranchOfficeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	ipPhone := vars["ip_phone"]

	var req models.BranchOfficeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.IpPhone == "" || req.Branchoffice == "" {
		response := models.BranchOfficeResponse{
			Success: false,
			Message: "Missing required fields: ipPhone, branchoffice",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Update database
	query := `UPDATE branch_office SET branchoffice = ? WHERE ip_phone = ?`
	result, err := db.DB.Exec(query, req.Branchoffice, ipPhone)
	if err != nil {
		log.Printf("Error updating branch office: %v", err)
		response := models.BranchOfficeResponse{
			Success: false,
			Message: "Failed to update branch office",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Check if any row was affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		response := models.BranchOfficeResponse{
			Success: false,
			Message: "Failed to verify update",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if rowsAffected == 0 {
		response := models.BranchOfficeResponse{
			Success: false,
			Message: "Branch office not found",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := models.BranchOfficeResponse{
		Success: true,
		Message: "Branch office updated successfully",
	}
	json.NewEncoder(w).Encode(response)
}

// DeleteBranchOfficeHandler handles DELETE requests to delete branch offices
func DeleteBranchOfficeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	ipPhone := vars["ip_phone"]

	// Delete from database
	query := `DELETE FROM branch_office WHERE ip_phone = ?`
	result, err := db.DB.Exec(query, ipPhone)
	if err != nil {
		log.Printf("Error deleting branch office: %v", err)
		response := models.BranchOfficeResponse{
			Success: false,
			Message: "Failed to delete branch office",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Check if any row was affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		response := models.BranchOfficeResponse{
			Success: false,
			Message: "Failed to verify deletion",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if rowsAffected == 0 {
		response := models.BranchOfficeResponse{
			Success: false,
			Message: "Branch office not found",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := models.BranchOfficeResponse{
		Success: true,
		Message: "Branch office deleted successfully",
	}
	json.NewEncoder(w).Encode(response)
}

// GetBranchOfficesHandler handles GET requests to retrieve all branch offices
func GetBranchOfficesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Query branch offices
	branchQuery := `SELECT ip_phone, branchoffice FROM branch_office ORDER BY branchoffice`
	branchRows, err := db.DB.Query(branchQuery)
	if err != nil {
		log.Printf("Error querying branch offices: %v", err)
		response := models.GetBranchOfficesResponse{
			Success: false,
			Message: "Failed to retrieve branch offices",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	defer branchRows.Close()

	var branchOffices []models.BranchOffice
	for branchRows.Next() {
		var bo models.BranchOffice
		err := branchRows.Scan(&bo.IpPhone, &bo.Branchoffice)
		if err != nil {
			log.Printf("Error scanning branch office: %v", err)
			continue
		}
		branchOffices = append(branchOffices, bo)
	}

	// Query programs
	programQuery := `SELECT id, name FROM program ORDER BY name`
	programRows, err := db.DB.Query(programQuery)
	if err != nil {
		log.Printf("Error querying programs: %v", err)
		response := models.GetBranchOfficesResponse{
			Success:       true,
			Message:       "Retrieved branch offices but failed to get programs",
			BranchOffices: branchOffices,
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	defer programRows.Close()

	var programs []models.Program
	for programRows.Next() {
		var program models.Program
		err := programRows.Scan(&program.ID, &program.Name)
		if err != nil {
			log.Printf("Error scanning program: %v", err)
			continue
		}
		programs = append(programs, program)
	}

	response := models.GetBranchOfficesResponse{
		Success:       true,
		Message:       "Branch offices and programs retrieved successfully",
		BranchOffices: branchOffices,
		Programs:      programs,
	}
	json.NewEncoder(w).Encode(response)
}

// DeleteAllBranchOffices ลบข้อมูลสาขาทั้งหมด
func DeleteAllBranchOfficesHandler(w http.ResponseWriter, r *http.Request) {
	database := db.DB

	// ลบข้อมูลทั้งหมดจากตาราง branch_office และรีเซ็ต auto-increment
	// ใช้ DELETE FROM แทน TRUNCATE TABLE เพื่อหลีกเลี่ยง foreign key constraints
	_, err := database.Exec("DELETE FROM branch_office")
	if err != nil {
		log.Printf("Error deleting all branch offices: %v", err)
		response := models.DeleteAllBranchOfficesResponse{
			Success: false,
			Message: "Failed to delete all branch offices: " + err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// รีเซ็ต auto-increment counter
	_, err = database.Exec("ALTER TABLE branch_office AUTO_INCREMENT = 1")
	if err != nil {
		log.Printf("Error resetting auto increment: %v", err)
		// ไม่ return error เพราะข้อมูลถูกลบแล้ว
	}

	response := models.DeleteAllBranchOfficesResponse{
		Success: true,
		Message: "All branch offices deleted successfully",
	}
	json.NewEncoder(w).Encode(response)
}
