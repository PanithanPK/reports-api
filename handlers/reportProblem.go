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

// GetTasksHandler returns a handler for listing all tasks with details
func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT t.id, t.phone_id, COALESCE(p.number, 0), COALESCE(p.name, ''), t.system_id, COALESCE(s.name, ''),
		COALESCE(p.department_id, 0), COALESCE(d.name, ''), COALESCE(d.branch_id, 0), COALESCE(b.name, ''),
		t.text, t.status, t.created_at, t.updated_at
		FROM tasks t
		LEFT JOIN ip_phones p ON t.phone_id = p.id
		LEFT JOIN departments d ON p.department_id = d.id
		LEFT JOIN branches b ON d.branch_id = b.id
		LEFT JOIN systems_program s ON t.system_id = s.id
	`
	rows, err := db.DB.Query(query)
	if err != nil {
		log.Printf("Error querying tasks with join: %v", err)
		http.Error(w, "Failed to query tasks", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []models.TaskWithDetails
	for rows.Next() {
		var t models.TaskWithDetails
		err := rows.Scan(&t.ID, &t.PhoneID, &t.Number, &t.PhoneName, &t.SystemID, &t.SystemName, &t.DepartmentID, &t.DepartmentName, &t.BranchID, &t.BranchName, &t.Text, &t.Status, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			log.Printf("Error scanning task: %v", err)
			continue
		}
		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Row error: %v", err)
		http.Error(w, "Failed to read tasks", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    tasks,
	})
}

// CreateTaskHandler เพิ่ม task ใหม่
func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var req models.TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	res, err := db.DB.Exec(`INSERT INTO tasks (phone_id, system_id, text, status, created_by) VALUES (?, ?, ?, 0, ?)`, req.PhoneID, req.SystemID, req.Text, req.CreatedBy)
	if err != nil {
		http.Error(w, "Failed to insert task", http.StatusInternalServerError)
		return
	}
	id, _ := res.LastInsertId()
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "id": id})
}

// UpdateTaskHandler แก้ไข task
func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}
	var req models.TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	_, err = db.DB.Exec(`UPDATE tasks SET phone_id=?, system_id=?, text=?, status=?, updated_at=CURRENT_TIMESTAMP, updated_by=? WHERE id=?`, req.PhoneID, req.SystemID, req.Text, req.Status, req.UpdatedBy, id)
	if err != nil {
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}

// DeleteTaskHandler (soft delete)
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}
	_, err = db.DB.Exec(`DELETE FROM tasks WHERE id=?`, id)
	if err != nil {
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}

func GetTaskDetailHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	var task models.TaskWithDetails
	err = db.DB.QueryRow(`
		SELECT t.id, t.phone_id, COALESCE(p.number, 0), COALESCE(p.name, ''), t.system_id, COALESCE(s.name, ''),
		COALESCE(p.department_id, 0), COALESCE(d.name, ''), COALESCE(d.branch_id, 0), COALESCE(b.name, ''),
		t.text, t.status, t.created_at, t.updated_at
		FROM tasks t
		LEFT JOIN ip_phones p ON t.phone_id = p.id
		LEFT JOIN departments d ON p.department_id = d.id
		LEFT JOIN branches b ON d.branch_id = b.id
		LEFT JOIN systems_program s ON t.system_id = s.id
		WHERE t.id = ?
	`, id).Scan(&task.ID, &task.PhoneID, &task.Number, &task.PhoneName, &task.SystemID, &task.SystemName, &task.DepartmentID, &task.DepartmentName, &task.BranchID, &task.BranchName, &task.Text, &task.Status, &task.CreatedAt, &task.UpdatedAt)

	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    task,
	})
}
