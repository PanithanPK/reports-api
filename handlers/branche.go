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
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "branches": branches})
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
	idStr := r.URL.Query().Get("id")
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
	_, err = db.DB.Exec(`UPDATE branches SET name=?, updated_by=?, updated_at=? WHERE id=? AND deleted_at IS NULL`, req.Name, req.UpdatedBy, time.Now(), id)
	if err != nil {
		http.Error(w, "Failed to update branch", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}

// DeleteBranchHandler returns a handler for deleting a branch
func DeleteBranchHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}
	deletedByStr := r.URL.Query().Get("deleted_by")
	deletedBy, _ := strconv.Atoi(deletedByStr)
	_, err = db.DB.Exec(`UPDATE branches SET deleted_at=?, deleted_by=? WHERE id=? AND deleted_at IS NULL`, time.Now(), deletedBy, id)
	if err != nil {
		http.Error(w, "Failed to delete branch", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}
