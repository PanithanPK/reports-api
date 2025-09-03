package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"reports-api/db"
	"reports-api/models"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func handleFileUploadsResolution(files []*multipart.FileHeader, ticketno string) ([]fiber.Map, []string) {
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
		contentType := "image/jpeg"
		if filepath.Ext(file.Filename) == ".png" {
			contentType = "image/png"
		}
		if err != nil {
			log.Printf("Failed to open %s: %v", file.Filename, err)
			errors = append(errors, fmt.Sprintf("Failed to open %s: %v", file.Filename, err))
			continue
		}

		dateStr := time.Now().Add(7 * time.Hour).Format("01022006")
		filenameSafe := strings.ReplaceAll(file.Filename, " ", "-")
		filename := "solution"
		objectName := fmt.Sprintf("%s-%s-%02d-%s-%s", ticketno, filename, i+1, dateStr, filenameSafe)

		_, err = minioClient.PutObject(
			context.Background(),
			bucketName,
			objectName,
			src,
			file.Size,
			minio.PutObjectOptions{ContentType: contentType},
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

// @Summary Get resolution by task ID
// @Description Get resolution details for a specific task
// @Tags resolutions
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/resolution/{id} [get]
func GetResolutionHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	var solution, filePaths string
	var telegramID int
	var SolutionID int
	var resolvedAt time.Time

	err := db.DB.QueryRow(`
		SELECT IFNULL(t.solution_id, 0) as solution_id
		FROM tasks t
		WHERE t.id = ?
	`, id).Scan(&SolutionID)

	if err != nil {
		log.Printf("Failed to retrieve resolution: %v", err)
		return c.Status(404).JSON(fiber.Map{"error": "Resolution not found SolutionID"})
	}

	err = db.DB.QueryRow(`
		SELECT IFNULL(r.text, '') as text, IFNULL(r.telegram_id, 0) as telegram_id, 
		IFNULL(r.file_paths, '[]') as file_paths, r.resolved_at
		FROM resolutions r
		WHERE r.id = ?
	`, SolutionID).Scan(&solution, &telegramID, &filePaths, &resolvedAt)

	if err != nil {
		log.Printf("Failed to retrieve resolution: %v", err)
		return c.Status(404).JSON(fiber.Map{"error": "Resolution not found"})
	}

	fileMap := make(map[string]string)
	if filePaths != "" && filePaths != "[]" {
		var filePathsArray []fiber.Map
		if err := json.Unmarshal([]byte(filePaths), &filePathsArray); err == nil {
			for i, fp := range filePathsArray {
				if url, ok := fp["url"].(string); ok {
					fileMap[fmt.Sprintf("image_%d", i)] = url
				}
			}
		}
	}

	response := fiber.Map{
		"solution":    solution,
		"telegram_id": telegramID,
		"file_paths":  fileMap,
		"resolved_at": resolvedAt,
	}

	return c.JSON(fiber.Map{"success": true, "data": response})
}

// @Summary Create resolution for task
// @Description Create a new resolution for a specific task with optional file uploads
// @Tags resolutions
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Task ID"
// @Param solution formData string false "Resolution text"
// @Param image formData file false "Resolution image files"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/resolution/create/{id} [post]
func CreateResolutionHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	var req models.ResolutionReq
	var uploadedFiles []fiber.Map

	// ดึงข้อมูล task ก่อน
	var ticketno string
	var assignto string
	var reportedby string
	var telegramID int
	var createdAt time.Time
	err := db.DB.QueryRow(`
			SELECT ticket_no, IFNULL(assignto, '') AS assignto, IFNULL(reported_by, '') AS reported_by, telegram_id, created_at
			FROM tasks
			WHERE id = ?
		`, id).Scan(&ticketno, &assignto, &reportedby, &telegramID, &createdAt)

	if err != nil {
		log.Printf("Failed to retrieve task data: %v", err)
		return c.Status(404).JSON(fiber.Map{"error": "Task not found"})
	}

	var reportID int
	err = db.DB.QueryRow(`
			SELECT report_id
			FROM telegram_chat
			WHERE id = ?
		`, telegramID).Scan(&reportID)

	if err != nil {
		log.Printf("Failed to retrieve report ID: %v", err)
	}

	// ลองแยกการ parse ข้อมูล
	form, err := c.MultipartForm()
	if err != nil {
		// ถ้าไม่ใช่ multipart form ให้ใช้ BodyParser ปกติ
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request format"})
		}
	} else {
		// ถ้าเป็น multipart form ให้ดึงข้อมูลจาก form values
		req.Solution = c.FormValue("solution")

		// จัดการไฟล์ที่อัปโหลด
		var allFiles []*multipart.FileHeader
		for key, files := range form.File {
			if strings.HasPrefix(key, "image_") || key == "image" {
				allFiles = append(allFiles, files...)
			}
		}

		if len(allFiles) > 0 {
			uploadedFiles, _ = handleFileUploadsResolution(allFiles, ticketno)
		}
	}

	// ตรวจสอบว่ามี solution text หรือไฟล์รูป อย่างน้อยอย่างใดอย่างหนึ่ง

	// เตรียม file paths JSON
	var filePathsJSON interface{}
	if len(uploadedFiles) > 0 {
		log.Printf("Uploaded %d files", len(uploadedFiles))
		filePathsBytes, _ := json.Marshal(uploadedFiles)
		filePathsJSON = string(filePathsBytes)
		log.Printf("Saving file_paths: %s", filePathsJSON)
	} else {
		filePathsJSON = nil
	}

	// บันทึก resolution ลงฐานข้อมูล
	res, err := db.DB.Exec(`INSERT INTO resolutions (tasks_id, text, telegram_id, file_paths) VALUES (?, ?, ?, ?)`, id, req.Solution, telegramID, filePathsJSON)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to insert resolution"})
	}

	resolutionID, _ := res.LastInsertId()

	// อัพเดต solution_id ใน tasks
	_, err = db.DB.Exec(`UPDATE tasks SET solution_id = ?, status = 1 WHERE id = ?`, resolutionID, id)
	if err != nil {
		log.Printf("Failed to update solution_id in tasks: %v", err)
	}

	// ส่ง solution ไปยัง Telegram ถ้ามี reportID
	// ดึง resolved_at จากฐานข้อมูล resolutions
	var resolvedAt time.Time
	err = db.DB.QueryRow(`SELECT resolved_at FROM resolutions WHERE id = ?`, resolutionID).Scan(&resolvedAt)
	if err != nil {
		log.Printf("Failed to get resolved_at: %v", err)
	}
	var Urlenv string
	env := os.Getenv("env")
	if env == "dev" {
		Urlenv = "http://helpdesk-dev.nopadol.com/tasks/show/" + id
	} else {
		Urlenv = "http://helpdesk.nopadol.com/tasks/show/" + id
	}

	// เตรียมข้อมูล response
	req.TicketNo = ticketno
	req.Assignto = assignto
	req.CreatedAt = createdAt.Add(7 * time.Hour).Format("02-01-2006 15:04:05")
	req.Url = Urlenv
	req.ResolvedAt = resolvedAt.Add(7 * time.Hour).Format("02-01-2006 15:04:05")

	// ส่ง solution ไปยัง Telegram ถ้ามี reportID
	if reportID > 0 {
		// ดึง MessageID จาก telegram_chat
		err = db.DB.QueryRow(`SELECT report_id FROM telegram_chat WHERE id = ?`, telegramID).Scan(&req.MessageID)
		if err != nil {
			log.Printf("Failed to get message ID: %v", err)
		}

		// ดึงข้อมูลเพิ่มเติมสำหรับ UpdateTelegram
		var phoneNumber int
		var departmentName, branchName, programName string
		var phoneID *int
		var systemID, departmentID int
		var text string

		err = db.DB.QueryRow(`
			SELECT phone_id, system_id, department_id, text
			FROM tasks WHERE id = ?
		`, id).Scan(&phoneID, &systemID, &departmentID, &text)
		if err != nil {
			log.Printf("Failed to get task details: %v", err)
		}
		if phoneID != nil && *phoneID > 0 {
			db.DB.QueryRow(`
				SELECT p.number, d.name, b.name 
				FROM ip_phones p 
				JOIN departments d ON p.department_id = d.id 
				JOIN branches b ON d.branch_id = b.id 
				WHERE p.id = ?
			`, *phoneID).Scan(&phoneNumber, &departmentName, &branchName)
		} else {
			db.DB.QueryRow(`
				SELECT d.name, b.name 
				FROM departments d 
				JOIN branches b ON d.branch_id = b.id 
				WHERE d.id = ?
			`, departmentID).Scan(&departmentName, &branchName)
		}

		if systemID > 0 {
			db.DB.QueryRow(`SELECT name FROM systems_program WHERE id = ?`, systemID).Scan(&programName)
		}

		// อัปเดตสถานะใน Telegram message
		taskReq := models.TaskRequest{
			PhoneID:        phoneID,
			SystemID:       systemID,
			DepartmentID:   departmentID,
			Text:           text,
			MessageID:      req.MessageID,
			Ticket:         ticketno,
			Assignto:       assignto,
			ReportedBy:     reportedby,
			CreatedAt:      req.CreatedAt,
			UpdatedAt:      req.ResolvedAt,
			Status:         1,
			Url:            req.Url,
			PhoneNumber:    phoneNumber,
			DepartmentName: departmentName,
			BranchName:     branchName,
			ProgramName:    programName,
		}

		// ดึง file paths จาก task เดิม
		var existingFilePathsJSON string
		db.DB.QueryRow(`SELECT IFNULL(file_paths, '[]') FROM tasks WHERE id = ?`, id).Scan(&existingFilePathsJSON)

		var photoURLs []string
		if existingFilePathsJSON != "" && existingFilePathsJSON != "[]" {
			var existingFiles []fiber.Map
			if err := json.Unmarshal([]byte(existingFilePathsJSON), &existingFiles); err == nil {
				for _, file := range existingFiles {
					if url, ok := file["url"].(string); ok {
						photoURLs = append(photoURLs, url)
					}
				}
			}
		}

		// อัปเดตสถานะใน Telegram
		if len(photoURLs) > 0 {
			_, err = UpdateTelegram(taskReq, photoURLs...)
		} else {
			_, err = UpdateTelegram(taskReq)
		}
		if err != nil {
			log.Printf("Failed to update Telegram status: %v", err)
		}

		// เตรียม photo URLs สำหรับ reply message
		var replyPhotoURLs []string
		for _, file := range uploadedFiles {
			if url, ok := file["url"].(string); ok {
				replyPhotoURLs = append(replyPhotoURLs, url)
			}
		}

		// ส่ง reply message ไปยัง Telegram
		replyMessageID, err := replyToSpecificMessage(req, replyPhotoURLs...)
		if err != nil {
			log.Printf("Failed to send solution to Telegram: %v", err)
		} else {
			log.Printf("Solution sent to Telegram with reply message ID: %d", replyMessageID)
		}

		_, err = db.DB.Exec(`UPDATE telegram_chat SET solution_id = ? WHERE id = ?`, replyMessageID, telegramID)
		if err != nil {
			log.Printf("Failed to update telegram_chat with message ID: %v", err)
		}
	}

	log.Printf("Inserted new resolution with ID: %d and updated tasks.solution_id", resolutionID)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Resolution created successfully",
	})
}

