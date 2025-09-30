package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"reports-api/config"
	"reports-api/db"
	"reports-api/handlers/common"
	"reports-api/models"
	"slices"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// parseProgressText parses progress_text and populates ProgressEntry text field
// Now progress_text contains only text, files are stored separately in file_paths
func parseProgressText(progressText string, entry *models.ProgressEntry) {
	var progressData map[string]any
	if err := json.Unmarshal([]byte(progressText), &progressData); err == nil {
		// Check if it's old format with files included
		if text, ok := progressData["text"].(string); ok {
			entry.Text = text
		}
		// Handle legacy format that might still have files in progress_text
		if fileList, ok := progressData["files"].([]any); ok {
			if entry.FilePaths == nil {
				entry.FilePaths = make(map[string]string)
			}
			for i, file := range fileList {
				if fileMap, ok := file.(map[string]any); ok {
					if url, ok := fileMap["url"].(string); ok {
						entry.FilePaths[strconv.Itoa(i)] = url
					}
				}
			}
		}
	} else {
		// Not JSON, treat as plain text
		entry.Text = progressText
	}
}

// parseProgressFilePaths parses file_paths JSON array and populates ProgressEntry FilePaths
func parseProgressFilePaths(filePathsJSON sql.NullString, entry *models.ProgressEntry) {
	if filePathsJSON.Valid && filePathsJSON.String != "" {
		// Try to parse as array of objects [{"url": "..."}]
		var fileObjects []map[string]any
		if err := json.Unmarshal([]byte(filePathsJSON.String), &fileObjects); err == nil {
			if entry.FilePaths == nil {
				entry.FilePaths = make(map[string]string)
			}
			for i, fileObj := range fileObjects {
				if url, ok := fileObj["url"].(string); ok {
					entry.FilePaths[fmt.Sprintf("image_%d", i)] = url
				}
			}
		} else {
			// Fallback: try to parse as array of strings (legacy format)
			var fileURLs []string
			if err := json.Unmarshal([]byte(filePathsJSON.String), &fileURLs); err == nil {
				if entry.FilePaths == nil {
					entry.FilePaths = make(map[string]string)
				}
				for i, url := range fileURLs {
					entry.FilePaths[fmt.Sprintf("image_%d", i)] = url
				}
			}
		}
	}
}

// parseProgressTextAndDeleteFiles parses progress_text and deletes files from MinIO (legacy format)
func parseProgressTextAndDeleteFiles(progressText string, entry *models.ProgressEntry) {
	var progressData map[string]any
	if err := json.Unmarshal([]byte(progressText), &progressData); err == nil {
		// Successfully parsed as JSON, extract text and files
		if text, ok := progressData["text"].(string); ok {
			entry.Text = text
		}
		if fileList, ok := progressData["files"].([]any); ok {
			// Convert files to map[string]string format expected by FilePaths
			entry.FilePaths = make(map[string]string)

			// ลบไฟล์จาก MinIO และ populate FilePaths
			for i, file := range fileList {
				if fileMap, ok := file.(map[string]any); ok {
					if url, ok := fileMap["url"].(string); ok {
						// Add to FilePaths for response
						entry.FilePaths[strconv.Itoa(i)] = url

						// Delete from MinIO
						if strings.Contains(url, "prefix=") {
							parts := strings.Split(url, "prefix=")
							if len(parts) > 1 {
								objectName := parts[1]
								if deleteErr := common.DeleteImage(objectName); deleteErr != nil {
									log.Printf("Warning: Failed to delete file from MinIO: %s, error: %v", objectName, deleteErr)
								} else {
									log.Printf("Successfully deleted file from MinIO: %s", objectName)
								}
							}
						}
					}
				}
			}
		}
	} else {
		// Not JSON, treat as plain text
		entry.Text = progressText
	}
}

