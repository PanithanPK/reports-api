package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/url"
	"os"
	"reports-api/db"
	"reports-api/models"
	"reports-api/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func generateticketno(no int) string {
	ticket := fmt.Sprintf("ticket-"+"%04d", no)
	return ticket
}

func deleteImage(objectName string) error {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	endpoint := os.Getenv("End_POINT")
	accessKeyID := os.Getenv("ACCESS_KEY")
	secretAccessKey := os.Getenv("SECRET_ACCESSKEY")
	useSSL := false
	bucketName := os.Getenv("BUCKET_NAME")

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return err
	}

	err = minioClient.RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		log.Printf("Failed to delete %s: %v", objectName, err)
		return err
	}

	log.Printf("Successfully deleted %s", objectName)
	return nil
}

func handleFileUploads(files []*multipart.FileHeader, ticketno string) ([]fiber.Map, []string) {
	var uploadedFiles []fiber.Map
	var errors []string

	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}
	// MinIO configuration
	endpoint := os.Getenv("End_POINT")
	accessKeyID := os.Getenv("ACCESS_KEY")
	secretAccessKey := os.Getenv("SECRET_ACCESSKEY")
	useSSL := false
	bucketName := os.Getenv("BUCKET_NAME")

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Printf("Failed to create MinIO client: %v", err)
		return uploadedFiles, []string{"Failed to initialize storage client"}
	}

	for i, file := range files {
		src, err := file.Open()
		if err != nil {
			log.Printf("Failed to open %s: %v", file.Filename, err)
			errors = append(errors, fmt.Sprintf("Failed to open %s: %v", file.Filename, err))
			continue
		}

		dateStr := time.Now().Add(7 * time.Hour).Format("01022006")
		objectName := fmt.Sprintf("%s-%02d-%s-%s", ticketno, i+1, dateStr, file.Filename)

		_, err = minioClient.PutObject(
			context.Background(),
			bucketName,
			objectName,
			src,
			file.Size,
			minio.PutObjectOptions{ContentType: file.Header.Get("Content-Type")},
		)
		src.Close()

		if err != nil {
			log.Printf("Failed to upload %s: %v", file.Filename, err)
			errors = append(errors, fmt.Sprintf("Failed to upload %s: %v", file.Filename, err))
			continue
		}

		fileURL := fmt.Sprintf("https://minio.sys9.co/api/v1/buckets/%s/objects/download?preview=true&prefix=%s", bucketName, objectName)
		uploadedFiles = append(uploadedFiles, fiber.Map{
			"url": fileURL,
		})
	}

	return uploadedFiles, errors
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
		SELECT t.id, IFNULL(t.ticket_no, ''), IFNULL(t.phone_id, 0) as phone_id, IFNULL(p.number, 0) as number , IFNULL(p.name, '') as phone_name, t.system_id, IFNULL(s.name, '') as system_name, IFNULL(it.name, '') as system_type,
		IFNULL(t.department_id, 0) as department_id, IFNULL(d.name, '') as department_name, IFNULL(d.branch_id, 0) as branch_id, IFNULL(b.name, '') as branch_name,
		t.text, IFNULL(t.assignto, ''), IFNULL(t.reported_by, '') as reported_by, t.status, t.created_at, t.updated_at, IFNULL(t.file_paths, '[]') as file_paths
		FROM tasks t
		LEFT JOIN ip_phones p ON t.phone_id = p.id
		LEFT JOIN departments d ON t.department_id = d.id
		LEFT JOIN branches b ON d.branch_id = b.id
		LEFT JOIN systems_program s ON t.system_id = s.id
		LEFT JOIN issue_types it ON s.type = it.id
		ORDER BY id DESC
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
		var filePathsJSON string
		err := rows.Scan(&t.ID, &t.Ticket, &t.PhoneID, &t.Number, &t.PhoneName, &t.SystemID, &t.SystemName, &t.SystemType, &t.DepartmentID, &t.DepartmentName, &t.BranchID, &t.BranchName, &t.Text, &t.Assignto, &t.ReportedBy, &t.Status, &t.CreatedAt, &t.UpdatedAt, &filePathsJSON)

		if err != nil {
			log.Printf("Error scanning task: %v", err)
			continue
		}

		// Parse file_paths JSON
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
	var uploadedFiles []fiber.Map
	// Get latest ID and add 1 for ticket number
	var lastID int
	err := db.DB.QueryRow("SELECT COALESCE(MAX(id), 0) FROM tasks").Scan(&lastID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get last ID"})
	}

	nextID := lastID + 1
	ticketno := generateticketno(nextID)

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

		uploadedFiles, _ = handleFileUploads(allFiles, ticketno)

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
		if createdByStr := c.FormValue("created_by"); createdByStr != "" {
			req.CreatedBy, _ = strconv.Atoi(createdByStr)
		}
		if telegramStr := c.FormValue("telegram"); telegramStr != "" {
			req.Telegram = telegramStr == "true"
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
		res, err = db.DB.Exec(`INSERT INTO tasks (phone_id, ticket_no, system_id, department_id, text, reported_by, status, created_by, file_paths) VALUES (?, ?, ?, ?, ?, ?, 0, ?, ?)`, req.PhoneID, ticketno, req.SystemID, req.DepartmentID, req.Text, req.ReportedBy, req.CreatedBy, string(filePathsBytes))
	} else {
		res, err = db.DB.Exec(`INSERT INTO tasks (phone_id, ticket_no, system_id, department_id, text, reported_by, status, created_by) VALUES (?, ?, ?, ?, ?, ?, 0, ?)`, req.PhoneID, ticketno, req.SystemID, req.DepartmentID, req.Text, req.ReportedBy, req.CreatedBy)
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
		req.Ticket = ticketno
		req.BranchName = branchName
		req.ProgramName = programName
		req.Url = "http://helpdesk.nopadol.com/"
		req.CreatedAt = time.Now().Add(7 * time.Hour).Format("2006-01-02 15:04:05")
		req.Status = 0

		var messageID int
		// Send with photo if files were uploaded
		if len(uploadedFiles) > 0 {
			// Get first image URL
			if firstFile, ok := uploadedFiles[0]["url"].(string); ok {
				messageID, err = SendTelegram(req, firstFile)
			} else {
				messageID, err = SendTelegram(req)
			}
		} else {
			messageID, err = SendTelegram(req)
		}

		if err == nil {
			// Update task with message_id
			_, err = db.DB.Exec(`UPDATE tasks SET message_id = ? WHERE id = ?`, messageID, id)
			if err != nil {
				log.Printf("Failed to update message_id: %v", err)
			}
		}
	}
	return c.JSON(fiber.Map{"success": true, "id": id})
}

