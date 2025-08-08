package handlers

import (
	"fmt"
	"log"
	"reports-api/db"
	"reports-api/models"
	"reports-api/utils"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func generateticketno(no int) string {
	ticket := fmt.Sprintf("ticket-"+"%04d", no)
	return ticket
}

// GetTasksHandler returns a handler for listing all tasks with details and pagination
func GetTasksHandler(c *fiber.Ctx) error {
	pagination := utils.GetPaginationParams(c)
	offset := utils.CalculateOffset(pagination.Page, pagination.Limit)

	// Get total count
	var total int
	err := db.DB.QueryRow(`SELECT COUNT(*) FROM tasks`).Scan(&total)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to count tasks"})
	}

	// Get paginated data
	query := `
		SELECT t.id, IFNULL(t.ticket_no, ''), IFNULL(t.phone_id, 0) as phone_id, IFNULL(p.number, 0) as number , IFNULL(p.name, '') as phone_name, t.system_id, IFNULL(s.name, '') as system_name,
		IFNULL(t.department_id, 0), IFNULL(d.name, '') as department_name, IFNULL(d.branch_id, 0) as branch_id, IFNULL(b.name, '') as branch_name,
		t.text, t.status, t.created_at, t.updated_at
		FROM tasks t
		LEFT JOIN ip_phones p ON t.phone_id = p.id
		LEFT JOIN departments d ON t.department_id = d.id
		LEFT JOIN branches b ON d.branch_id = b.id
		LEFT JOIN systems_program s ON t.system_id = s.id
		ORDER BY t.id DESC
		LIMIT ? OFFSET ?
	`
	rows, err := db.DB.Query(query, pagination.Limit, offset)
	if err != nil {
		log.Printf("Error querying tasks with join: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to query tasks"})
	}
	defer rows.Close()

	var tasks []models.TaskWithDetails
	for rows.Next() {
		var t models.TaskWithDetails
		err := rows.Scan(&t.ID, &t.Ticket, &t.PhoneID, &t.Number, &t.PhoneName, &t.SystemID, &t.SystemName, &t.DepartmentID, &t.DepartmentName, &t.BranchID, &t.BranchName, &t.Text, &t.Status, &t.CreatedAt, &t.UpdatedAt)

		if err != nil {
			log.Printf("Error scanning task: %v", err)
			continue
		}

		// Calculate overdue
		createdAt, err := time.Parse(time.RFC3339, t.CreatedAt)
		if err == nil {
			createdAt = createdAt.Add(7 * time.Hour)
			now := time.Now().Add(7 * time.Hour)
			duration := now.Sub(createdAt)
			if createdAt.Format("2006-01-02") == now.Format("2006-01-02") {
				hours := int(duration.Hours())
				minutes := int(duration.Minutes()) % 60
				seconds := int(duration.Seconds()) % 60
				t.Overdue = fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
			} else {
				days := int(duration.Hours() / 24)
				if days == 0 {
					days = 1
				}
				t.Overdue = days
			}
		}

		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Row error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to read tasks"})
	}

	log.Printf("Getting tasks Success")
	return c.JSON(models.PaginatedResponse{
		Success: true,
		Data:    tasks,
		Pagination: models.PaginationResponse{
			Page:       pagination.Page,
			Limit:      pagination.Limit,
			Total:      total,
			TotalPages: utils.CalculateTotalPages(total, pagination.Limit),
		},
	})
}