// parseFilePathsAndDelete parses file_paths from database and deletes files from MinIO
func parseFilePathsAndDelete(filePathsJSON sql.NullString, entry *models.ProgressEntry) {
	if filePathsJSON.Valid && filePathsJSON.String != "" {
		// Try to parse as array of objects [{"url": "..."}]
		var fileObjects []map[string]any
		if err := json.Unmarshal([]byte(filePathsJSON.String), &fileObjects); err == nil {
			if entry.FilePaths == nil {
				entry.FilePaths = make(map[string]string)
			}

			for i, fileObj := range fileObjects {
				if url, ok := fileObj["url"].(string); ok {
					// Add to FilePaths for response
					entry.FilePaths[strconv.Itoa(i)] = url

					// Delete from MinIO
					if strings.Contains(url, "prefix=") {
						parts := strings.Split(url, "prefix=")
						if len(parts) > 1 {
							objectName := parts[1]
							if deleteErr := common.DeleteImage(objectName); deleteErr != nil {
								log.Printf("Warning: Failed to delete file from MinIO: %s, error: %v", objectName, deleteErr)
							} else {
								log.Printf("Successfully deleted file from MinIO: %s", objectName)
							}
						}
					}
				}
			}
		} else {
			// Fallback: try to parse as array of strings (legacy format)
			var fileURLs []string
			if err := json.Unmarshal([]byte(filePathsJSON.String), &fileURLs); err == nil {
				if entry.FilePaths == nil {
					entry.FilePaths = make(map[string]string)
				}

				for i, url := range fileURLs {
					// Add to FilePaths for response
					entry.FilePaths[strconv.Itoa(i)] = url

					// Delete from MinIO
					if strings.Contains(url, "prefix=") {
						parts := strings.Split(url, "prefix=")
						if len(parts) > 1 {
							objectName := parts[1]
							if deleteErr := common.DeleteImage(objectName); deleteErr != nil {
								log.Printf("Warning: Failed to delete file from MinIO: %s, error: %v", objectName, deleteErr)
							} else {
								log.Printf("Successfully deleted file from MinIO: %s", objectName)
							}
						}
					}
				}
			} else {
				log.Printf("Warning: Failed to parse file_paths JSON: %v", err)
			}
		}
	}
}

// ProgressEntry represents a single progress entry

