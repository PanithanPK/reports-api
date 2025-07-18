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

// ListBranchesHandler returns a handler for listing all branches
func ListBranchesHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query(`SELECT id, name, created_at, updated_at, deleted_at, created_by, updated_by, deleted_by FROM branches WHERE deleted_at IS NULL`)
	if err != nil {
		http.Error(w, "Failed to query branches", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var branches []models.Branch
	for rows.Next() {
		var b models.Branch
		err := rows.Scan(&b.ID, &b.Name, &b.CreatedAt, &b.UpdatedAt, &b.DeletedAt, &b.CreatedBy, &b.UpdatedBy, &b.DeletedBy)
		if err != nil {
			log.Printf("Error scanning branch: %v", err)
			continue
		}
		branches = append(branches, b)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "data": branches})
}

// CreateBranchHandler returns a handler for creating a new branch
func CreateBranchHandler(w http.ResponseWriter, r *http.Request) {
	var req models.BranchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	res, err := db.DB.Exec(`INSERT INTO branches (name, created_by, updated_by) VALUES (?, ?, ?)`, req.Name, req.CreatedBy, req.UpdatedBy)
	if err != nil {
		http.Error(w, "Failed to insert branch", http.StatusInternalServerError)
		return
	}
	id, _ := res.LastInsertId()
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "id": id})
}

// UpdateBranchHandler returns a handler for updating an existing branch
func UpdateBranchHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}
	var req models.BranchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	_, err = db.DB.Exec(`UPDATE branches SET name=?, updated_by=?, updated_at=CURRENT_TIMESTAMP WHERE id=? AND deleted_at IS NULL`, req.Name, req.UpdatedBy, id)
	if err != nil {
		http.Error(w, "Failed to update branch", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}

// DeleteBranchHandler returns a handler for deleting a branch
func DeleteBranchHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}
	_, err = db.DB.Exec(`DELETE FROM branches WHERE id=?`, id)
	if err != nil {
		http.Error(w, "Failed to delete branch", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}

// GetBranchDetailHandler returns detailed information about a specific branch
func GetBranchDetailHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	// 1. Get branch information
	var branchDetail models.BranchDetail
	err = db.DB.QueryRow(`
		SELECT id, name, created_at, updated_at 
		FROM branches 
		WHERE id = ? AND deleted_at IS NULL
	`, id).Scan(&branchDetail.ID, &branchDetail.Name, &branchDetail.CreatedAt, &branchDetail.UpdatedAt)

	if err != nil {
		log.Printf("Error fetching branch details: %v", err)
		http.Error(w, "Branch not found", http.StatusNotFound)
		return
	}

	// 2. Count departments in this branch
	err = db.DB.QueryRow(`
		SELECT COUNT(*) 
		FROM departments 
		WHERE branch_id = ? AND deleted_at IS NULL
	`, id).Scan(&branchDetail.DepartmentsCount)

	if err != nil {
		log.Printf("Error counting departments: %v", err)
		branchDetail.DepartmentsCount = 0 // Default to 0 if error
	}

	// 3. Count IP phones in this branch
	err = db.DB.QueryRow(`
		SELECT COUNT(*) FROM ip_phones ip
		JOIN departments d ON ip.department_id = d.id
		WHERE d.branch_id = ? AND ip.deleted_at IS NULL
	`, id).Scan(&branchDetail.IPPhonesCount)

	if err != nil {
		log.Printf("Error counting IP phones: %v", err)
		branchDetail.IPPhonesCount = 0 // Default to 0 if error
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    branchDetail,
	})
}
