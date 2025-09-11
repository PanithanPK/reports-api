package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/url"
	"os"
	"reports-api/config"
	"reports-api/db"
	"reports-api/handlers/common"
	"reports-api/models"
	"reports-api/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// GetTasksHandler returns a handler for listing all tasks with details and pagination
// @Summary Get all problems
// @Description Get list of all problems with pagination
// @Tags problems
// @Accept json
// @Produce json
// @Success 200 {object} models.PaginatedResponse
// @Router /api/v1/problem/list [get]
func GetTasksHandler(c *fiber.Ctx) error {
	// Get pagination params
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
		SELECT t.id, IFNULL(t.ticket_no, ''), IFNULL(t.phone_id, 0), IFNULL(p.number, 0), IFNULL(p.name, ''), t.system_id, IFNULL(s.name, ''), IFNULL(t.issue_type, 0), IFNULL(t.issue_else, ''), IFNULL(it.name, ''), IFNULL(t.department_id, 0), IFNULL(d.name, ''), IFNULL(d.branch_id, 0), IFNULL(b.name, ''), t.text, IFNULL(t.assignto_id, 0), IFNULL(t.assignto, ''), IFNULL(t.reported_by, ''), t.status, t.created_at, t.updated_at, IFNULL(t.file_paths, '[]')
		FROM tasks t
		LEFT JOIN ip_phones p ON t.phone_id = p.id
		LEFT JOIN departments d ON t.department_id = d.id
		LEFT JOIN branches b ON d.branch_id = b.id
		LEFT JOIN systems_program s ON t.system_id = s.id
		LEFT JOIN issue_types it ON t.issue_type = it.id
		ORDER BY id DESC
		LIMIT ? OFFSET ?
	`
	rows, err := db.DB.Query(query, pagination.Limit, offset)
	if err != nil {
		log.Printf("Error querying tasks with join: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to query tasks"})
	}
	defer rows.Close()

	// task model Mapping
	var tasks []models.TaskWithDetails
	for rows.Next() {
		var t models.TaskWithDetails
		var issueTypeName string
		var filePathsJSON string
		// Scan row into task model
		err := rows.Scan(&t.ID, &t.Ticket, &t.PhoneID, &t.Number, &t.PhoneName, &t.SystemID, &t.SystemName, &t.IssueTypeID, &t.IssueElse, &issueTypeName, &t.DepartmentID, &t.DepartmentName, &t.BranchID, &t.BranchName, &t.Text, &t.AssignedtoID, &t.Assignto, &t.ReportedBy, &t.Status, &t.CreatedAt, &t.UpdatedAt, &filePathsJSON)
		if err != nil {
			log.Printf("Error scanning task: %v", err)
			continue
		}
		// Set SystemType based on SystemID
		if t.SystemID > 0 {
			t.SystemType = issueTypeName
		} else {
			t.SystemType = issueTypeName
		}

		// Parse file_paths JSON
		fileMap := make(map[string]string)
		// Check file Path
		if filePathsJSON != "" && filePathsJSON != "[]" {
			var filePaths []fiber.Map
			// Parse JSON array file paths
			if err := json.Unmarshal([]byte(filePathsJSON), &filePaths); err == nil {
				for i, fp := range filePaths {
					if url, ok := fp["url"].(string); ok {
						fileMap[fmt.Sprintf("image_%d", i)] = url
					}
				}
			}
		}
		// Set FilePaths as map instead of array
		t.FilePaths = fileMap

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
	// Return paginated response
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

// @Summary Get problem details
// @Description Get detailed information of a specific problem
// @Tags problems
// @Accept json
// @Produce json
// @Param id path string true "Problem ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/problem/{id} [get]
func GetTaskDetailHandler(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid id"})
	}
	var filePathsJSON string
	var task models.TaskWithDetails
	var issueTypeName string
	err = db.DB.QueryRow(`
		SELECT t.id, IFNULL(t.ticket_no, ''), IFNULL(t.phone_id, 0), IFNULL(p.number, 0), IFNULL(p.name, ''), IFNULL(t.system_id, 0), IFNULL(s.name, ''), IFNULL(t.issue_type, 0), IFNULL(t.issue_else, ''), IFNULL(it.name, ''), IFNULL(t.department_id, 0), IFNULL(d.name, ''), IFNULL(d.branch_id, 0), IFNULL(b.name, ''), IFNULL(t.text, ''), IFNULL(t.assignto_id, 0), IFNULL(t.assignto, ''), IFNULL(t.reported_by, ''), IFNULL(t.status, 0), IFNULL(t.created_at, ''), IFNULL(t.updated_at, ''), IFNULL(t.file_paths, '[]')
		FROM tasks t
		LEFT JOIN ip_phones p ON t.phone_id = p.id
		LEFT JOIN departments d ON t.department_id = d.id
		LEFT JOIN branches b ON d.branch_id = b.id
		LEFT JOIN systems_program s ON t.system_id = s.id
		LEFT JOIN issue_types it ON t.issue_type = it.id
		WHERE t.id = ?
	`, id).Scan(&task.ID, &task.Ticket, &task.PhoneID, &task.Number, &task.PhoneName, &task.SystemID, &task.SystemName, &task.IssueTypeID, &task.IssueElse, &issueTypeName, &task.DepartmentID, &task.DepartmentName, &task.BranchID, &task.BranchName, &task.Text, &task.AssignedtoID, &task.Assignto, &task.ReportedBy, &task.Status, &task.CreatedAt, &task.UpdatedAt, &filePathsJSON)

	if task.SystemID > 0 {
		task.SystemType = issueTypeName
	} else {
		task.SystemType = issueTypeName
	}

	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Task not found"})
	}

	// Parse file_paths JSON and convert to image_{index} format
	fileMap := make(map[string]string)
	if filePathsJSON != "" && filePathsJSON != "[]" {
		var filePaths []fiber.Map
		if err := json.Unmarshal([]byte(filePathsJSON), &filePaths); err == nil {
			for i, fp := range filePaths {
				if url, ok := fp["url"].(string); ok {
					fileMap[fmt.Sprintf("image_%d", i)] = url
				}
			}
		}
	}
	// Set FilePaths as map instead of array
	task.FilePaths = fileMap

	log.Printf("Getting task ID: %d details", id)
	return c.JSON(fiber.Map{"success": true, "data": task})
}

// CreateTaskHandler เพิ่ม task ใหม่
// @Summary Create new problem
// @Description Create a new problem report
// @Tags problems
// @Accept multipart/form-data
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/problem/create [post]
func CreateTaskHandler(c *fiber.Ctx) error {
	var req models.TaskRequest
	var uploadedFiles []fiber.Map
	// Get latest ID and add 1 for ticket number
	ticketno := common.Generateticketno()

	form, err := c.MultipartForm()
	if err != nil {
		// If multipart parsing fails, try regular body parser
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
		}
	} else {
		// Parse body from multipart form
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
		}

		// Handle file uploads if present (support image_{index} format)
		var allFiles []*multipart.FileHeader

		// Check for indexed files (image_0, image_1, image_2, etc.)
		for key, files := range form.File {
			if strings.HasPrefix(key, "image_") || key == "image" {
				allFiles = append(allFiles, files...)
			}
		}

		uploadedFiles, _ = common.HandleFileUploads(allFiles, ticketno)

		// Convert string form values to int for multipart data
		if phoneIDStr := c.FormValue("phone_id"); phoneIDStr != "" && phoneIDStr != "0" {
			if phoneID, err := strconv.Atoi(phoneIDStr); err == nil {
				req.PhoneID = &phoneID
			}
		}
		if systemIDStr := c.FormValue("system_id"); systemIDStr != "" {
			req.SystemID, _ = strconv.Atoi(systemIDStr)
		}
		if departmentIDStr := c.FormValue("department_id"); departmentIDStr != "" {
			req.DepartmentID, _ = strconv.Atoi(departmentIDStr)
		}
		if issueTypeIDStr := c.FormValue("issue_type"); issueTypeIDStr != "" {
			req.IssueTypeID, _ = strconv.Atoi(issueTypeIDStr)
		}
		if createdByStr := c.FormValue("created_by"); createdByStr != "" {
			req.CreatedBy, _ = strconv.Atoi(createdByStr)
		}

		if reportedByStr := c.FormValue("reported_by"); reportedByStr != "" {
			req.ReportedBy = reportedByStr
		}
		if issueElseStr := c.FormValue("issue_else"); issueElseStr != "" {
			req.IssueElse = issueElseStr
		}
		if telegramStr := c.FormValue("telegram"); telegramStr != "" {
			req.Telegram = telegramStr == "true"
			log.Printf("📥 Parsed telegram field from form: '%s' -> %t", telegramStr, req.Telegram)
		} else {
			log.Printf("📥 No telegram field found in form, defaulting to false")
		}
		if textStr := c.FormValue("text"); textStr != "" {
			req.Text = textStr
		}

	}

	// Get department_id from ip_phones if phone_id is provided
	if req.PhoneID != nil && *req.PhoneID > 0 {
		err := db.DB.QueryRow("SELECT department_id FROM ip_phones WHERE id = ?", *req.PhoneID).Scan(&req.DepartmentID)
		if err != nil {
			log.Printf("Warning: Could not get department_id from phone_id %d: %v", *req.PhoneID, err)
		}
	}

	// log.Printf("Request data: PhoneID=%v, SystemID=%d, DepartmentID=%d, Text=%s, CreatedBy=%d", req.PhoneID, req.SystemID, req.DepartmentID, req.Text, req.CreatedBy)

	if len(uploadedFiles) > 0 {
		log.Printf("Uploaded %d files", len(uploadedFiles))
	}

	var res sql.Result
	if len(uploadedFiles) > 0 {
		filePathsBytes, _ := json.Marshal(uploadedFiles)
		log.Printf("Saving file_paths: %s", string(filePathsBytes))
		if req.SystemID > 0 {
			// ดึง typeid จาก systems_program
			var typeid int
			db.DB.QueryRow(`SELECT type FROM systems_program WHERE id = ?`, req.SystemID).Scan(&typeid)
			res, err = db.DB.Exec(`INSERT INTO tasks (phone_id, ticket_no, system_id, issue_type, department_id, text, reported_by, status, created_by, file_paths) VALUES (?, ?, ?, ?, ?, ?, ?, 0, ?, ?)`, req.PhoneID, ticketno, req.SystemID, typeid, req.DepartmentID, req.Text, req.ReportedBy, req.CreatedBy, string(filePathsBytes))
		} else {
			// กรณีไม่มีระบบ ให้เก็บ issue_type และ issue_else
			if req.IssueTypeID == 0 || req.IssueElse == "" {
				return c.Status(400).JSON(fiber.Map{"error": "issue_type and issue_else must be provided when system_id is 0"})
			}
			// รับค่า issue_type จาก frontend โดยตรง
			res, err = db.DB.Exec(`INSERT INTO tasks (phone_id, ticket_no, system_id, issue_type, issue_else, department_id, text, reported_by, status, created_by, file_paths) VALUES (?, ?, 0, ?, ?, ?, ?, ?, 0, ?, ?)`, req.PhoneID, ticketno, req.IssueTypeID, req.IssueElse, req.DepartmentID, req.Text, req.ReportedBy, req.CreatedBy, string(filePathsBytes))
		}
	} else {
		if req.SystemID > 0 {
			var typeid int
			db.DB.QueryRow(`SELECT type FROM systems_program WHERE id = ?`, req.SystemID).Scan(&typeid)
			res, err = db.DB.Exec(`INSERT INTO tasks (phone_id, ticket_no, system_id, issue_type, department_id, text, reported_by, status, created_by) VALUES (?, ?, ?, ?, ?, ?, ?, 0, ?)`, req.PhoneID, ticketno, req.SystemID, typeid, req.DepartmentID, req.Text, req.ReportedBy, req.CreatedBy)
		} else {
			// กรณีไม่มีระบบ ให้เก็บ issue_type และ issue_else
			if req.IssueTypeID == 0 || req.IssueElse == "" {
				return c.Status(400).JSON(fiber.Map{"error": "issue_type and issue_else must be provided when system_id is 0"})
			}
			// รับค่า issue_type จาก frontend โดยตรง
			res, err = db.DB.Exec(`INSERT INTO tasks (phone_id, ticket_no, system_id, issue_type, issue_else, department_id, text, reported_by, status, created_by) VALUES (?, ?, 0, ?, ?, ?, ?, ?, 0, ?)`, req.PhoneID, ticketno, req.IssueTypeID, req.IssueElse, req.DepartmentID, req.Text, req.ReportedBy, req.CreatedBy)
		}
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to insert task"})
	}
	id, _ := res.LastInsertId()

	// Update department score
	if err := updateDepartmentScore(req.DepartmentID); err != nil {
		log.Printf("Failed to update department score: %v", err)
	}

	log.Printf("Inserted new task with ID: %d", id)
	var Urlenv string
	env := os.Getenv("env")

	if env == "dev" {
		Urlenv = "http://helpdesk-dev.nopadol.com/tasks/show/" + strconv.Itoa(int(id))
	} else {
		Urlenv = "http://helpdesk.nopadol.com/tasks/show/" + strconv.Itoa(int(id))
	}
	log.Printf("🔔 Telegram flag: %t", req.Telegram)
	if req.Telegram {
		log.Printf("🚀 Preparing to send Telegram notification for task ID: %d", id)
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
		req.Ticket = ticketno
		req.BranchName = branchName
		req.ProgramName = programName
		req.Url = Urlenv
		req.CreatedAt = time.Now().Add(7 * time.Hour).Format("02/01/2006 15:04:05")
		req.Status = 0

		var messageID int
		var messageName string
		// Send with photo if files were uploaded
		if len(uploadedFiles) > 0 {
			// Get all image URLs
			var photoURLs []string
			for _, file := range uploadedFiles {
				if url, ok := file["url"].(string); ok {
					photoURLs = append(photoURLs, url)
				}
			}
			messageID, messageName, err = common.SendTelegram(req, photoURLs...)
			if err != nil {
				log.Printf("❌ Error sending Telegram: %v", err)
			}
		} else {
			messageID, messageName, err = common.SendTelegram(req)
			if err != nil {
				log.Printf("❌ Error sending Telegram: %v", err)
			}
		}

		if err == nil {
			log.Printf("✅ Telegram sent successfully, updating message_id: %d for task: %d", messageID, id)
			// Update task with message_id
			chatID, _ := strconv.Atoi(os.Getenv("CHAT_ID"))
			resTG, errTG := db.DB.Exec(`INSERT INTO telegram_chat (chat_id, chat_name, report_id) VALUES (?, ?, ?)`, chatID, messageName, messageID)
			var telegramChatID int64
			if errTG != nil {
				log.Printf("❌ Failed to Chat telegram: %v", errTG)
			} else {
				telegramChatID, _ = resTG.LastInsertId()
				log.Printf("✅ Message ID Chat telegram successfully in database, telegram_chat.id: %d", telegramChatID)
			}
			_, err = db.DB.Exec(`UPDATE tasks SET telegram_id = ? WHERE id = ?`, telegramChatID, id)
			if err != nil {
				log.Printf("❌ Failed to update telegram_id: %v", err)
			} else {
				log.Printf("✅ Telegram ID updated successfully in database")
			}
		} else {
			log.Printf("❌ Failed to send Telegram notification: %v", err)
		}
	} else {
		log.Printf("⚠️ Telegram notification skipped - Telegram flag is false")
	}
	return c.JSON(fiber.Map{"success": true, "id": id})
}

// UpdateTaskHandler แก้ไข task
// @Summary Update problem
// @Description Update an existing problem
// @Tags problems
// @Accept json
// @Produce json
// @Param id path string true "Problem ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/problem/update/{id} [put]
func UpdateTaskHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	var err error
	var ticketno string
	var uploadedFiles []fiber.Map

	var req models.TaskRequestUpdate
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Handle phone_id = 0 as null
	if req.PhoneID != nil && *req.PhoneID == 0 {
		req.PhoneID = nil
	}

	log.Printf("Looking for task with ID: %s", id)
	err = db.DB.QueryRow("SELECT ticket_no FROM tasks WHERE id = ?", id).Scan(&ticketno)
	if err != nil {
		log.Printf("Task not found error: %v", err)
		return c.Status(404).JSON(fiber.Map{"error": "Task not found"})
	}
	log.Printf("Found task with ticket_no: %s", ticketno)

	form, err := c.MultipartForm()
	if err != nil {
		// If multipart parsing fails, try regular body parser
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
		}
	} else {
		// Parse body from multipart form
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
		}

		// Handle file uploads if present (support image_{index} format)
		var allFiles []*multipart.FileHeader

		// Check for indexed files (image_0, image_1, image_2, etc.)
		for key, files := range form.File {
			if strings.HasPrefix(key, "image_") || key == "image" {
				allFiles = append(allFiles, files...)
			}
		}

		uploadedFiles, _ = common.HandleFileUploads(allFiles, ticketno)
	}

	// Convert string form values to int for multipart data
	if phoneIDStr := c.FormValue("phone_id"); phoneIDStr != "" && phoneIDStr != "0" {
		if phoneID, err := strconv.Atoi(phoneIDStr); err == nil {
			req.PhoneID = &phoneID
		}
	}
	if systemIDStr := c.FormValue("system_id"); systemIDStr != "" {
		req.SystemID, _ = strconv.Atoi(systemIDStr)
	}
	if departmentIDStr := c.FormValue("department_id"); departmentIDStr != "" {
		req.DepartmentID, _ = strconv.Atoi(departmentIDStr)
	}
	if issueTypeIDStr := c.FormValue("issue_type"); issueTypeIDStr != "" {
		req.IssueTypeID, _ = strconv.Atoi(issueTypeIDStr)
	}
	if issueElseStr := c.FormValue("issue_else"); issueElseStr != "" {
		req.IssueElse = issueElseStr
	}
	if assignedToIDStr := c.FormValue("assignedto_id"); assignedToIDStr != "" {
		req.AssignedtoID, _ = strconv.Atoi(assignedToIDStr)
	}
	if textStr := c.FormValue("text"); textStr != "" {
		req.Text = textStr
	}
	if reportedByStr := c.FormValue("reported_by"); reportedByStr != "" {
		req.ReportedBy = &reportedByStr
	}
	var previousAssigntoNull sql.NullString // เพิ่มตัวแปรเก็บ assignto เดิม
	err = db.DB.QueryRow(`SELECT assignto FROM tasks WHERE id = ?`, id).Scan(&previousAssigntoNull)
	if err != nil {
		log.Println("Error fetching previous_assignto:", err)
	}
	var previousAssignto string
	if previousAssigntoNull.Valid {
		previousAssignto = previousAssigntoNull.String
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
	var createdAtStr string
	var CreatedAt time.Time

	log.Printf("Updating task ID: %s", id)

	// Handle file uploads
	if len(uploadedFiles) > 0 {
		// Delete existing files first
		var existingFilePathsJSON string
		db.DB.QueryRow(`SELECT IFNULL(file_paths, '[]') FROM tasks WHERE id = ?`, id).Scan(&existingFilePathsJSON)

		if existingFilePathsJSON != "" && existingFilePathsJSON != "[]" {
			var existingFiles []fiber.Map
			if err := json.Unmarshal([]byte(existingFilePathsJSON), &existingFiles); err == nil {
				for _, fp := range existingFiles {
					if url, ok := fp["url"].(string); ok {
						// Extract object name from URL
						if strings.Contains(url, "prefix=") {
							parts := strings.Split(url, "prefix=")
							if len(parts) > 1 {
								objectName := parts[1]
								common.DeleteImage(objectName)
							}
						}
					}
				}
			}
		}

		// Upload new files
		filePathsBytes, _ := json.Marshal(uploadedFiles)
		log.Printf("Updating file_paths: %s", string(filePathsBytes))

		if req.SystemID > 0 {
			var typeid int
			db.DB.QueryRow(`SELECT type FROM systems_program WHERE id = ?`, req.SystemID).Scan(&typeid)
			if req.ReportedBy != nil {
				_, err = db.DB.Exec(`UPDATE tasks SET phone_id=?, system_id=?, issue_type=?, issue_else=NULL, department_id=?, assignto_id=?, assignto=?, reported_by=?, text=?, status=?, updated_at=CURRENT_TIMESTAMP, updated_by=?, file_paths=? WHERE id=?`, req.PhoneID, req.SystemID, typeid, req.DepartmentID, req.AssignedtoID, req.Assignto, req.ReportedBy, req.Text, req.Status, req.UpdatedBy, string(filePathsBytes), id)
			} else {
				_, err = db.DB.Exec(`UPDATE tasks SET phone_id=?, system_id=?, issue_type=?, issue_else=NULL, department_id=?, assignto_id=?, assignto=?, text=?, status=?, updated_at=CURRENT_TIMESTAMP, updated_by=?, file_paths=? WHERE id=?`, req.PhoneID, req.SystemID, typeid, req.DepartmentID, req.AssignedtoID, req.Assignto, req.Text, req.Status, req.UpdatedBy, string(filePathsBytes), id)
			}
		} else {
			if req.ReportedBy != nil {
				_, err = db.DB.Exec(`UPDATE tasks SET phone_id=?, system_id=0, issue_type=?, issue_else=?, department_id=?, assignto_id=?, assignto=?, reported_by=?, text=?, status=?, updated_at=CURRENT_TIMESTAMP, updated_by=?, file_paths=? WHERE id=?`, req.PhoneID, req.IssueTypeID, req.IssueElse, req.DepartmentID, req.AssignedtoID, req.Assignto, req.ReportedBy, req.Text, req.Status, req.UpdatedBy, string(filePathsBytes), id)
			} else {
				_, err = db.DB.Exec(`UPDATE tasks SET phone_id=?, system_id=0, issue_type=?, issue_else=?, department_id=?, assignto_id=?, assignto=?, text=?, status=?, updated_at=CURRENT_TIMESTAMP, updated_by=?, file_paths=? WHERE id=?`, req.PhoneID, req.IssueTypeID, req.IssueElse, req.DepartmentID, req.AssignedtoID, req.Assignto, req.Text, req.Status, req.UpdatedBy, string(filePathsBytes), id)
			}
		}
	} else {
		if req.SystemID > 0 {
			var typeid int
			db.DB.QueryRow(`SELECT type FROM systems_program WHERE id = ?`, req.SystemID).Scan(&typeid)
			if req.ReportedBy != nil {
				_, err = db.DB.Exec(`UPDATE tasks SET phone_id=?, system_id=?, issue_type=?, issue_else=NULL, department_id=?, assignto_id=?, assignto=?, reported_by=?, text=?, status=?, updated_at=CURRENT_TIMESTAMP, updated_by=? WHERE id=?`, req.PhoneID, req.SystemID, typeid, req.DepartmentID, req.AssignedtoID, req.Assignto, req.ReportedBy, req.Text, req.Status, req.UpdatedBy, id)
			} else {
				_, err = db.DB.Exec(`UPDATE tasks SET phone_id=?, system_id=?, issue_type=?, issue_else=NULL, department_id=?, assignto_id=?, assignto=?, text=?, status=?, updated_at=CURRENT_TIMESTAMP, updated_by=? WHERE id=?`, req.PhoneID, req.SystemID, typeid, req.DepartmentID, req.AssignedtoID, req.Assignto, req.Text, req.Status, req.UpdatedBy, id)
			}
		} else {
			var assignto string
			if req.Assignto != nil {
				assignto = *req.Assignto
			}
			if req.ReportedBy != nil {
				_, err = db.DB.Exec(`UPDATE tasks SET phone_id=?, system_id=0, issue_type=?, issue_else=?, department_id=?, assignto_id=?, assignto=?, reported_by=?, text=?, status=?, updated_at=CURRENT_TIMESTAMP, updated_by=? WHERE id=?`, req.PhoneID, req.IssueTypeID, req.IssueElse, req.DepartmentID, req.AssignedtoID, assignto, req.ReportedBy, req.Text, req.Status, req.UpdatedBy, id)
			} else {
				_, err = db.DB.Exec(`UPDATE tasks SET phone_id=?, system_id=0, issue_type=?, issue_else=?, department_id=?, assignto_id=?, assignto=?, text=?, status=?, updated_at=CURRENT_TIMESTAMP, updated_by=? WHERE id=?`, req.PhoneID, req.IssueTypeID, req.IssueElse, req.DepartmentID, req.AssignedtoID, assignto, req.Text, req.Status, req.UpdatedBy, id)
			}
		}
	}

	if req.Status == 1 {
		_, err = db.DB.Exec(`UPDATE tasks SET resolved_at=CURRENT_TIMESTAMP WHERE id=?`, id)
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update task"})
	}
	var Urlenv string
	env := config.AppConfig.Environment
	if env == "dev" {
		Urlenv = "http://helpdesk-dev.nopadol.com/tasks/show/" + id
	} else {
		Urlenv = "http://helpdesk.nopadol.com/tasks/show/" + id
	}
	// Get message_id and update Telegram if exists
	var messageID int
	var reported string
	var existingFilePathsJSON string
	var telegramUser string
	var assigntoID int
	var ResolvedAt sql.NullTime
	var telegramID int

	err = db.DB.QueryRow(`
		SELECT IFNULL(t.ticket_no, ''),IFNULL(tc.report_id, 0), IFNULL(t.reported_by, ''), 
		IFNULL(t.file_paths, '[]'), IFNULL(rs.telegram_username, ''), IFNULL(tc.assignto_id, 0), t.created_at, IFNULL(t.telegram_id, 0), t.resolved_at
		FROM tasks t
		LEFT JOIN telegram_chat tc ON t.telegram_id = tc.id
		LEFT JOIN responsibilities rs ON t.assignto_id = rs.id
		WHERE t.id = ?
		`, id).Scan(&ticketno, &messageID, &reported, &existingFilePathsJSON, &telegramUser, &assigntoID, &createdAtStr, &telegramID, &ResolvedAt)

	// Parse created_at string to time
	if createdAtStr != "" {
		CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			log.Printf("Error parsing created_at: %v", err)
			CreatedAt = time.Now() // fallback
		}
	} else {
		CreatedAt = time.Now() // fallback
	}

	log.Printf("Query result - err: %v, messageID: %d, telegramID: %d", err, messageID, telegramID)

	if err == nil && messageID > 0 {
		// Create TaskRequest from TaskRequestUpdate for Telegram
		telegramReq := models.TaskRequest{
			PhoneID:          req.PhoneID,
			SystemID:         req.SystemID,
			DepartmentID:     req.DepartmentID,
			Text:             req.Text,
			Status:           req.Status,
			ReportedBy:       reported,
			TelegramUser:     telegramUser,
			AssigntoID:       assigntoID,
			Assignto:         "",
			PreviousAssignto: previousAssignto,
			MessageID:        messageID,

			Ticket: ticketno,
		}

		// Get additional data for Telegram
		var phoneNumber int
		var departmentName, branchName, programName string

		if req.PhoneID != nil {
			db.DB.QueryRow(`SELECT p.number, d.name, b.name FROM ip_phones p JOIN departments d ON p.department_id = d.id JOIN branches b ON d.branch_id = b.id WHERE p.id = ?`, *req.PhoneID).Scan(&phoneNumber, &departmentName, &branchName)
		} else {
			db.DB.QueryRow(`SELECT d.name, b.name FROM departments d JOIN branches b ON d.branch_id = b.id WHERE d.id = ?`, req.DepartmentID).Scan(&departmentName, &branchName)
		}

		if req.SystemID != 0 {
			db.DB.QueryRow(`SELECT name FROM systems_program WHERE id = ?`, req.SystemID).Scan(&programName)
		}

		telegramReq.PhoneNumber = phoneNumber
		telegramReq.DepartmentName = departmentName
		telegramReq.BranchName = branchName
		telegramReq.ProgramName = programName
		telegramReq.Url = Urlenv
		telegramReq.PreviousAssignto = previousAssignto
		telegramReq.CreatedAt = CreatedAt.Add(7 * time.Hour).Format("02/01/2006 15:04:05")
		if ResolvedAt.Valid {
			telegramReq.UpdatedAt = ResolvedAt.Time.Add(7 * time.Hour).Format("02/01/2006 15:04:05")
		} else {
			telegramReq.UpdatedAt = ""
		}
		telegramReq.Ticket = ticketno
		telegramReq.ReportedBy = reported
		telegramReq.TelegramUser = telegramUser
		telegramReq.AssigntoID = assigntoID
		if req.Assignto != nil {
			telegramReq.Assignto = *req.Assignto
		}

		// Get first image URL from existing files for Telegram
		var photoURLs []string
		if existingFilePathsJSON != "" && existingFilePathsJSON != "[]" {
			var existingFiles []fiber.Map
			if err := json.Unmarshal([]byte(existingFilePathsJSON), &existingFiles); err == nil && len(existingFiles) > 0 {
				for _, file := range existingFiles {
					if url, ok := file["url"].(string); ok {
						photoURLs = append(photoURLs, url)
					}
				}
			}
		}
		log.Printf("Previous assignto: %s, New assignto: %v", previousAssignto, req.Assignto)

		// ตรวจสอบการเปลี่ยนผู้รับผิดชอบ
		var currentAssignto string
		if req.Assignto != nil {
			currentAssignto = *req.Assignto
		}

		if previousAssignto != currentAssignto && previousAssignto != "" {
			// ลบ notification message เก่าถ้ามีการเปลี่ยนผู้รับผิดชอบ
			var oldNotificationID int
			db.DB.QueryRow(`SELECT IFNULL(assignto_id, 0) FROM telegram_chat WHERE id = ?`, telegramID).Scan(&oldNotificationID)
			if oldNotificationID > 0 {
				_, _ = common.DeleteTelegram(oldNotificationID)
			}
		}

		// ส่งหลายไฟล์
		var notificationResp int
		if len(photoURLs) > 0 {
			notificationResp, _ = common.UpdateTelegram(telegramReq, photoURLs...)
		} else {
			notificationResp, _ = common.UpdateTelegram(telegramReq)
		}

		_, err = db.DB.Exec(`UPDATE telegram_chat SET assignto_id = ? WHERE id = ?`, notificationResp, telegramID)
		if err != nil {
			log.Printf("❌ Failed to update assignto_id: %v", err)
		} else {
			log.Printf("✅ Assignto ID updated successfully in database")
		}

		// ใช้ UpdatereplyToSpecificMessage เมื่อมีการอัปเดต assignto
		if previousAssignto != currentAssignto {
			// ดึงข้อมูลจาก resolutions table
			var resolutionID sql.NullInt64
			var solutionMessageID int
			db.DB.QueryRow(`SELECT solution_id FROM tasks WHERE id = ?`, id).Scan(&resolutionID)
			db.DB.QueryRow(`SELECT IFNULL(solution_id, 0) FROM telegram_chat WHERE id = ?`, telegramID).Scan(&solutionMessageID)

			if resolutionID.Valid && solutionMessageID > 0 {
				var resolutionText string
				var resolutionFilePathsJSON string
				var resolutionResolvedAt string
				err = db.DB.QueryRow(`
					SELECT IFNULL(text, ''), IFNULL(file_paths, '[]'), 
					DATE_FORMAT(resolved_at, '%d/%m/%Y %H:%i:%s') 
					FROM resolutions WHERE id = ?
				`, resolutionID.Int64).Scan(&resolutionText, &resolutionFilePathsJSON, &resolutionResolvedAt)

				if err == nil {
					// สร้าง ResolutionReq
					resolutionReq := models.ResolutionReq{
						Solution:     resolutionText,
						TelegramUser: telegramUser,
						MessageID:    messageID,
						Url:          Urlenv,
						Assignto:     currentAssignto,
						TicketNo:     ticketno,
						CreatedAt:    CreatedAt.Add(7 * time.Hour).Format("02/01/2006 15:04:05"),
						ResolvedAt:   resolutionResolvedAt,
					}

					// ดึง photo URLs จาก resolution files
					var resolutionPhotoURLs []string
					if resolutionFilePathsJSON != "" && resolutionFilePathsJSON != "[]" {
						var resolutionFiles []fiber.Map
						if err := json.Unmarshal([]byte(resolutionFilePathsJSON), &resolutionFiles); err == nil {
							for _, file := range resolutionFiles {
								if url, ok := file["url"].(string); ok {
									resolutionPhotoURLs = append(resolutionPhotoURLs, url)
								}
							}
						}
					}

					// เรียกใช้ UpdatereplyToSpecificMessage
					newSolutionMessageID, err := common.UpdatereplyToSpecificMessage(solutionMessageID, resolutionReq, resolutionPhotoURLs...)
					if err != nil {
						log.Printf("❌ Failed to update resolution message: %v", err)
					} else {
						// อัปเดต solution_id ใน telegram_chat
						_, err = db.DB.Exec(`UPDATE telegram_chat SET solution_id = ? WHERE id = ?`, newSolutionMessageID, telegramID)
						if err != nil {
							log.Printf("❌ Failed to update solution_id: %v", err)
						} else {
							log.Printf("✅ Resolution message updated successfully")
						}
					}
				}
			}
		}
	}
	log.Printf("Updating task ID: %s", id)
	return c.JSON(fiber.Map{"success": true})
}

// DeleteTaskHandler (soft delete)
// @Summary Delete problem
// @Description Delete a problem and related data
// @Tags problems
// @Accept json
// @Produce json
// @Param id path string true "Problem ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/problem/delete/{id} [delete]
func DeleteTaskHandler(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid id"})
	}

	// Get message_id, file_paths and resolution data before deleting
	var messageID int
	var filePathsJSON string
	var solutionID *int
	var telegramID int
	var solutionMessageID int
	var assigntoID int
	err = db.DB.QueryRow(`
		SELECT IFNULL(tc.report_id, 0), IFNULL(t.file_paths, '[]'), t.solution_id, IFNULL(t.telegram_id, 0), IFNULL(tc.assignto_id, 0) FROM tasks t
		LEFT JOIN telegram_chat tc ON t.telegram_id = tc.id
		WHERE t.id = ?
		`, id).Scan(&messageID, &filePathsJSON, &solutionID, &telegramID, &assigntoID)
	if err != nil {
		log.Printf("Failed to get task data: %v", err)
	}

	// Get solution message ID if exists
	if telegramID > 0 {
		db.DB.QueryRow(`SELECT IFNULL(solution_id, 0) FROM telegram_chat WHERE id = ?`, telegramID).Scan(&solutionMessageID)
	}

	// Delete resolution files from MinIO if they exist
	if solutionID != nil {
		var resolutionFilePathsJSON string
		err = db.DB.QueryRow(`SELECT IFNULL(file_paths, '[]') FROM resolutions WHERE id = ?`, *solutionID).Scan(&resolutionFilePathsJSON)
		if err == nil && resolutionFilePathsJSON != "" && resolutionFilePathsJSON != "[]" {
			var resolutionFiles []fiber.Map
			if err := json.Unmarshal([]byte(resolutionFilePathsJSON), &resolutionFiles); err == nil {
				for _, fp := range resolutionFiles {
					if url, ok := fp["url"].(string); ok {
						if strings.Contains(url, "prefix=") {
							parts := strings.Split(url, "prefix=")
							if len(parts) > 1 {
								objectName := parts[1]
								common.DeleteImage(objectName)
								log.Printf("Deleted resolution file: %s", objectName)
							}
						}
					}
				}
			}
		}
	}

	// Delete task files from MinIO if they exist
	if filePathsJSON != "" && filePathsJSON != "[]" {
		var filePaths []fiber.Map
		if err := json.Unmarshal([]byte(filePathsJSON), &filePaths); err == nil {
			for _, fp := range filePaths {
				if url, ok := fp["url"].(string); ok {
					if strings.Contains(url, "prefix=") {
						parts := strings.Split(url, "prefix=")
						if len(parts) > 1 {
							objectName := parts[1]
							common.DeleteImage(objectName)
						}
					}
				}
			}
		}
	}
	log.Printf("Deleting Telegram message ID: %d", assigntoID)
	if assigntoID > 0 {
		_, _ = common.DeleteTelegram(assigntoID)
	}
	// Delete resolution if exists
	if solutionID != nil {
		_, err = db.DB.Exec(`DELETE FROM resolutions WHERE id = ?`, *solutionID)
		if err != nil {
			log.Printf("Failed to delete resolution: %v", err)
		}
	}

	// Delete telegram_chat

	if telegramID > 0 {
		_, err = db.DB.Exec(`DELETE FROM telegram_chat WHERE id = ?`, telegramID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to delete telegram chat"})
		}
	}

	// Delete task
	_, err = db.DB.Exec(`DELETE FROM tasks WHERE id = ?`, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete task"})
	}

	// Delete Telegram messages if they exist
	if messageID > 0 {
		_, _ = common.DeleteTelegram(messageID)
	}
	if solutionMessageID > 0 {
		_, _ = common.DeleteTelegram(solutionMessageID)
	}

	log.Printf("Deleted task ID: %d", id)
	return c.JSON(fiber.Map{"success": true})
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

// @Summary Search problems
// @Description Search problems by query string with pagination
// @Tags problems
// @Accept json
// @Produce json
// @Param query path string true "Search query (use 'all' for all problems)"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} models.PaginatedResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/problem/list/{query} [get]
func GetTasksWithQueryHandler(c *fiber.Ctx) error {
	query := c.Params("query")

	// If query is empty or "all", return all tasks with pagination
	if query == "" || query == "all" {
		return GetTasksHandler(c)
	}

	// Otherwise, search tasks
	pagination := utils.GetPaginationParams(c)
	offset := utils.CalculateOffset(pagination.Page, pagination.Limit)

	// URL decode for Thai language support
	decodedQuery, err := url.QueryUnescape(query)
	if err != nil {
		decodedQuery = query
	}

	// Clean and validate query
	decodedQuery = strings.TrimSpace(decodedQuery)
	if len(decodedQuery) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Search query cannot be empty"})
	}
	if len(decodedQuery) > 100 {
		return c.Status(400).JSON(fiber.Map{"error": "Search query too long"})
	}

	// Clean query - remove SQL wildcards and quotes
	decodedQuery = strings.ReplaceAll(decodedQuery, "  ", " ")
	decodedQuery = strings.ReplaceAll(decodedQuery, "%", "")
	decodedQuery = strings.ReplaceAll(decodedQuery, "_", "")
	decodedQuery = strings.ReplaceAll(decodedQuery, "'", "")
	decodedQuery = strings.ReplaceAll(decodedQuery, "\"", "")
	decodedQuery = strings.ReplaceAll(decodedQuery, ";", "")
	decodedQuery = strings.ReplaceAll(decodedQuery, "--", "")
	decodedQuery = strings.ReplaceAll(decodedQuery, "/*", "")
	decodedQuery = strings.ReplaceAll(decodedQuery, "*/", "")

	searchPattern := "%" + decodedQuery + "%"

	// Get total count
	var total int
	err = db.DB.QueryRow(`
		SELECT COUNT(*) FROM tasks t
		LEFT JOIN issue_types it ON t.issue_type = it.id
		LEFT JOIN ip_phones p ON t.phone_id = p.id
		LEFT JOIN departments d ON t.department_id = d.id
		LEFT JOIN branches b ON d.branch_id = b.id
		LEFT JOIN systems_program s ON t.system_id = s.id
		WHERE (t.ticket_no LIKE ? OR p.number LIKE ? OR p.name LIKE ? OR d.name LIKE ? OR b.name LIKE ? OR s.name LIKE ? OR t.text LIKE ? OR it.name LIKE ? OR t.status LIKE ?)
	`, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern).Scan(&total)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to count search results"})
	}

	// Get search results with pagination
	rows, err := db.DB.Query(`
		SELECT t.id, IFNULL(t.ticket_no, ''), IFNULL(t.phone_id, 0), IFNULL(p.number, 0) as number, IFNULL(p.name, ''), t.system_id, IFNULL(s.name, ''), IFNULL(t.issue_type, 0), IFNULL(t.issue_else, ''), IFNULL(it.name, ''), IFNULL(t.department_id, 0), IFNULL(d.name, ''), IFNULL(d.branch_id, 0), IFNULL(b.name, ''), t.text, IFNULL(t.assignto, ''), t.status, t.created_at, t.updated_at
		FROM tasks t
		LEFT JOIN ip_phones p ON t.phone_id = p.id
		LEFT JOIN departments d ON t.department_id = d.id
		LEFT JOIN branches b ON d.branch_id = b.id
		LEFT JOIN systems_program s ON t.system_id = s.id
		LEFT JOIN issue_types it ON t.issue_type = it.id
		WHERE (t.ticket_no LIKE ? OR p.number LIKE ? OR p.name LIKE ? OR d.name LIKE ? OR b.name LIKE ? OR s.name LIKE ? OR t.text LIKE ? OR it.name LIKE ? OR t.assignto LIKE ? OR t.status LIKE ?)
		ORDER BY t.id DESC
		LIMIT ? OFFSET ?
	`, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, pagination.Limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to search tasks"})
	}
	defer rows.Close()

	var tasks []models.TaskWithDetails
	for rows.Next() {
		var t models.TaskWithDetails
		var issueTypeName string
		err := rows.Scan(&t.ID, &t.Ticket, &t.PhoneID, &t.Number, &t.PhoneName, &t.SystemID, &t.SystemName, &t.IssueTypeID, &t.IssueElse, &issueTypeName, &t.DepartmentID, &t.DepartmentName, &t.BranchID, &t.BranchName, &t.Text, &t.Assignto, &t.Status, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			log.Printf("Error scanning task: %v", err)
			continue
		}
		if t.SystemID > 0 {
			t.SystemType = issueTypeName
		} else {
			t.SystemType = issueTypeName
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

	log.Printf("Searching tasks with query: %s, found %d results", query, len(tasks))
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

// @Summary Search problems by column
// @Description Search problems by specific column and value with pagination
// @Tags problems
// @Accept json
// @Produce json
// @Param column path string true "Column name to search in"
// @Param query path string true "Search value"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} models.PaginatedResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/problem/list/{column}/{query} [get]
func GetTasksWithColumnQueryHandler(c *fiber.Ctx) error {
	column := c.Params("column")
	query := c.Params("query")

	if column == "" || query == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Column and query parameters are required"})
	}

	// Get pagination parameters
	pagination := utils.GetPaginationParams(c)
	offset := utils.CalculateOffset(pagination.Page, pagination.Limit)

	// Define column types and validation
	intColumns := map[string]bool{
		"phone_id":      true,
		"issue_type":    true,
		"system_id":     true,
		"department_id": true,
		"branch_id":     true,
		"status":        true,
		"created_by":    true,
		"updated_by":    true,
		"telegram_id":   true,
	}

	stringColumns := map[string]bool{
		"ticket_no":       true,
		"number":          true,
		"phone_name":      true,
		"system_name":     true,
		"issue_else":      true,
		"department_name": true,
		"branch_name":     true,
		"text":            true,
		"assignto":        true,
		"reported_by":     true,
		"solution":        true,
	}

	// Validate column exists
	if !intColumns[column] && !stringColumns[column] {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid column name"})
	}

	// Validate query parameter based on column type
	var queryParam interface{}
	var err error

	if intColumns[column] {
		// For integer columns, convert string to int
		intVal, err := strconv.Atoi(query)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": fmt.Sprintf("Invalid value for column %s: must be integer", column)})
		}
		queryParam = intVal
	} else {
		// For string columns, use as is
		queryParam = query
	}

	var queryStr string
	var rows *sql.Rows

	// Declare total variable for pagination
	var total int

	// Column mapping to prevent SQL injection
	columnMap := map[string]string{
		"phone_id":        "t.phone_id",
		"issue_type":      "t.issue_type",
		"system_id":       "t.system_id",
		"department_id":   "t.department_id",
		"branch_id":       "d.branch_id",
		"status":          "t.status",
		"created_by":      "t.created_by",
		"updated_by":      "t.updated_by",
		"telegram_id":     "t.telegram_id",
		"phone_name":      "p.name",
		"number":          "p.number",
		"system_name":     "s.name",
		"department_name": "d.name",
		"branch_name":     "b.name",
		"solution":        "t.solution",
		"reported_by":     "t.reported_by",
		"ticket_no":       "t.ticket_no",
		"issue_else":      "t.issue_else",
		"text":            "t.text",
		"assignto":        "t.assignto",
	}

	sqlColumn, exists := columnMap[column]
	if !exists {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid column name"})
	}

	baseQuery := `FROM tasks t
		LEFT JOIN ip_phones p ON t.phone_id = p.id
		LEFT JOIN departments d ON t.department_id = d.id
		LEFT JOIN branches b ON d.branch_id = b.id
		LEFT JOIN systems_program s ON t.system_id = s.id
		LEFT JOIN issue_types it ON t.issue_type = it.id`

	selectFields := `t.id, IFNULL(t.ticket_no, ''), IFNULL(t.phone_id, 0), IFNULL(p.number, 0), IFNULL(p.name, ''), 
		t.system_id, IFNULL(s.name, ''), IFNULL(t.issue_type, 0), IFNULL(t.issue_else, ''), 
		IFNULL(it.name, ''), IFNULL(t.department_id, 0), IFNULL(d.name, ''), IFNULL(d.branch_id, 0), 
		IFNULL(b.name, ''), t.text, IFNULL(t.assignto, ''), t.status, t.created_at, t.updated_at`

	// Build query based on column type
	if intColumns[column] {
		// Get total count for pagination
		countQuery := fmt.Sprintf("SELECT COUNT(*) %s WHERE %s = ?", baseQuery, sqlColumn)
		db.DB.QueryRow(countQuery, queryParam).Scan(&total)

		queryStr = fmt.Sprintf("SELECT %s %s WHERE %s = ? ORDER BY t.id DESC LIMIT ? OFFSET ?",
			selectFields, baseQuery, sqlColumn)
		rows, err = db.DB.Query(queryStr, queryParam, pagination.Limit, offset)
	} else {
		// For string columns, use LIKE search
		searchPattern := "%" + query + "%"

		// Get total count for pagination
		countQuery := fmt.Sprintf("SELECT COUNT(*) %s WHERE %s LIKE ?", baseQuery, sqlColumn)
		db.DB.QueryRow(countQuery, searchPattern).Scan(&total)

		queryStr = fmt.Sprintf("SELECT %s %s WHERE %s LIKE ? ORDER BY t.id DESC LIMIT ? OFFSET ?",
			selectFields, baseQuery, sqlColumn)
		rows, err = db.DB.Query(queryStr, searchPattern, pagination.Limit, offset)
	}

	if err != nil {
		log.Printf("Query error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to query tasks"})
	}
	defer rows.Close()

	var tasks []models.TaskWithDetails
	for rows.Next() {
		var t models.TaskWithDetails
		var issueTypeName string
		err := rows.Scan(&t.ID, &t.Ticket, &t.PhoneID, &t.Number, &t.PhoneName, &t.SystemID, &t.SystemName,
			&t.IssueTypeID, &t.IssueElse, &issueTypeName, &t.DepartmentID, &t.DepartmentName,
			&t.BranchID, &t.BranchName, &t.Text, &t.Assignto, &t.Status, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			log.Printf("Error scanning task: %v", err)
			continue
		}

		t.SystemType = issueTypeName

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

// @Summary Update assigned person
// @Description Update the assigned person for a specific problem
// @Tags problems
// @Accept json
// @Produce json
// @Param id path string true "Problem ID"
// @Param request body object true "Assignment update data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/problem/update/assignto/{id} [put]
func UpdateAssignedTo(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Task ID is required"})
	}
	var Urlenv string
	env := config.AppConfig.Environment

	var req models.AssignRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var messageID int

	err := db.DB.QueryRow(`
		SELECT IFNULL(tc.assignto_id, 0) as assignto_id
		FROM tasks t
		LEFT JOIN telegram_chat tc ON t.telegram_id = tc.id
		WHERE t.id = ?
		`, id).Scan(&messageID)
	if err != nil {
		log.Printf("Failed to get task data: %v", err)
	}

	if messageID > 0 {
		_, _ = common.DeleteTelegram(messageID)
	}

	_, err = db.DB.Exec(`UPDATE tasks SET assignto_id = ?, assignto = ?, updated_by = ?, updated_at = NOW() WHERE id = ?`, req.AssignedtoID, req.Assignto, req.UpdatedBy, id)
	if err != nil {
		log.Printf("Database error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update assigned person"})
	}
	// เฉพาะกรณีที่ต้องการอัพเดต Telegram
	if req.UpdateTelegram {

		// ดึงข้อมูล task สำหรับอัพเดต Telegram
		var ticket, text, issueElse, reportedBy, assignto, createdAt, updatedAt, branchName, departmentName, programName string
		var phoneID, systemID, departmentID, status, messageID, phoneNumber, branchID, telegramID int
		var filePathsJSON string
		var telegramUser string
		// Fix SQL: JOINs before WHERE, select tc.report_id as messageID
		err := db.DB.QueryRow(`
			SELECT IFNULL(t.ticket_no, ''), IFNULL(t.phone_id, 0), IFNULL(t.system_id, 0), IFNULL(t.issue_else, ''), IFNULL(t.department_id, 0),
			IFNULL(t.text, ''), IFNULL(t.status, 0), IFNULL(t.reported_by, ''), IFNULL(t.assignto, ''), IFNULL(rs.telegram_username, ''), 
			IFNULL(tc.report_id, 0), IFNULL(t.file_paths, '[]'), IFNULL(d.branch_id, 0), IFNULL(t.created_at, ''), IFNULL(t.updated_at, ''),
			IFNULL(t.telegram_id, 0)
			FROM tasks t
			LEFT JOIN telegram_chat tc ON t.telegram_id = tc.id
			LEFT JOIN departments d ON t.department_id = d.id
			LEFT JOIN branches b ON d.branch_id = b.id
			LEFT JOIN systems_program s ON t.system_id = s.id
			LEFT JOIN responsibilities rs ON t.assignto_id = rs.id
			WHERE t.id = ?
		`, id).Scan(&ticket, &phoneID, &systemID, &issueElse, &departmentID, &text, &status, &reportedBy, &assignto, &telegramUser, &messageID, &filePathsJSON, &branchID, &createdAt, &updatedAt, &telegramID)
		log.Printf("Fetched task for Telegram update, ID: %s, MessageID: %d", telegramUser, messageID)
		// Query extra info for Telegram
		db.DB.QueryRow(`SELECT name FROM branches WHERE id = ?`, branchID).Scan(&branchName)
		db.DB.QueryRow(`SELECT name FROM departments WHERE id = ?`, departmentID).Scan(&departmentName)
		db.DB.QueryRow(`SELECT name FROM systems_program WHERE id = ?`, systemID).Scan(&programName)
		db.DB.QueryRow(`SELECT number FROM ip_phones WHERE id = ?`, phoneID).Scan(&phoneNumber)

		if env == "dev" {
			Urlenv = "http://helpdesk-dev.nopadol.com/tasks/show/" + id
		} else {
			Urlenv = "http://helpdesk.nopadol.com/tasks/show/" + id
		}

		if err == nil {
			// Parse file_paths JSON
			var photoURLs []string
			if filePathsJSON != "" && filePathsJSON != "[]" {
				var filePaths []fiber.Map
				if err := json.Unmarshal([]byte(filePathsJSON), &filePaths); err == nil && len(filePaths) > 0 {
					for _, file := range filePaths {
						if url, ok := file["url"].(string); ok {
							photoURLs = append(photoURLs, url)
						}
					}
				}
			}
			telegramReq := models.TaskRequest{
				PhoneID:        &phoneID,
				SystemID:       systemID,
				IssueElse:      issueElse,
				DepartmentID:   departmentID,
				Text:           text,
				Status:         status,
				ReportedBy:     reportedBy,
				Assignto:       assignto,
				TelegramUser:   telegramUser,
				MessageID:      messageID,
				Ticket:         ticket,
				BranchName:     branchName,
				DepartmentName: departmentName,
				PhoneNumber:    phoneNumber,
				ProgramName:    programName,
				Url:            Urlenv,
				CreatedAt:      createdAt,
				UpdatedAt:      updatedAt,
			}
			var notificationResp int
			if len(photoURLs) > 0 {
				notificationResp, _ = common.UpdateTelegram(telegramReq, photoURLs...)
			} else {
				notificationResp, _ = common.UpdateTelegram(telegramReq)
			}
			_, err = db.DB.Exec(`UPDATE telegram_chat SET assignto_id = ? WHERE id = ?`, notificationResp, telegramID)
			if err != nil {
				log.Printf("❌ Failed to update assignto_id: %v", err)
			} else {
				log.Printf("✅ Assignto ID updated successfully in database")
			}
		}
	}
	return c.JSON(fiber.Map{"success": true, "message": "Assigned person updated successfully"})
}
