package handlers

import (
	"log"
	"reports-api/db"
	"reports-api/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// GetTasksHandler returns a handler for listing all tasks with details
func GetTasksHandler(c *fiber.Ctx) error {
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
		return c.Status(500).JSON(fiber.Map{"error": "Failed to query tasks"})
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
		return c.Status(500).JSON(fiber.Map{"error": "Failed to read tasks"})
	}
	log.Printf("Getting tasks Success")

	return c.JSON(fiber.Map{"success": true, "data": tasks})
}

// CreateTaskHandler เพิ่ม task ใหม่
func CreateTaskHandler(c *fiber.Ctx) error {
	var req models.TaskRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Validate department_id or phone_id
	if req.DepartmentID == 0 {
		if req.PhoneID == nil {
			return c.Status(400).JSON(fiber.Map{"error": "Either department_id or phone_id is required"})
		}
		// Get department_id from phone_id
		err := db.DB.QueryRow("SELECT department_id FROM ip_phones WHERE id = ?", *req.PhoneID).Scan(&req.DepartmentID)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid phone_id"})
		}
	}

	res, err := db.DB.Exec(`INSERT INTO tasks (phone_id, system_id, department_id, text, status, created_by) VALUES (?, ?, ?, ?, 0, ?)`, req.PhoneID, req.SystemID, req.DepartmentID, req.Text, req.CreatedBy)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to insert task"})
	}
	id, _ := res.LastInsertId()

	// Update department score
	updateDepartmentScore(req.DepartmentID)
	
	log.Printf("Inserted new task with ID: %d", id)
	if req.Telegram == true {
		// Get additional data for Telegram
		var phoneNumber int
		var departmentName, branchName string
		
		if req.PhoneID != nil {
			// Get data from phone if phone_id exists
			db.DB.QueryRow(`
				SELECT p.number, d.name, b.name 
				FROM ip_phones p 
				JOIN departments d ON p.department_id = d.id 
				JOIN branches b ON d.branch_id = b.id 
				WHERE p.id = ?
			`, *req.PhoneID).Scan(&phoneNumber, &departmentName, &branchName)
		} else {
			// Get data from department_id if no phone_id
			db.DB.QueryRow(`
				SELECT d.name, b.name 
				FROM departments d 
				JOIN branches b ON d.branch_id = b.id 
				WHERE d.id = ?
			`, req.DepartmentID).Scan(&departmentName, &branchName)
		}
		
		req.PhoneNumber = phoneNumber
		req.DepartmentName = departmentName
		req.BranchName = branchName
		_ = SendTelegram(req)
	}
	return c.JSON(fiber.Map{"success": true, "id": id})
}

// UpdateTaskHandler แก้ไข task
func UpdateTaskHandler(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid id"})
	}

	var req models.TaskRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	_, err = db.DB.Exec(`UPDATE tasks SET phone_id=?, system_id=?, department_id=?, text=?, status=?, updated_at=CURRENT_TIMESTAMP, updated_by=? WHERE id=?`, req.PhoneID, req.SystemID, req.DepartmentID, req.Text, req.Status, req.UpdatedBy, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update task"})
	}

	log.Printf("Updating task ID: %d", id)
	return c.JSON(fiber.Map{"success": true})
}

func UpdateTaskStatusHandler(c *fiber.Ctx) error {
	var req models.TaskStatusUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	_, err := db.DB.Exec(`UPDATE tasks SET status=?, updated_at=CURRENT_TIMESTAMP, updated_by=? WHERE id=?`, req.Status, req.UpdatedBy, req.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update task status"})
	}

	log.Printf("Updating task ID: %d status to: %d", req.ID, req.Status)
	return c.JSON(fiber.Map{"success": true})
}

// DeleteTaskHandler (soft delete)
func DeleteTaskHandler(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid id"})
	}

	_, err = db.DB.Exec(`DELETE FROM tasks WHERE id=?`, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete task"})
	}

	log.Printf("Deleted task ID: %d", id)
	return c.JSON(fiber.Map{"success": true})
}

func GetTaskDetailHandler(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid id"})
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
		return c.Status(404).JSON(fiber.Map{"error": "Task not found"})
	}

	log.Printf("Getting task ID: %d details", id)
	return c.JSON(fiber.Map{"success": true, "data": task})
}

// updateDepartmentScore updates the department score based on problem count
func updateDepartmentScore(departmentID int) {
	now := time.Now()
	year, month := now.Year(), int(now.Month())

	// 1. Check if record exists for this department/month
	var exists bool
	err := db.DB.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM scores WHERE department_id = ? AND year = ? AND month = ?)
	`, departmentID, year, month).Scan(&exists)
	if err != nil {
		log.Printf("Error checking score record: %v", err)
		return
	}

	// Insert new record if it doesn't exist
	if !exists {
		_, err := db.DB.Exec(`
			INSERT INTO scores (department_id, year, month, score)
			VALUES (?, ?, ?, 100)
		`, departmentID, year, month)
		if err != nil {
			log.Printf("Error creating score record: %v", err)
			return
		}
	}

	// 2. Check number of problems in that month
	var problemCount int
	err = db.DB.QueryRow(`
		SELECT COUNT(*) FROM tasks t
		WHERE t.department_id = ? AND YEAR(t.created_at) = ? AND MONTH(t.created_at) = ?
	`, departmentID, year, month).Scan(&problemCount)
	if err != nil {
		log.Printf("Error counting problems: %v", err)
		return
	}
	log.Printf("Department %d has %d problems in %d/%d", departmentID, problemCount, month, year)
	// 3. If problem count > 3, deduct score
	if problemCount > 3 {
		_, err := db.DB.Exec(`
			UPDATE scores
			SET score = GREATEST(score - 1, 0)
			WHERE department_id = ? AND year = ? AND month = ?
		`, departmentID, year, month)
		if err != nil {
			log.Printf("Error updating score: %v", err)
		}
		log.Printf("Updated department %d score for month %d/%d, problem count: %d", departmentID, month, year, problemCount)
	}
}