// @Summary Update resolution
// @Description Update an existing resolution with optional file uploads
// @Tags resolutions
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Task ID"
// @Param solution formData string false "Updated resolution text"
// @Param image formData file false "Updated resolution image files"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/resolution/update/{id} [put]
func UpdateResolutionHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	var req models.ResolutionReq
	var uploadedFiles []fiber.Map

	// ดึงข้อมูล resolution เดิม
	var existingResolution models.ResolutionReq
	var telegramID int
	var existingFilePathsJSON string
	var resolutions int

	err := db.DB.QueryRow(`
		SELECT solution_id
		FROM tasks WHERE id = ?
	`, id).Scan(&resolutions)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "task not found"})
	}

	err = db.DB.QueryRow(`
		SELECT text, telegram_id, IFNULL(file_paths, '[]')
		FROM resolutions WHERE id = ?
	`, resolutions).Scan(&existingResolution.Solution, &telegramID, &existingFilePathsJSON)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Resolution not found"})
	}

	// ดึงข้อมูล task
	var ticketno, assignto, reportedby string
	var createdAt time.Time
	var taskID int
	err = db.DB.QueryRow(`
		SELECT r.tasks_id, t.ticket_no, IFNULL(t.assignto, ''), IFNULL(t.reported_by, ''), t.created_at
		FROM resolutions r
		JOIN tasks t ON r.tasks_id = t.id
		WHERE r.id = ?
	`, resolutions).Scan(&taskID, &ticketno, &assignto, &reportedby, &createdAt)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Task not found"})
	}

	// Parse ข้อมูลจาก request
	form, err := c.MultipartForm()
	if err != nil {
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request format"})
		}
	} else {
		req.Solution = c.FormValue("solution")

		// จัดการไฟล์ใหม่
		var allFiles []*multipart.FileHeader
		for key, files := range form.File {
			if strings.HasPrefix(key, "image_") || key == "image" {
				allFiles = append(allFiles, files...)
			}
		}

		// ถ้ามีไฟล์ใหม่ ให้ลบไฟล์เก่าจาก resolutions ก่อน
		if len(allFiles) > 0 {
			// ลบไฟล์เก่าจาก resolutions file_paths
			if existingFilePathsJSON != "" && existingFilePathsJSON != "[]" {
				var existingFiles []fiber.Map
				if err := json.Unmarshal([]byte(existingFilePathsJSON), &existingFiles); err == nil {
					for _, file := range existingFiles {
						if url, ok := file["url"].(string); ok {
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

			// อัปโหลดไฟล์ใหม่
			uploadedFiles, _ = handleFileUploadsResolution(allFiles, ticketno)
		}
	}

	// ใช้ solution เดิมถ้าไม่ได้ส่งมาใหม่
	if req.Solution == "" {
		req.Solution = existingResolution.Solution
	}

	// เตรียม file paths JSON
	var filePathsJSON interface{}
	if len(uploadedFiles) > 0 {
		filePathsBytes, _ := json.Marshal(uploadedFiles)
		filePathsJSON = string(filePathsBytes)
	} else {
		// ใช้ไฟล์เดิมถ้าไม่มีไฟล์ใหม่
		filePathsJSON = existingFilePathsJSON
		if filePathsJSON == "[]" {
			filePathsJSON = nil
		}
	}

	// อัปเดต resolution
	_, err = db.DB.Exec(`UPDATE resolutions SET text = ?, file_paths = ? WHERE id = ?`, req.Solution, filePathsJSON, resolutions)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update resolution"})
	}

	// ดึง resolved_at
	var resolvedAt time.Time
	err = db.DB.QueryRow(`SELECT resolved_at FROM resolutions WHERE id = ?`, resolutions).Scan(&resolvedAt)
	if err != nil {
		log.Printf("Failed to get resolved_at: %v", err)
	}

	err = db.DB.QueryRow(`SELECT resolved_at FROM tasks WHERE id = ?`, id).Scan(&resolvedAt)
	if err != nil {
		log.Printf("Failed to get resolved_at: %v", err)
	}

	// เตรียม URL
	var Urlenv string
	env := os.Getenv("env")
	if env == "dev" {
		Urlenv = "http://helpdesk-dev.nopadol.com/tasks/show/" + fmt.Sprintf("%d", taskID)
	} else {
		Urlenv = "http://helpdesk.nopadol.com/tasks/show/" + fmt.Sprintf("%d", taskID)
	}

	// เตรียมข้อมูลสำหรับ Telegram
	req.TicketNo = ticketno
	req.Assignto = assignto
	req.CreatedAt = createdAt.Add(7 * time.Hour).Format("02-01-2006 15:04:05")
	req.Url = Urlenv
	req.ResolvedAt = resolvedAt.Add(7 * time.Hour).Format("02-01-2006 15:04:05")

	// อัปเดต Telegram reply message ถ้ามี solution_id
	var solutionMessageID int
	err = db.DB.QueryRow(`SELECT solution_id FROM telegram_chat WHERE id = ?`, telegramID).Scan(&solutionMessageID)
	if err == nil && solutionMessageID > 0 {
		// เตรียม photo URLs
		var photoURLs []string
		if len(uploadedFiles) > 0 {
			// ใช้ไฟล์ใหม่
			for _, file := range uploadedFiles {
				if url, ok := file["url"].(string); ok {
					photoURLs = append(photoURLs, url)
				}
			}
		} else if filePathsJSON != nil {
			// ใช้ไฟล์เดิม
			var existingFiles []fiber.Map
			if err := json.Unmarshal([]byte(filePathsJSON.(string)), &existingFiles); err == nil {
				for _, file := range existingFiles {
					if url, ok := file["url"].(string); ok {
						photoURLs = append(photoURLs, url)
					}
				}
			}
		}
		var messageID int
		// อัปเดต reply message
		messageID, err = UpdatereplyToSpecificMessage(solutionMessageID, req, photoURLs...)
		if err != nil {
			log.Printf("Failed to update Telegram reply: %v", err)
		}

		_, err = db.DB.Exec(`UPDATE telegram_chat SET solution_id = ? WHERE id = ?`, messageID, telegramID)
		if err != nil {
			log.Printf("Failed to update telegram_chat with message ID: %v", err)
		}

	}

	log.Printf("Updated resolution ID: %s", resolutions)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Resolution updated successfully",
	})
}