// UpdateTaskHandler แก้ไข task
func UpdateTaskHandler(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	var ticketno string
	var uploadedFiles []fiber.Map
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

	log.Printf("Looking for task with ID: %d", id)
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

		uploadedFiles, _ = handleFileUploads(allFiles, ticketno)
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
	if textStr := c.FormValue("text"); textStr != "" {
		req.Text = textStr
	}
	if reportedByStr := c.FormValue("reported_by"); reportedByStr != "" {
		req.ReportedBy = &reportedByStr
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
	var CreatedAt time.Time
	err = db.DB.QueryRow(`SELECT created_at FROM tasks WHERE id = ?`, id).Scan(&CreatedAt)

	if err != nil {
		log.Println("Error fetching created_at:", err)
	}

	log.Printf("Updating task ID: %s", CreatedAt.Format("2006-01-02 15:04:05"))

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
								deleteImage(objectName)
							}
						}
					}
				}
			}
		}

		// Upload new files
		filePathsBytes, _ := json.Marshal(uploadedFiles)
		log.Printf("Updating file_paths: %s", string(filePathsBytes))

		if req.ReportedBy != nil {
			_, err = db.DB.Exec(`UPDATE tasks SET phone_id=?, system_id=?, department_id=?, assignto=?, reported_by=?, text=?, status=?, updated_at=CURRENT_TIMESTAMP, updated_by=?, file_paths=? WHERE id=?`, req.PhoneID, req.SystemID, req.DepartmentID, req.Assignto, req.ReportedBy, req.Text, req.Status, req.UpdatedBy, string(filePathsBytes), id)
		} else {
			_, err = db.DB.Exec(`UPDATE tasks SET phone_id=?, system_id=?, department_id=?, assignto=?, text=?, status=?, updated_at=CURRENT_TIMESTAMP, updated_by=?, file_paths=? WHERE id=?`, req.PhoneID, req.SystemID, req.DepartmentID, req.Assignto, req.Text, req.Status, req.UpdatedBy, string(filePathsBytes), id)
		}
	} else {
		if req.ReportedBy != nil {
			_, err = db.DB.Exec(`UPDATE tasks SET phone_id=?, system_id=?, department_id=?, assignto=?, reported_by=?, text=?, status=?, updated_at=CURRENT_TIMESTAMP, updated_by=? WHERE id=?`, req.PhoneID, req.SystemID, req.DepartmentID, req.Assignto, req.ReportedBy, req.Text, req.Status, req.UpdatedBy, id)
		} else {
			_, err = db.DB.Exec(`UPDATE tasks SET phone_id=?, system_id=?, department_id=?, assignto=?, text=?, status=?, updated_at=CURRENT_TIMESTAMP, updated_by=? WHERE id=?`, req.PhoneID, req.SystemID, req.DepartmentID, req.Assignto, req.Text, req.Status, req.UpdatedBy, id)
		}
	}
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update task"})
	}

	// Get message_id and update Telegram if exists
	var messageID int
	var reported string
	var existingFilePathsJSON string
	err = db.DB.QueryRow(`SELECT IFNULL(ticket_no, ''),IFNULL(message_id, 0), IFNULL(reported_by, ''), IFNULL(file_paths, '[]') FROM tasks WHERE id = ?`, id).Scan(&ticketno, &messageID, &reported, &existingFilePathsJSON)
	if err == nil && messageID > 0 {
		// Create TaskRequest from TaskRequestUpdate for Telegram
		telegramReq := models.TaskRequest{
			PhoneID:      req.PhoneID,
			SystemID:     req.SystemID,
			DepartmentID: req.DepartmentID,
			Text:         req.Text,
			Status:       req.Status,
			ReportedBy:   reported,
			Assignto:     "",
			MessageID:    messageID,

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
		telegramReq.Url = "http://helpdesk.nopadol.com/"
		telegramReq.CreatedAt = CreatedAt.Add(7 * time.Hour).Format("2006-01-02 15:04:05")
		telegramReq.UpdatedAt = time.Now().Add(7 * time.Hour).Format("2006-01-02 15:04:05")
		telegramReq.Ticket = ticketno
		telegramReq.ReportedBy = reported
		if req.Assignto != nil {
			telegramReq.Assignto = *req.Assignto
		}

		// Get first image URL from existing files for Telegram
		var firstImageURL string
		if existingFilePathsJSON != "" && existingFilePathsJSON != "[]" {
			var existingFiles []fiber.Map
			if err := json.Unmarshal([]byte(existingFilePathsJSON), &existingFiles); err == nil && len(existingFiles) > 0 {
				if url, ok := existingFiles[0]["url"].(string); ok {
					firstImageURL = url
				}
			}
		}

		// Send update to Telegram with photo if available
		if firstImageURL != "" {
			_, _ = UpdateTelegram(telegramReq, firstImageURL)
		} else {
			_, _ = UpdateTelegram(telegramReq)
		}
	}
	log.Printf("Updating task ID: %d", id)
	return c.JSON(fiber.Map{"success": true})
}

