package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"mime/multipart"
	"reports-api/db"
	"reports-api/handlers/common"
	"reports-api/models"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

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

	// ตรวจสอบว่า task_id มีอยู่ในตาราง tasks หรือไม่ และดึง ticket_no
	var ticketNo string
	err = db.DB.QueryRow("SELECT IFNULL(ticket_no, '') FROM tasks WHERE id = ? AND deleted_at IS NULL", taskID).Scan(&ticketNo)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{"error": "Task not found"})
		}
		log.Printf("Error checking task existence: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
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
		if len(allFiles) > 0 {
			uploadedFiles, _ = common.HandleFileUploadsProgress(allFiles, ticketNo)
			log.Printf("Uploaded %d files for progress entry", len(uploadedFiles))
		}
	}

	// Validate that we have progress text
	if progressText == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Progress text is required"})
	}

	// Prepare the final progress text with file information if files were uploaded
	finalProgressText := progressText
	if len(uploadedFiles) > 0 {
		// Create a structured progress entry with text and files
		progressData := fiber.Map{
			"text":  progressText,
			"files": uploadedFiles,
		}
		progressDataBytes, _ := json.Marshal(progressData)
		finalProgressText = string(progressDataBytes)
	}

	// บันทึกข้อมูลลงในตาราง progress
	result, err := db.DB.Exec(
		"INSERT INTO progress (task_id, progress_text, created_at, updated_at) VALUES (?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)",
		taskID, finalProgressText,
	)
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
	err = db.DB.QueryRow(
		"SELECT id, progress_text, created_at, updated_at FROM progress WHERE id = ?",
		progressID,
	).Scan(&createdEntry.ID, &finalProgressText, &createdAt, &updatedAt)

	if err != nil {
		log.Printf("Error retrieving created progress entry: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve created progress entry"})
	}

	createdEntry.CreatedAt = createdAt
	createdEntry.UpdateAt = updatedAt

	// Parse the stored progress text to populate the model fields
	var progressData map[string]interface{}
	if err := json.Unmarshal([]byte(finalProgressText), &progressData); err == nil {
		// Successfully parsed as JSON, extract text and files
		if text, ok := progressData["text"].(string); ok {
			createdEntry.Text = text
		}
		if fileList, ok := progressData["files"].([]interface{}); ok {
			// Convert files to map[string]string format expected by FilePaths
			createdEntry.FilePaths = make(map[string]string)
			for i, file := range fileList {
				if fileMap, ok := file.(map[string]interface{}); ok {
					if url, ok := fileMap["url"].(string); ok {
						createdEntry.FilePaths[strconv.Itoa(i)] = url
					}
				}
			}
		}
	} else {
		// Not JSON, treat as plain text
		createdEntry.Text = finalProgressText
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Progress entry created successfully",
		"data":    createdEntry,
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
	err = db.DB.QueryRow("SELECT '1' FROM tasks WHERE id = ? AND deleted_at IS NULL LIMIT 1", taskID).Scan(&taskExists)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{"error": "Task not found"})
		}
		log.Printf("Error checking task existence: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}

	// ดึงข้อมูล progress entries สำหรับ task นี้
	rows, err := db.DB.Query(
		"SELECT id, progress_text, created_at, updated_at FROM progress WHERE task_id = ? AND deleted_at IS NULL ORDER BY created_at DESC",
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

		err := rows.Scan(&entry.ID, &progressText, &entry.CreatedAt, &entry.UpdateAt)
		if err != nil {
			log.Printf("Error scanning progress row: %v", err)
			continue
		}

		// Try to parse progress_text as JSON to extract text and files
		var progressData map[string]interface{}

		if err := json.Unmarshal([]byte(progressText), &progressData); err == nil {
			// Successfully parsed as JSON, extract text and files
			if text, ok := progressData["text"].(string); ok {
				entry.Text = text
			}
			if fileList, ok := progressData["files"].([]interface{}); ok {
				// Convert files to map[string]string format expected by FilePaths
				entry.FilePaths = make(map[string]string)
				for i, file := range fileList {
					if fileMap, ok := file.(map[string]interface{}); ok {
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
	id := c.Params("id")
	progressid := c.Params("pgid")
	var req models.ProgressEntry
	var existingProgressText string
	var uploadedFiles []fiber.Map
	var ticketno string

	log.Printf("Looking for task with ID: %s", id)
	err := db.DB.QueryRow("SELECT ticket_no FROM tasks WHERE id = ?", id).Scan(&ticketno)
	if err != nil {
		log.Printf("Task not found error: %v", err)
		return c.Status(404).JSON(fiber.Map{"error": "Task not found"})
	}
	log.Printf("Found task with ticket_no: %s", ticketno)

	err = db.DB.QueryRow(`
		SELECT progress_text
		FROM progress WHERE id = ?
	`, progressid).Scan(&existingProgressText)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Progress not found"})
	}

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
			var progressData map[string]interface{}

			if err := json.Unmarshal([]byte(existingProgressText), &progressData); err == nil {
				if fileList, ok := progressData["files"].([]interface{}); ok {
					for _, file := range fileList {
						if fileMap, ok := file.(map[string]interface{}); ok {
							if url, ok := fileMap["url"].(string); ok {
								existingURLs = append(existingURLs, url)
							}
						}
					}
				}
			}

			// ถ้า URLs ตรงกันทั้งหมด และไม่มีไฟล์ใหม่
			if len(existingURLs) == len(keepImageURLs) {
				allMatch := true
				for _, keepURL := range keepImageURLs {
					found := false
					for _, existingURL := range existingURLs {
						if keepURL == existingURL {
							found = true
							break
						}
					}
					if !found {
						allMatch = false
						break
					}
				}
				if allMatch {
					// URLs ตรงกันทั้งหมด ใช้ text เดิมถ้าไม่ได้ส่งมาใหม่
					if req.Text == "" {
						if text, ok := progressData["text"].(string); ok {
							req.Text = text
						} else {
							req.Text = existingProgressText
						}
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
		var progressData map[string]interface{}
		if err := json.Unmarshal([]byte(existingProgressText), &progressData); err == nil {
			if fileList, ok := progressData["files"].([]interface{}); ok {
				for _, file := range fileList {
					if fileMap, ok := file.(map[string]interface{}); ok {
						if url, ok := fileMap["url"].(string); ok {
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
								keepImage := false
								for _, keepURL := range keepImageURLs {
									if url == keepURL {
										keepImage = true
										break
									}
								}
								// ถ้าไม่ต้องการเก็บ ให้ลบออกจาก MinIO
								if !keepImage {
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
		// ดึง text จาก progress_text เดิม
		var progressData map[string]interface{}
		if err := json.Unmarshal([]byte(existingProgressText), &progressData); err == nil {
			if text, ok := progressData["text"].(string); ok {
				req.Text = text
			}
		} else {
			req.Text = existingProgressText
		}
	}

	// เตรียม progress text ใหม่
	var finalProgressText string
	if len(uploadedFiles) > 0 {
		// มีไฟล์ใหม่ สร้าง JSON ที่มี text และ files
		progressData := fiber.Map{
			"text":  req.Text,
			"files": uploadedFiles,
		}
		progressDataBytes, err := json.Marshal(progressData)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to marshal progress data"})
		}
		finalProgressText = string(progressDataBytes)
	} else {
		// ไม่มีไฟล์ ใช้แค่ text อย่างเดียว
		finalProgressText = req.Text
	}

	// อัปเดต progress
	_, err = db.DB.Exec(`UPDATE progress SET progress_text = ? WHERE id = ?`, finalProgressText, progressid)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update progress"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Progress updated successfully",
	})

}
