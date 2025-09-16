package handlers

import (
	"log"
	"reports-api/db"
	"strconv"

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

	// ตรวจสอบว่า task_id มีอยู่ในตาราง tasks หรือไม่
	var exists int
	err = db.DB.QueryRow("SELECT COUNT(*) FROM tasks WHERE id = ? AND deleted_at IS NULL", taskID).Scan(&exists)
	if err != nil {
		log.Printf("Error checking task existence: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}

	if exists == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Task not found"})
	}

	// รับข้อมูล progress_text จาก form
	progressText := c.FormValue("text")
	if progressText == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Progress text is required"})
	}

	// บันทึกข้อมูลลงในตาราง progress
	result, err := db.DB.Exec(
		"INSERT INTO progress (task_id, progress_text, created_at, updated_at) VALUES (?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)",
		taskID, progressText,
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

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Progress entry created successfully",
		"data": fiber.Map{
			"id":      progressID,
			"task_id": taskID,
			"text":    progressText,
		},
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
	var exists int
	err = db.DB.QueryRow("SELECT COUNT(*) FROM tasks WHERE id = ? AND deleted_at IS NULL", taskID).Scan(&exists)
	if err != nil {
		log.Printf("Error checking task existence: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}

	if exists == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Task not found"})
	}

	// ดึงข้อมูล progress entries สำหรับ task นี้
	rows, err := db.DB.Query(
		"SELECT id, task_id, progress_text, created_at, updated_at FROM progress WHERE task_id = ? AND deleted_at IS NULL ORDER BY created_at DESC",
		taskID,
	)
	if err != nil {
		log.Printf("Error querying progress entries: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get progress entries"})
	}
	defer rows.Close()

	var progressEntries []fiber.Map
	for rows.Next() {
		var id, taskID int
		var progressText, createdAt, updatedAt string

		err := rows.Scan(&id, &taskID, &progressText, &createdAt, &updatedAt)
		if err != nil {
			log.Printf("Error scanning progress row: %v", err)
			continue
		}

		progressEntries = append(progressEntries, fiber.Map{
			"id":            id,
			"task_id":       taskID,
			"progress_text": progressText,
			"created_at":    createdAt,
			"updated_at":    updatedAt,
		})
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