// @Summary Create progress entry for task
// @Description Add a new progress entry to a specific task
// @Tags progress
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Task ID"
// @Param text formData string true "Progress text"
// @Param image formData file false "Progress image files"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/progress/create/{id} [post]
func CreateProgressHandler(c *fiber.Ctx) error {
	idStr := c.Params("id")
	taskID, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid task ID"})
	}

	// ตรวจสอบว่า task_id มีอยู่ในตาราง tasks หรือไม่ และดึง ticket_no และ status
	var ticketNo string
	var status int
	var Urlenv string
	env := config.AppConfig.Environment

	err = db.DB.QueryRow("SELECT IFNULL(ticket_no, ''), IFNULL(status, 0) FROM tasks WHERE id = ?", taskID).Scan(&ticketNo, &status)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{"error": "Task not found"})
		}
		log.Printf("Error checking task existence: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}

	// ตรวจสอบว่า task เสร็จสิ้นแล้วหรือไม่ (status = 2)
	if status == 2 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Cannot add progress to completed task",
		})
	}

	if status != 2 {
		_, err = db.DB.Exec(`UPDATE tasks SET status = 1 WHERE id = ?`, idStr)
		if err != nil {
			log.Printf("Database error: %v", err)
			return c.Status(500).JSON(fiber.Map{"error": "Failed to update status progress"})
		}
	}

	var uploadedFiles []fiber.Map
	var progressText string

	// Try to parse as multipart form first (for file uploads)
	form, err := c.MultipartForm()
	if err != nil {
		// If multipart parsing fails, get text from regular form or JSON body
		progressText = c.FormValue("text")
		if progressText == "" {
			// Try to get from JSON body
			var reqBody map[string]interface{}
			if err := c.BodyParser(&reqBody); err == nil {
				if text, ok := reqBody["text"].(string); ok {
					progressText = text
				}
			}
		}
	} else {
		// Handle multipart form data
		progressText = c.FormValue("text")

		// Handle file uploads if present
		var allFiles []*multipart.FileHeader

		// Check for indexed files (image_0, image_1, image_2, etc.) or single image field
		for key, files := range form.File {
			if strings.HasPrefix(key, "image_") || key == "image" {
				allFiles = append(allFiles, files...)
			}
		}

		// Upload files if any were provided
		log.Printf("progress ticketNo : %s", idStr)
		log.Printf("progress ticketNo : %s", ticketNo)
		if len(allFiles) > 0 {
			uploadedFiles, _ = common.HandleFileUploadsProgress(allFiles, ticketNo)
			log.Printf("Uploaded %d files for progress entry", len(uploadedFiles))
		}
	}

	// Validate that we have progress text
	if progressText == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Progress text is required"})
	}

	// Prepare file paths JSON if files were uploaded
	var filePathsJSON string
	if len(uploadedFiles) > 0 {
		// Keep the original format with url objects: [{"url": "..."}]
		filePathsBytes, _ := json.Marshal(uploadedFiles)
		filePathsJSON = string(filePathsBytes)
	}

	// บันทึกข้อมูลลงในตาราง progress
	var result sql.Result
	if filePathsJSON != "" {
		result, err = db.DB.Exec(
			"INSERT INTO progress (task_id, progress_text, file_paths, created_at, updated_at) VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)",
			taskID, progressText, filePathsJSON,
		)
	} else {
		result, err = db.DB.Exec(
			"INSERT INTO progress (task_id, progress_text, created_at, updated_at) VALUES (?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)",
			taskID, progressText,
		)
	}
	if err != nil {
		log.Printf("Error inserting progress: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create progress entry"})
	}

	progressID, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert ID: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get progress ID"})
	}

	log.Printf("Created progress entry with ID: %d for task ID: %d", progressID, taskID)

	// ดึงข้อมูลที่เพิ่งสร้างเพื่อส่งกลับในรูปแบบ ProgressEntry
	var createdEntry models.ProgressEntry
	var createdAt, updatedAt string
	var retrievedFilePathsJSON sql.NullString
	err = db.DB.QueryRow(
		"SELECT id, progress_text, file_paths, created_at, updated_at FROM progress WHERE id = ?",
		progressID,
	).Scan(&createdEntry.ID, &createdEntry.Text, &retrievedFilePathsJSON, &createdAt, &updatedAt)

	if err != nil {
		log.Printf("Error retrieving created progress entry: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve created progress entry"})
	}

	createdEntry.CreatedAt = createdAt
	createdEntry.UpdateAt = updatedAt

	// Parse file_paths จากฐานข้อมูล (ถ้ามี)
	parseProgressFilePaths(retrievedFilePathsJSON, &createdEntry)

	// ดึงข้อมูล task สำหรับอัพเดต Telegram
	var ticket, text, issueElse, reportedBy, assignto, branchName, departmentName, programName string
	var phoneID, systemID, departmentID, messageID, phoneNumber, branchID, telegramID int
	var phoneElse *string
	var telegramUser string
	// Fix SQL: JOINs before WHERE, select tc.report_id as messageID
	err = db.DB.QueryRow(`
			SELECT IFNULL(t.ticket_no, ''), IFNULL(t.phone_id, 0), IFNULL(t.phone_else, ''), IFNULL(t.system_id, 0), IFNULL(t.issue_else, ''), IFNULL(t.department_id, 0),
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
		`, idStr).Scan(&ticket, &phoneID, &phoneElse, &systemID, &issueElse, &departmentID, &text, &status, &reportedBy, &assignto, &telegramUser, &messageID, &filePathsJSON, &branchID, &createdAt, &updatedAt, &telegramID)
	log.Printf("Fetched task for Telegram update, ID: %s, MessageID: %d", telegramUser, messageID)
	// Query extra info for Telegram
	db.DB.QueryRow(`SELECT name FROM branches WHERE id = ?`, branchID).Scan(&branchName)
	db.DB.QueryRow(`SELECT name FROM departments WHERE id = ?`, departmentID).Scan(&departmentName)
	db.DB.QueryRow(`SELECT name FROM systems_program WHERE id = ?`, systemID).Scan(&programName)
	db.DB.QueryRow(`SELECT number FROM ip_phones WHERE id = ?`, phoneID).Scan(&phoneNumber)

	if env == "dev" {
		Urlenv = "http://helpdesk-dev.nopadol.com/tasks/show/" + idStr
	} else {
		Urlenv = "http://helpdesk.nopadol.com/tasks/show/" + idStr
	}

	CreatedAt := common.Fixtimefeature(createdAt)
	UpdatedAt := common.Fixtimefeature(updatedAt)

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
			PhoneElse:      phoneElse,
			SystemID:       systemID,
			IssueElse:      issueElse,
			DepartmentID:   departmentID,
			Text:           text,
			Status:         1,
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
			CreatedAt:      CreatedAt,
			UpdatedAt:      UpdatedAt,
		}
		if len(photoURLs) > 0 {
			assigntoID, _ := common.UpdateTelegram(telegramReq, photoURLs...)
			_, err = db.DB.Exec(`UPDATE telegram_chat SET assignto_id = ? WHERE id = ?`, assigntoID, idStr)
			if err != nil {
				log.Printf("Database error: %v", err)
				return c.Status(500).JSON(fiber.Map{"error": "Failed to update telegram chat"})
			}
		} else {
			assigntoID, _ := common.UpdateTelegram(telegramReq)
			_, err = db.DB.Exec(`UPDATE telegram_chat SET assignto_id = ? WHERE id = ?`, assigntoID, idStr)
			if err != nil {
				log.Printf("Database error: %v", err)
				return c.Status(500).JSON(fiber.Map{"error": "Failed to update telegram chat"})
			}
		}

	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Progress entry created successfully",
	})
}

