package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"reports-api/db"
	"reports-api/models"
	"strconv"
	"time"
)

// ListDepartmentsHandler returns a handler for listing all departments
func ListDepartmentsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query(`SELECT id, name, branch_id, created_at, updated_at, deleted_at FROM departments WHERE deleted_at IS NULL`)
	if err != nil {
		http.Error(w, "Failed to query departments", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var departments []models.Department
	for rows.Next() {
		var d models.Department
		err := rows.Scan(&d.ID, &d.Name, &d.BranchID, &d.CreatedAt, &d.UpdatedAt, &d.DeletedAt)
		if err != nil {
			log.Printf("Error scanning department: %v", err)
			continue
		}
		departments = append(departments, d)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "departments": departments})
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
	idStr := r.URL.Query().Get("id")
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
	_, err = db.DB.Exec(`UPDATE departments SET name=?, branch_id=?, updated_at=? WHERE id=? AND deleted_at IS NULL`, req.Name, req.BranchID, time.Now(), id)
	if err != nil {
		http.Error(w, "Failed to update department", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}

// DeleteDepartmentHandler returns a handler for deleting a department
func DeleteDepartmentHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}
	_, err = db.DB.Exec(`UPDATE departments SET deleted_at=? WHERE id=? AND deleted_at IS NULL`, time.Now(), id)
	if err != nil {
		http.Error(w, "Failed to delete department", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}
