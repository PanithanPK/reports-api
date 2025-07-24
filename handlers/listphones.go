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

// ListIPPhonesHandler returns a handler for listing all IP phones
func ListIPPhonesHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query(`
   SELECT ip.id, ip.number, ip.name, ip.department_id,
          d.name as department_name, d.branch_id, b.name as branch_name,
          ip.created_at, ip.updated_at, ip.deleted_at, ip.created_by, ip.updated_by, ip.deleted_by
   FROM ip_phones ip
   LEFT JOIN departments d ON ip.department_id = d.id
   LEFT JOIN branches b ON d.branch_id = b.id
   WHERE ip.deleted_at IS NULL
 `)
	if err != nil {
		http.Error(w, "Failed to query ip_phones", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var phones []models.IPPhone
	for rows.Next() {
		var p models.IPPhone
		err := rows.Scan(
			&p.ID, &p.Number, &p.Name, &p.DepartmentID,
			&p.DepartmentName, &p.BranchID, &p.BranchName,
			&p.CreatedAt, &p.UpdatedAt, &p.DeletedAt, &p.CreatedBy, &p.UpdatedBy, &p.DeletedBy,
		)
		if err != nil {
			log.Printf("Error scanning ip_phone: %v", err)
			continue
		}
		phones = append(phones, p)
	}
	log.Printf("Getting IP phones Success")
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "data": phones})
}

// CreateIPPhoneHandler returns a handler for creating a new IP phone
func CreateIPPhoneHandler(w http.ResponseWriter, r *http.Request) {
	var req models.IPPhoneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	res, err := db.DB.Exec(`INSERT INTO ip_phones (number, name, department_id, created_by) VALUES (?, ?, ?, ?)`, req.Number, req.Name, req.DepartmentID, req.CreatedBy)
	if err != nil {
		http.Error(w, "Failed to insert ip_phone", http.StatusInternalServerError)
		return
	}
	log.Printf("Inserted new IP phone: %d", req.Number)
	id, _ := res.LastInsertId()
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "id": id})
}

// UpdateIPPhoneHandler returns a handler for updating an existing IP phone
func UpdateIPPhoneHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
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
	_, err = db.DB.Exec(`UPDATE ip_phones SET number=?, name=?, department_id=?, updated_by=?, updated_at=CURRENT_TIMESTAMP WHERE id=? AND deleted_at IS NULL`, req.Number, req.Name, req.DepartmentID, req.UpdatedBy, id)
	if err != nil {
		http.Error(w, "Failed to update ip_phone", http.StatusInternalServerError)
		return
	}
	log.Printf("Updating IP phone ID: %d with number: %d", id, req.Number)
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}

// DeleteIPPhoneHandler returns a handler for deleting an IP phone
func DeleteIPPhoneHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}
	_, err = db.DB.Exec(`DELETE FROM ip_phones WHERE id=?`, id)
	if err != nil {
		http.Error(w, "Failed to delete ip_phone", http.StatusInternalServerError)
		return
	}
	log.Printf("Deleted IP phone ID: %d", id)
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}