// @Summary Get progress entries for task
// @Description Get all progress entries for a specific task
// @Tags progress
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/progress/{id} [get]
func GetProgressHandler(c *fiber.Ctx) error {
	idStr := c.Params("id")
	taskID, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid task ID"})
	}

	// ตรวจสอบว่า task_id มีอยู่ในตาราง tasks หรือไม่
	var taskExists string
	var assignto string
	var ticketno string
	err = db.DB.QueryRow("SELECT '1', IFNULL(ticket_no, ''),IFNULL(assignto,'') FROM tasks WHERE id = ? LIMIT 1", taskID).Scan(&taskExists, &ticketno, &assignto)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{"error": "Task not found"})
		}
		log.Printf("Error checking task existence: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}

	// ดึงข้อมูล progress entries สำหรับ task นี้
	rows, err := db.DB.Query(
		"SELECT id, progress_text, file_paths, created_at, updated_at FROM progress WHERE task_id = ?",
		taskID,
	)
	if err != nil {
		log.Printf("Error querying progress entries: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get progress entries"})
	}
	defer rows.Close()

	var progressEntries []models.ProgressEntry

	for rows.Next() {
		var entry models.ProgressEntry
		var progressText string
		var filePathsJSON sql.NullString
		var CreatedAt string
		var UpdateAt string

		err := rows.Scan(&entry.ID, &progressText, &filePathsJSON, &CreatedAt, &UpdateAt)
		if err != nil {
			log.Printf("Error scanning progress row: %v", err)
			continue
		}

		// Parse progress_text and file_paths
		parseProgressText(progressText, &entry)
		parseProgressFilePaths(filePathsJSON, &entry)

		// Set assignto from task
		entry.CreatedAt = common.Fixtimefeature(CreatedAt)
		entry.UpdateAt = common.Fixtimefeature(UpdateAt)
		entry.Ticketno = ticketno
		entry.AssignTo = assignto

		progressEntries = append(progressEntries, entry)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Row error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to read progress entries"})
	}

	log.Printf("Retrieved %d progress entries for task ID: %d", len(progressEntries), taskID)

	return c.JSON(fiber.Map{
		"success": true,
		"data":    progressEntries,
		"count":   len(progressEntries),
	})
}