// DeleteTaskHandler (soft delete)
func DeleteTaskHandler(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid id"})
	}

	// Get message_id and file_paths before deleting
	var messageID int
	var filePathsJSON string
	err = db.DB.QueryRow(`SELECT IFNULL(message_id, 0), IFNULL(file_paths, '[]') FROM tasks WHERE id = ?`, id).Scan(&messageID, &filePathsJSON)
	if err != nil {
		log.Printf("Failed to get task data: %v", err)
	}

	// Delete files from MinIO if they exist
	if filePathsJSON != "" && filePathsJSON != "[]" {
		var filePaths []fiber.Map
		if err := json.Unmarshal([]byte(filePathsJSON), &filePaths); err == nil {
			for _, fp := range filePaths {
				if url, ok := fp["url"].(string); ok {
					if strings.Contains(url, "prefix=") {
						parts := strings.Split(url, "prefix=")
						if len(parts) > 1 {
							objectName := parts[1]
							deleteImage(objectName)
						}
					}
				}
			}
		}
	}

	_, err = db.DB.Exec(`DELETE FROM tasks WHERE id=?`, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete task"})
	}

	// Delete Telegram message if exists
	if messageID > 0 {
		_, _ = DeleteTelegram(messageID)
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
	var filePathsJSON string
	var task models.TaskWithDetails
	err = db.DB.QueryRow(`
		SELECT t.id, IFNULL(t.ticket_no, '') as ticket_to, IFNULL(t.phone_id, 0) as phone_id, IFNULL(p.number, 0) as number, IFNULL(p.name, '') as phone_name, IFNULL(t.system_id, 0) as system_id, IFNULL(s.name, '') as system_name, IFNULL(it.name, '') as system_type,
		IFNULL(t.department_id, 0) as department_id, IFNULL(d.name, '') as department_name, IFNULL(d.branch_id, 0) as branch_id, IFNULL(b.name, '') as branch_name,
		IFNULL(t.text, '') as text, IFNULL(t.assignto, ''), IFNULL(t.reported_by, '') as reported_by, IFNULL(t.status, 0) as status, IFNULL(t.created_at, '') as created_at, IFNULL(t.updated_at, '') as updated_at, IFNULL(t.file_paths, '[]') as file_paths
		FROM tasks t
		LEFT JOIN ip_phones p ON t.phone_id = p.id
		LEFT JOIN departments d ON t.department_id = d.id
		LEFT JOIN branches b ON d.branch_id = b.id
		LEFT JOIN systems_program s ON t.system_id = s.id
		LEFT JOIN issue_types it ON s.type = it.id
		WHERE t.id = ?
	`, id).Scan(&task.ID, &task.Ticket, &task.PhoneID, &task.Number, &task.PhoneName, &task.SystemID, &task.SystemName, &task.SystemType, &task.DepartmentID, &task.DepartmentName, &task.BranchID, &task.BranchName, &task.Text, &task.Assignto, &task.ReportedBy, &task.Status, &task.CreatedAt, &task.UpdatedAt, &filePathsJSON)

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

// SearchTasksHandler returns a handler for searching tasks
func SearchTasksHandler(c *fiber.Ctx) error {
	query := c.Params("query")
	if query == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Query parameter is required"})
	}

	// URL decode for Thai language support
	decodedQuery, err := url.QueryUnescape(query)
	if err != nil {
		decodedQuery = query
	}
	decodedQuery = strings.TrimSpace(decodedQuery)
	searchPattern := "%" + decodedQuery + "%"

	// Get total count
	var total int
	err = db.DB.QueryRow(`
		SELECT COUNT(*) FROM tasks t
		LEFT JOIN ip_phones p ON t.phone_id = p.id AND p.deleted_at IS NULL
		LEFT JOIN departments d ON t.department_id = d.id AND d.deleted_at IS NULL
		LEFT JOIN branches b ON d.branch_id = b.id AND b.deleted_at IS NULL
		LEFT JOIN systems_program s ON t.system_id = s.id AND s.deleted_at IS NULL
		WHERE (t.ticket_no LIKE ? OR p.number LIKE ? OR p.name LIKE ? OR d.name LIKE ? OR b.name LIKE ? OR s.name LIKE ? OR t.text LIKE ?)
	`, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern).Scan(&total)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to count search results"})
	}

	// Get search results
	rows, err := db.DB.Query(`
		SELECT t.id, IFNULL(t.ticket_no, ''), IFNULL(t.phone_id, 0) as phone_id, IFNULL(p.number, 0) as number, IFNULL(p.name, '') as phone_name, t.system_id, IFNULL(s.name, '') as system_name,
		IFNULL(t.department_id, 0) as department_id, IFNULL(d.name, '') as department_name, IFNULL(d.branch_id, 0) as branch_id, IFNULL(b.name, '') as branch_name,
		t.text, t.status, t.created_at, t.updated_at
		FROM tasks t
		LEFT JOIN ip_phones p ON t.phone_id = p.id AND p.deleted_at IS NULL
		LEFT JOIN departments d ON t.department_id = d.id AND d.deleted_at IS NULL
		LEFT JOIN branches b ON d.branch_id = b.id AND b.deleted_at IS NULL
		LEFT JOIN systems_program s ON t.system_id = s.id AND s.deleted_at IS NULL
		WHERE (t.ticket_no LIKE ? OR p.number LIKE ? OR p.name LIKE ? OR d.name LIKE ? OR b.name LIKE ? OR s.name LIKE ? OR t.text LIKE ?)
		ORDER BY t.id DESC
	`, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to search tasks"})
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

	log.Printf("Searching tasks with query: %s", query)
	return c.JSON(models.PaginatedResponse{
		Success: true,
		Data:    tasks,
	})
}