// @Summary Delete resolution
// @Description Delete a resolution and reset task status
// @Tags resolutions
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/resolution/delete/{id} [delete]
func DeleteResolutionHandler(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid id"})
	}

	// Get message_id and file_paths before deleting
	var resolutions int
	var existingResolution models.ResolutionReq
	var telegramID int
	var existingFilePathsJSON string
	var messageID int

	err = db.DB.QueryRow(`
		SELECT solution_id
		FROM tasks WHERE id = ?
	`, id).Scan(&resolutions)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "task not found"})
	}

	err = db.DB.QueryRow(`
		SELECT text, telegram_id, IFNULL(file_paths, '[]')
		FROM resolutions 
		WHERE id = ?
	`, resolutions).Scan(&existingResolution.Solution, &telegramID, &existingFilePathsJSON)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Resolution not found"})
	}

	err = db.DB.QueryRow(`
		SELECT solution_id
		FROM telegram_chat
		WHERE id = ?
	`, telegramID).Scan(&messageID)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Resolution not found"})
	}

	// Delete files from MinIO if they exist
	if existingFilePathsJSON != "" && existingFilePathsJSON != "[]" {
		var filePaths []fiber.Map
		if err := json.Unmarshal([]byte(existingFilePathsJSON), &filePaths); err == nil {
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

	_, err = db.DB.Exec(`DELETE FROM resolutions WHERE id=?`, resolutions)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete resolutions"})
	}

	// อัปเดต solution_id เป็น NULL ใน telegram_chat
	_, err = db.DB.Exec(`UPDATE telegram_chat SET solution_id = NULL WHERE id = ?`, telegramID)
	if err != nil {
		log.Printf("Failed to update telegram_chat solution_id to NULL: %v", err)
	}

	// อัปเดต solution_id และ status ใน tasks
	_, err = db.DB.Exec(`UPDATE tasks SET solution_id = NULL, status = 0, resolved_at=NULL WHERE id = ?`, id)
	if err != nil {
		log.Printf("Failed to update tasks solution_id to NULL: %v", err)
	}

	// อัปเดตสถานะใน Telegram message กลับเป็น "รอดำเนินการ"
	var reportID int
	err = db.DB.QueryRow(`SELECT report_id FROM telegram_chat WHERE id = ?`, telegramID).Scan(&reportID)
	if err == nil && reportID > 0 {
		// ดึงข้อมูล task สำหรับอัปเดต Telegram
		var ticketno, assignto, reportedby string
		var createdAt time.Time
		var phoneID *int
		var systemID, departmentID int
		var text string

		err = db.DB.QueryRow(`
			SELECT ticket_no, IFNULL(assignto, ''), IFNULL(reported_by, ''), phone_id, system_id, department_id, text, created_at
			FROM tasks WHERE id = ?
		`, id).Scan(&ticketno, &assignto, &reportedby, &phoneID, &systemID, &departmentID, &text, &createdAt)

		if err == nil {
			// ดึงข้อมูลเพิ่มเติม
			var phoneNumber int
			var departmentName, branchName, programName string

			if phoneID != nil && *phoneID > 0 {
				db.DB.QueryRow(`
					SELECT p.number, d.name, b.name 
					FROM ip_phones p 
					JOIN departments d ON p.department_id = d.id 
					JOIN branches b ON d.branch_id = b.id 
					WHERE p.id = ?
				`, *phoneID).Scan(&phoneNumber, &departmentName, &branchName)
			} else {
				db.DB.QueryRow(`
					SELECT d.name, b.name 
					FROM departments d 
					JOIN branches b ON d.branch_id = b.id 
					WHERE d.id = ?
				`, departmentID).Scan(&departmentName, &branchName)
			}

			if systemID > 0 {
				db.DB.QueryRow(`SELECT name FROM systems_program WHERE id = ?`, systemID).Scan(&programName)
			}

			// เตรียม URL
			var Urlenv string
			env := os.Getenv("env")
			if env == "dev" {
				Urlenv = "http://helpdesk-dev.nopadol.com/tasks/show/" + fmt.Sprintf("%d", id)
			} else {
				Urlenv = "http://helpdesk.nopadol.com/tasks/show/" + fmt.Sprintf("%d", id)
			}

			// สร้าง TaskRequest สำหรับอัปเดต Telegram
			taskReq := models.TaskRequest{
				PhoneID:        phoneID,
				SystemID:       systemID,
				DepartmentID:   departmentID,
				Text:           text,
				MessageID:      reportID,
				Ticket:         ticketno,
				Assignto:       assignto,
				ReportedBy:     reportedby,
				CreatedAt:      createdAt.Add(7 * time.Hour).Format("02-01-2006 15:04:05"),
				Status:         0, // เปลี่ยนกลับเป็นรอดำเนินการ
				Url:            Urlenv,
				PhoneNumber:    phoneNumber,
				DepartmentName: departmentName,
				BranchName:     branchName,
				ProgramName:    programName,
			}

			// ดึง file paths จาก task
			var taskFilePathsJSON string
			db.DB.QueryRow(`SELECT IFNULL(file_paths, '[]') FROM tasks WHERE id = ?`, id).Scan(&taskFilePathsJSON)

			var photoURLs []string
			if taskFilePathsJSON != "" && taskFilePathsJSON != "[]" {
				var taskFiles []fiber.Map
				if err := json.Unmarshal([]byte(taskFilePathsJSON), &taskFiles); err == nil {
					for _, file := range taskFiles {
						if url, ok := file["url"].(string); ok {
							photoURLs = append(photoURLs, url)
						}
					}
				}
			}

			// อัปเดตสถานะใน Telegram
			if len(photoURLs) > 0 {
				_, err = UpdateTelegram(taskReq, photoURLs...)
			} else {
				_, err = UpdateTelegram(taskReq)
			}
			if err != nil {
				log.Printf("Failed to update Telegram status: %v", err)
			}
		}
	}

	// Delete Telegram reply message if exists
	if messageID > 0 {
		_, _ = DeleteTelegram(messageID)
	}

	log.Printf("Deleted task ID: %d", id)
	return c.JSON(fiber.Map{"success": true})
}