// @Summary Update progress entry
// @Description Update an existing progress entry with new text and/or images
// @Tags progress
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Task ID"
// @Param pgid path string true "Progress ID"
// @Param text formData string false "Updated progress text"
// @Param image_urls formData string false "JSON array of image URLs to keep"
// @Param image formData file false "New image files"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/progress/update/{id}/{pgid} [put]
func UpdateProgressHandler(c *fiber.Ctx) error {
	idStr := c.Params("id")

	progressid := c.Params("pgid")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid task ID"})
	}
	var req models.UpdateProgress
	var existingProgressText string
	var uploadedFiles []fiber.Map
	var ticketno string
	var existingFilePathsJSON sql.NullString
	var taskid int
	err = db.DB.QueryRow(`
		SELECT task_id, progress_text, file_paths
		FROM progress WHERE id = ?
	`, progressid).Scan(&taskid, &existingProgressText, &existingFilePathsJSON)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Progress not found"})
	}

	if taskid != id {
		return c.Status(404).JSON(fiber.Map{"error": "tasks id not match"})
	}

	log.Printf("Looking for task with ID: %d", id)
	err = db.DB.QueryRow("SELECT ticket_no FROM tasks WHERE id = ?", id).Scan(&ticketno)
	if err != nil {
		log.Printf("Task not found error: %v", err)
		return c.Status(404).JSON(fiber.Map{"error": "Task not found"})
	}
	log.Printf("Found task with ticket_no: %s", ticketno)

	var keepImageURLs []string
	form, err := c.MultipartForm()
	if err != nil {
		// Handle JSON request
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request format"})
		}
		keepImageURLs = req.ImageURLs
	} else {
		// Handle multipart form request
		req.Text = c.FormValue("text")

		// รับ URL รูปเก่าที่ต้องการเก็บไว้
		imageURLsStr := c.FormValue("image_urls")
		if imageURLsStr != "" {
			if err := json.Unmarshal([]byte(imageURLsStr), &keepImageURLs); err != nil {
				log.Printf("Error parsing image_urls: %v", err)
			}
		}

		// จัดการไฟล์ใหม่
		var allFiles []*multipart.FileHeader
		for key, files := range form.File {
			if strings.HasPrefix(key, "image_") || key == "image" {
				allFiles = append(allFiles, files...)
			}
		}

		// ตรวจสอบว่า ImageURLs ที่ส่งมาตรงกับที่มีอยู่แล้วหรือไม่
		if len(allFiles) == 0 && len(keepImageURLs) > 0 {
			var existingURLs []string

			// ดึง URLs จาก file_paths column
			if existingFilePathsJSON.Valid && existingFilePathsJSON.String != "" {
				// Try to parse as array of objects [{"url": "..."}]
				var fileObjects []map[string]any
				if err := json.Unmarshal([]byte(existingFilePathsJSON.String), &fileObjects); err == nil {
					for _, fileObj := range fileObjects {
						if url, ok := fileObj["url"].(string); ok {
							existingURLs = append(existingURLs, url)
						}
					}
				} else {
					// Fallback: try to parse as array of strings
					if err := json.Unmarshal([]byte(existingFilePathsJSON.String), &existingURLs); err != nil {
						log.Printf("Error parsing existing file_paths: %v", err)
					}
				}
			}

			// ถ้า URLs ตรงกันทั้งหมด และไม่มีไฟล์ใหม่
			if len(existingURLs) == len(keepImageURLs) {
				allMatch := true
				for _, keepURL := range keepImageURLs {
					if !slices.Contains(existingURLs, keepURL) {
						allMatch = false
						break
					}
				}
				if allMatch {
					// URLs ตรงกันทั้งหมด ใช้ text เดิมถ้าไม่ได้ส่งมาใหม่
					if req.Text == "" {
						req.Text = existingProgressText
					}
					// อัปเดตเฉพาะ progress_text
					_, err = db.DB.Exec(`UPDATE progress SET progress_text = ? WHERE id = ?`, req.Text, progressid)
					if err != nil {
						return c.Status(500).JSON(fiber.Map{"error": "Failed to update progress"})
					}

					return c.JSON(fiber.Map{"success": true, "message": "Progress updated successfully"})
				}
			}
		}

		// ลบรูปเก่าที่ไม่ต้องการเก็บไว้
		if existingFilePathsJSON.Valid && existingFilePathsJSON.String != "" {
			var existingURLs []string

			// Try to parse as array of objects [{"url": "..."}]
			var fileObjects []map[string]any
			if err := json.Unmarshal([]byte(existingFilePathsJSON.String), &fileObjects); err == nil {
				for _, fileObj := range fileObjects {
					if url, ok := fileObj["url"].(string); ok {
						existingURLs = append(existingURLs, url)
					}
				}
			} else {
				// Fallback: try to parse as array of strings
				if err := json.Unmarshal([]byte(existingFilePathsJSON.String), &existingURLs); err != nil {
					log.Printf("Error parsing existing file_paths for deletion: %v", err)
				}
			}

			for _, url := range existingURLs {
				// ถ้าไม่ได้ส่ง image_urls มา ให้ลบทั้งหมด
				if len(keepImageURLs) == 0 {
					if strings.Contains(url, "prefix=") {
						parts := strings.Split(url, "prefix=")
						if len(parts) > 1 {
							objectName := parts[1]
							common.DeleteImage(objectName)
						}
					}
				} else {
					// ตรวจสอบว่า URL นี้อยู่ในรายการที่ต้องการเก็บไว้หรือไม่
					if !slices.Contains(keepImageURLs, url) {
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

		// อัปโหลดไฟล์ใหม่ถ้ามี
		if len(allFiles) > 0 {
			var uploadErrors []string
			uploadedFiles, uploadErrors = common.HandleFileUploadsProgress(allFiles, ticketno)
			if len(uploadErrors) > 0 {
				log.Printf("File upload errors: %v", uploadErrors)
				// Continue with partial success, but log the errors
			}
		}

		// รวมรูปเก่าที่เก็บไว้กับรูปใหม่
		for _, keepURL := range keepImageURLs {
			uploadedFiles = append(uploadedFiles, fiber.Map{"url": keepURL})
		}
	}

	// ใช้ text เดิมถ้าไม่ได้ส่งมาใหม่
	if req.Text == "" {
		req.Text = existingProgressText
	}

	// เตรียม file paths JSON ใหม่
	var newFilePathsJSON string
	if len(uploadedFiles) > 0 {
		filePathsBytes, _ := json.Marshal(uploadedFiles)
		newFilePathsJSON = string(filePathsBytes)
	}

	// อัปเดต progress
	if newFilePathsJSON != "" {
		_, err = db.DB.Exec(`UPDATE progress SET progress_text = ?, file_paths = ? WHERE id = ?`, req.Text, newFilePathsJSON, progressid)
	} else {
		_, err = db.DB.Exec(`UPDATE progress SET progress_text = ?, file_paths = NULL WHERE id = ?`, req.Text, progressid)
	}
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update progress"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Progress updated successfully",
	})

}

// @Summary Delete progress entry
// @Description Delete a specific progress entry and its associated files permanently
// @Tags progress
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Param pgid path string true "Progress ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/progress/delete/{id}/{pgid} [delete]
func DeleteProgressHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	progressid := c.Params("pgid")

	// Validate task ID
	taskID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid task ID"})
	}

	// Validate progress ID
	progressID, err := strconv.Atoi(progressid)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid progress ID"})
	}

	// ตรวจสอบว่า task มีอยู่จริงหรือไม่
	var taskExists string
	err = db.DB.QueryRow("SELECT '1' FROM tasks WHERE id = ? LIMIT 1", taskID).Scan(&taskExists)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{"error": "Task not found"})
		}
		log.Printf("Error checking task existence: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}

	// ดึงข้อมูล progress entry ที่จะลบ โดยใช้ models.ProgressEntry
	var progressEntry models.ProgressEntry
	var progressText string
	var filePathsJSON sql.NullString
	err = db.DB.QueryRow(
		"SELECT id, progress_text, file_paths, created_at, updated_at FROM progress WHERE id = ? AND task_id = ?",
		progressID, taskID,
	).Scan(&progressEntry.ID, &progressText, &filePathsJSON, &progressEntry.CreatedAt, &progressEntry.UpdateAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{"error": "Progress entry not found"})
		}
		log.Printf("Error retrieving progress entry: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}

	// Parse progress text เพื่อดึงข้อมูลไฟล์และ populate ProgressEntry
	parseProgressTextAndDeleteFiles(progressText, &progressEntry)

	// Parse file_paths จากฐานข้อมูล (ถ้ามี)
	parseFilePathsAndDelete(filePathsJSON, &progressEntry)

	// ทำ hard delete - ลบข้อมูลออกจากฐานข้อมูลจริง
	_, err = db.DB.Exec(
		"DELETE FROM progress WHERE id = ? AND task_id = ?",
		progressID, taskID,
	)
	if err != nil {
		log.Printf("Error deleting progress entry: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete progress entry"})
	}

	// Count deleted files for logging
	fileCount := 0
	if progressEntry.FilePaths != nil {
		fileCount = len(progressEntry.FilePaths)
	}

	log.Printf("Successfully deleted progress entry ID: %d for task ID: %d with %d associated files", progressID, taskID, fileCount)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Progress entry deleted successfully",
	})
}
