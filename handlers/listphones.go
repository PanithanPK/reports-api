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

// ListIPPhonesHandler returns a handler for listing all IP phones
func ListIPPhonesHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query(`SELECT id, number, name, department_id, created_at, updated_at, deleted_at, created_by, updated_by, deleted_by FROM ip_phones WHERE deleted_at IS NULL`)
	if err != nil {
		http.Error(w, "Failed to query ip_phones", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var phones []models.IPPhone
	for rows.Next() {
		var p models.IPPhone
		err := rows.Scan(&p.ID, &p.Number, &p.Name, &p.DepartmentID, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt, &p.CreatedBy, &p.UpdatedBy, &p.DeletedBy)
		if err != nil {
			log.Printf("Error scanning ip_phone: %v", err)
			continue
		}
		phones = append(phones, p)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "ip_phones": phones})
}

// CreateIPPhoneHandler returns a handler for creating a new IP phone
func CreateIPPhoneHandler(w http.ResponseWriter, r *http.Request) {
	var req models.IPPhoneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	res, err := db.DB.Exec(`INSERT INTO ip_phones (number, name, department_id, created_by, updated_by) VALUES (?, ?, ?, ?, ?)`, req.Number, req.Name, req.DepartmentID, req.CreatedBy, req.UpdatedBy)
	if err != nil {
		http.Error(w, "Failed to insert ip_phone", http.StatusInternalServerError)
		return
	}
	id, _ := res.LastInsertId()
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "id": id})
}

// UpdateIPPhoneHandler returns a handler for updating an existing IP phone
func UpdateIPPhoneHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}
	var req models.IPPhoneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	_, err = db.DB.Exec(`UPDATE ip_phones SET number=?, name=?, department_id=?, updated_by=?, updated_at=? WHERE id=? AND deleted_at IS NULL`, req.Number, req.Name, req.DepartmentID, req.UpdatedBy, time.Now(), id)
	if err != nil {
		http.Error(w, "Failed to update ip_phone", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}

// DeleteIPPhoneHandler returns a handler for deleting an IP phone
func DeleteIPPhoneHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}
	deletedByStr := r.URL.Query().Get("deleted_by")
	deletedBy, _ := strconv.Atoi(deletedByStr)
	_, err = db.DB.Exec(`UPDATE ip_phones SET deleted_at=?, deleted_by=? WHERE id=? AND deleted_at IS NULL`, time.Now(), deletedBy, id)
	if err != nil {
		http.Error(w, "Failed to delete ip_phone", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}
