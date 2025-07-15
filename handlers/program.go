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

// แสดงรายการโปรแกรมทั้งหมด
func ListProgramsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query(`SELECT id, name, created_at, updated_at, deleted_at, created_by, updated_by, deleted_by FROM systems_program WHERE deleted_at IS NULL`)
	if err != nil {
		http.Error(w, "Failed to query programs", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var programs []models.Program
	for rows.Next() {
		var p models.Program
		err := rows.Scan(&p.ID, &p.Name, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt, &p.CreatedBy, &p.UpdatedBy, &p.DeletedBy)
		if err != nil {
			log.Printf("Error scanning program: %v", err)
			continue
		}
		programs = append(programs, p)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "programs": programs})
}

// เพิ่มโปรแกรม
func CreateProgramHandler(w http.ResponseWriter, r *http.Request) {
	var req models.ProgramRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	res, err := db.DB.Exec(`INSERT INTO systems_program (name, created_by, updated_by) VALUES (?, ?, ?)`, req.Name, req.CreatedBy, req.UpdatedBy)
	if err != nil {
		http.Error(w, "Failed to insert program", http.StatusInternalServerError)
		return
	}
	id, _ := res.LastInsertId()
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "id": id})
}

// แก้ไขโปรแกรม
func UpdateProgramHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}
	var req models.ProgramRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	_, err = db.DB.Exec(`UPDATE systems_program SET name=?, updated_by=?, updated_at=? WHERE id=? AND deleted_at IS NULL`, req.Name, req.UpdatedBy, time.Now(), id)
	if err != nil {
		http.Error(w, "Failed to update program", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}

// ลบโปรแกรม (soft delete)
func DeleteProgramHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}
	deletedByStr := r.URL.Query().Get("deleted_by")
	deletedBy, _ := strconv.Atoi(deletedByStr)
	_, err = db.DB.Exec(`UPDATE systems_program SET deleted_at=?, deleted_by=? WHERE id=? AND deleted_at IS NULL`, time.Now(), deletedBy, id)
	if err != nil {
		http.Error(w, "Failed to delete program", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}
