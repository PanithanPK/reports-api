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