// GetTasksWithQueryHandler handles both list all and search functionality
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

	// Clean query
	decodedQuery = strings.TrimSpace(decodedQuery)
	decodedQuery = strings.ReplaceAll(decodedQuery, "  ", " ")
	decodedQuery = strings.ReplaceAll(decodedQuery, "%", "")
	decodedQuery = strings.ReplaceAll(decodedQuery, "_", "")
	decodedQuery = strings.ReplaceAll(decodedQuery, "'", "")
	decodedQuery = strings.ReplaceAll(decodedQuery, "\"", "")

	searchPattern := "%" + decodedQuery + "%"

	// Get total count
	var total int
	err = db.DB.QueryRow(`
		SELECT COUNT(*) FROM tasks t
		LEFT JOIN ip_phones p ON t.phone_id = p.id
		LEFT JOIN departments d ON t.department_id = d.id
		LEFT JOIN branches b ON d.branch_id = b.id
		LEFT JOIN systems_program s ON t.system_id = s.id
		WHERE (t.ticket_no LIKE ? OR p.number LIKE ? OR p.name LIKE ? OR d.name LIKE ? OR b.name LIKE ? OR s.name LIKE ? OR t.text LIKE ?)
	`, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern).Scan(&total)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to count search results"})
	}

	// Get search results with pagination
	rows, err := db.DB.Query(`
		SELECT t.id, IFNULL(t.ticket_no, ''), IFNULL(t.phone_id, 0) as phone_id, IFNULL(p.number, 0) as number, IFNULL(p.name, '') as phone_name, t.system_id, IFNULL(s.name, '') as system_name,
		IFNULL(t.department_id, 0) as department_id, IFNULL(d.name, '') as department_name, IFNULL(d.branch_id, 0) as branch_id, IFNULL(b.name, '') as branch_name,
		t.text, IFNULL(t.assignto, ''), t.status, t.created_at, t.updated_at
		FROM tasks t
		LEFT JOIN ip_phones p ON t.phone_id = p.id
		LEFT JOIN departments d ON t.department_id = d.id
		LEFT JOIN branches b ON d.branch_id = b.id
		LEFT JOIN systems_program s ON t.system_id = s.id
		WHERE (t.ticket_no LIKE ? OR p.number LIKE ? OR p.name LIKE ? OR d.name LIKE ? OR b.name LIKE ? OR s.name LIKE ? OR t.text LIKE ? OR t.assignto LIKE ?)
		ORDER BY t.id DESC
		LIMIT ? OFFSET ?
	`, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, pagination.Limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to search tasks"})
	}
	defer rows.Close()

	var tasks []models.TaskWithDetails
	for rows.Next() {
		var t models.TaskWithDetails
		err := rows.Scan(&t.ID, &t.Ticket, &t.PhoneID, &t.Number, &t.PhoneName, &t.SystemID, &t.SystemName, &t.DepartmentID, &t.DepartmentName, &t.BranchID, &t.BranchName, &t.Text, &t.Assignto, &t.Status, &t.CreatedAt, &t.UpdatedAt)
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