// CreateTaskHandler เพิ่ม task ใหม่
func CreateTaskHandler(c *fiber.Ctx) error {
	var req models.TaskRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Get department_id from phone_id if phone_id exists
	if req.PhoneID != nil {
		err := db.DB.QueryRow("SELECT department_id FROM ip_phones WHERE id = ?", *req.PhoneID).Scan(&req.DepartmentID)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid phone_id"})
		}
	}

	// Get latest ID and add 1 for ticket number
	var lastID int
	err := db.DB.QueryRow("SELECT COALESCE(MAX(id), 0) FROM tasks").Scan(&lastID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get last ID"})
	}
	nextID := lastID + 1
	ticketno := generateticketno(nextID)

	res, err := db.DB.Exec(`INSERT INTO tasks (phone_id, ticket_no, system_id, department_id, text, status, created_by) VALUES (?, ?, ?, ?, ?, 0, ?)`, req.PhoneID, ticketno, req.SystemID, req.DepartmentID, req.Text, req.CreatedBy)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to insert task"})
	}
	id, _ := res.LastInsertId()

	// Update department score
	if err := updateDepartmentScore(req.DepartmentID); err != nil {
		log.Printf("Failed to update department score: %v", err)
	}

	log.Printf("Inserted new task with ID: %d", id)
	if req.Telegram == true {
		// Get additional data for Telegram
		var phoneNumber int
		var departmentName, branchName string
		var programName string

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

		if req.SystemID != 0 {
			db.DB.QueryRow(`
				SELECT name
				FROM systems_program
				WHERE id = ?
			`, req.SystemID).Scan(&programName)
		}

		req.PhoneNumber = phoneNumber
		req.DepartmentName = departmentName
		req.BranchName = branchName
		req.ProgramName = programName
		req.Url = "http://helpdesk.nopadol.com/"
		req.CreatedAt = time.Now().Add(7 * time.Hour).Format("2006-01-02 15:04:05")
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

	var req models.TaskRequestUpdate
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Handle phone_id = 0 as null
	if req.PhoneID != nil && *req.PhoneID == 0 {
		req.PhoneID = nil
	}

	// Get department_id from phone_id if phone_id exists and is valid
	if req.PhoneID != nil {
		err := db.DB.QueryRow("SELECT department_id FROM ip_phones WHERE id = ?", *req.PhoneID).Scan(&req.DepartmentID)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid phone_id"})
		}
	} else if req.DepartmentID == 0 {
		// If no phone_id and no department_id provided, get from existing task
		err := db.DB.QueryRow("SELECT department_id FROM tasks WHERE id = ?", id).Scan(&req.DepartmentID)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Task not found"})
		}
	}

	_, err = db.DB.Exec(`UPDATE tasks SET phone_id=?, system_id=?, department_id=?, assignto=?, text=?, status=?, updated_at=CURRENT_TIMESTAMP, updated_by=? WHERE id=?`, req.PhoneID, req.SystemID, req.DepartmentID, req.Assignto, req.Text, req.Status, req.UpdatedBy, id)
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
		t.department_id, COALESCE(d.name, ''), COALESCE(d.branch_id, 0), COALESCE(b.name, ''),
		t.text, t.status, t.created_at, t.updated_at
		FROM tasks t
		LEFT JOIN ip_phones p ON t.phone_id = p.id
		LEFT JOIN departments d ON t.department_id = d.id
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
func updateDepartmentScore(departmentID int) error {
	now := time.Now()
	year, month := now.Year(), int(now.Month())

	// 1. Check if record exists for this department/month
	var exists bool
	err := db.DB.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM scores WHERE department_id = ? AND year = ? AND month = ?)
	`, departmentID, year, month).Scan(&exists)
	if err != nil {
		log.Printf("Error checking score record: %v", err)
		return err
	}

	// Insert new record if it doesn't exist
	if !exists {
		_, err := db.DB.Exec(`
			INSERT INTO scores (department_id, year, month, score)
			VALUES (?, ?, ?, 100)
		`, departmentID, year, month)
		if err != nil {
			log.Printf("Error creating score record: %v", err)
			return err
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
		return err
	}

	log.Printf("Department ID %d has %d problems in %02d/%d", departmentID, problemCount, month, year)

	// 3. If problem count > 3, deduct score
	if problemCount > 3 {
		_, err := db.DB.Exec(`
			UPDATE scores
			SET score = GREATEST(score - 1, 0)
			WHERE department_id = ? AND year = ? AND month = ?
		`, departmentID, year, month)
		if err != nil {
			log.Printf("Error updating score: %v", err)
			return err
		}
		log.Printf("Updated department ID %d score for month %02d/%d, problem count: %d", departmentID, month, year, problemCount)
	}

	return nil
}
