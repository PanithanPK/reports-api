package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"reports-api/db"
	"reports-api/models"
)

// GetTasksHandler คืนค่าข้อมูล tasks ทั้งหมด
func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := `
		SELECT t.id, t.phone_id, COALESCE(p.name, ''), t.system_id, COALESCE(s.name, ''),
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
		err := rows.Scan(&t.ID, &t.PhoneID, &t.PhoneName, &t.SystemID, &t.SystemName, &t.DepartmentID, &t.DepartmentName, &t.BranchID, &t.BranchName, &t.Text, &t.Status, &t.CreatedAt, &t.UpdatedAt)
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
		"tasks":   tasks,
	})
}
