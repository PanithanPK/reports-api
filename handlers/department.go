package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"reports-api/db"
	"reports-api/models"
	"strconv"

	"github.com/gorilla/mux"
)

// ListDepartmentsHandler returns a handler for listing all departments
func ListDepartmentsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query(`
   SELECT d.id, d.name, d.branch_id, b.name as branch_name, d.created_at, d.updated_at, d.deleted_at
   FROM departments d
   LEFT JOIN branches b ON d.branch_id = b.id
   WHERE d.deleted_at IS NULL
 `)
	if err != nil {
		http.Error(w, "Failed to query departments", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var departments []models.Department
	for rows.Next() {
		var d models.Department
		err := rows.Scan(&d.ID, &d.Name, &d.BranchID, &d.BranchName, &d.CreatedAt, &d.UpdatedAt, &d.DeletedAt)
		if err != nil {
			log.Printf("Error scanning department: %v", err)
			continue
		}
		departments = append(departments, d)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "data": departments})
}

// CreateDepartmentHandler returns a handler for creating a new department
func CreateDepartmentHandler(w http.ResponseWriter, r *http.Request) {
	var req models.DepartmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	res, err := db.DB.Exec(`INSERT INTO departments (name, branch_id) VALUES (?, ?)`, req.Name, req.BranchID)
	if err != nil {
		http.Error(w, "Failed to insert department", http.StatusInternalServerError)
		return
	}
	id, _ := res.LastInsertId()
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "id": id})
}

// UpdateDepartmentHandler returns a handler for updating an existing department
func UpdateDepartmentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}
	var req models.DepartmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	_, err = db.DB.Exec(`UPDATE departments SET name=?, branch_id=?, updated_at=CURRENT_TIMESTAMP WHERE id=? AND deleted_at IS NULL`, req.Name, req.BranchID, id)
	if err != nil {
		http.Error(w, "Failed to update department", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}

// DeleteDepartmentHandler returns a handler for deleting a department
func DeleteDepartmentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}
	_, err = db.DB.Exec(`DELETE FROM departments WHERE id=?`, id)
	if err != nil {
		http.Error(w, "Failed to delete department", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}

// GetDepartmentDetailHandler returns detailed information about a specific department
func GetDepartmentDetailHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	// 1. Get department information
	var departmentDetail models.DepartmentDetail
	err = db.DB.QueryRow(`
		SELECT d.id, d.name, d.branch_id, b.name, d.created_at, d.updated_at 
		FROM departments d
		LEFT JOIN branches b ON d.branch_id = b.id
		WHERE d.id = ? AND d.deleted_at IS NULL
	`, id).Scan(
		&departmentDetail.ID,
		&departmentDetail.Name,
		&departmentDetail.BranchID,
		&departmentDetail.BranchName,
		&departmentDetail.CreatedAt,
		&departmentDetail.UpdatedAt,
	)

	if err != nil {
		log.Printf("Error fetching department details: %v", err)
		http.Error(w, "Department not found", http.StatusNotFound)
		return
	}

	// 2. Count IP phones in this department
	err = db.DB.QueryRow(`
		SELECT COUNT(*) 
		FROM ip_phones 
		WHERE department_id = ? AND deleted_at IS NULL
	`, id).Scan(&departmentDetail.IPPhonesCount)

	if err != nil {
		log.Printf("Error counting IP phones: %v", err)
	}

	// 3. Count tasks related to this department
	err = db.DB.QueryRow(`
		SELECT COUNT(*) FROM tasks t
		JOIN ip_phones ip ON t.phone_id = ip.id
		WHERE ip.department_id = ? AND t.deleted_at IS NULL
	`, id).Scan(&departmentDetail.TasksCount)

	if err != nil {
		log.Printf("Error counting tasks: %v", err)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    departmentDetail,
	})
}
