package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"reports-api/config"
	"reports-api/db"
	"reports-api/handlers/common"
	"reports-api/models"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

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
		IFNULL(r.file_paths, '[]') as file_paths
		FROM resolutions r
		WHERE r.id = ?
	`, SolutionID).Scan(&solution, &telegramID, &filePaths)

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
	var AssignedtoID int

	if solutionByStr := c.FormValue("solution"); solutionByStr != "" {
		req.Solution = solutionByStr
	}
	if assigntoStr := c.FormValue("assignto"); assigntoStr != "" {
		req.Assignto = assigntoStr
	}
	if assignedtoIDStr := c.FormValue("assignedto_id"); assignedtoIDStr != "" {
		req.AssignedtoID, _ = strconv.Atoi(assignedtoIDStr)
	}
	var createdAtStr string
	err := db.DB.QueryRow(`
			SELECT ticket_no, IFNULL(assignto_id, 0),IFNULL(assignto, '') AS assignto, IFNULL(reported_by, '') AS reported_by, IFNULL(telegram_id, 0), IFNULL(created_at, '')
			FROM tasks
			WHERE id = ?
		`, id).Scan(&ticketno, &AssignedtoID, &assignto, &reportedby, &telegramID, &createdAtStr)
	if err != nil {
		log.Printf("Failed to retrieve task data for task ID %s: %v", id, err)
		return c.Status(404).JSON(fiber.Map{"error": "Task not found"})
	}

	if req.AssignedtoID == 0 && assignto != "" {
		req.AssignedtoID = AssignedtoID
		req.Assignto = assignto
	}

	if req.Assignto != "" || req.AssignedtoID != 0 {
		_, err := db.DB.Exec(`UPDATE tasks SET assignto_id = ?, assignto = ? WHERE id = ?`, req.AssignedtoID, req.Assignto, id)
		if err != nil {
			log.Printf("Failed to update task assignto: %v", err)
		}
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

	CreatedAt := common.Fixtimefeature(createdAtStr)

	// ลองแยกการ parse ข้อมูล
	form, err := c.MultipartForm()
	if err != nil {
		// ถ้าไม่ใช่ multipart form ให้ใช้ BodyParser ปกติ
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request format"})
		}
	} else {

		// จัดการไฟล์ที่อัปโหลด
		var allFiles []*multipart.FileHeader
		for key, files := range form.File {
			if strings.HasPrefix(key, "image_") || key == "image" {
				allFiles = append(allFiles, files...)
			}
		}

		if len(allFiles) > 0 {
			uploadedFiles, _ = common.HandleFileUploadsResolution(allFiles, ticketno)
		}
	}

	// ตรวจสอบว่ามี solution text หรือไฟล์รูป อย่างน้อยอย่างใดอย่างหนึ่ง

	// เตรียม file paths JSON
	var filePathsJSON interface{}
	if len(uploadedFiles) > 0 {

		filePathsBytes, _ := json.Marshal(uploadedFiles)
		filePathsJSON = string(filePathsBytes)
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
	_, err = db.DB.Exec(`UPDATE tasks SET solution_id = ?, status = 2, resolved_at=CURRENT_TIMESTAMP WHERE id = ?`, resolutionID, id)
	if err != nil {
		log.Printf("Failed to update solution_id in tasks: %q", err)
	}

	// ส่ง solution ไปยัง Telegram ถ้ามี reportID
	// ดึง resolved_at จากฐานข้อมูล resolutions
	resolvedAt, err := common.GetResolvedAtSafely(db.DB, int(resolutionID))
	if err != nil {
		log.Printf("Failed to get resolved_at: %v", err)
		resolvedAt = time.Now() // Fallback to current time
	}
	var Urlenv string
	env := config.AppConfig.Environment
	if env == "dev" {
		Urlenv = "http://helpdesk-dev.nopadol.com/tasks/show/" + id
	} else {
		Urlenv = "http://helpdesk.nopadol.com/tasks/show/" + id
	}

	// เตรียมข้อมูล response
	req.TicketNo = ticketno
	req.CreatedAt = CreatedAt
	req.Url = Urlenv
	req.ResolvedAt = resolvedAt.Add(7 * time.Hour).Format("2006/01/02/ 15:04:05")

	var assignmsgID int
	// ส่ง solution ไปยัง Telegram ถ้ามี reportID
	if reportID > 0 {
		// ดึง MessageID จาก telegram_chat
		err = db.DB.QueryRow(`SELECT IFNULL(report_id, 0), IFNULL(assignto_id, 0) FROM telegram_chat WHERE id = ?`, telegramID).Scan(&req.MessageID, &assignmsgID)
		if err != nil {
			log.Printf("Failed to get message ID: %v", err)
		}

		// ดึงข้อมูลเพิ่มเติมสำหรับ UpdateTelegram
		var phoneNumber int
		var departmentName, branchName, programName string
		var phoneID *int
		var systemID, departmentID int
		var text string
		var telegramUser string

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

		// ดึง telegram_user สำหรับ UpdateAssignedtoMsg
		if req.AssignedtoID > 0 {
			db.DB.QueryRow(`SELECT IFNULL(telegram_username, '') FROM responsibilities WHERE id = ?`, req.AssignedtoID).Scan(&telegramUser)
		}

		// อัปเดตสถานะใน Telegram message
		taskReq := models.TaskRequest{
			PhoneID:        phoneID,
			SystemID:       systemID,
			DepartmentID:   departmentID,
			Text:           text,
			MessageID:      req.MessageID,
			Ticket:         ticketno,
			Assignto:       req.Assignto,
			ReportedBy:     reportedby,
			CreatedAt:      CreatedAt,
			ResolvedAt:     req.ResolvedAt,
			Status:         2,
			Url:            req.Url,
			PhoneNumber:    phoneNumber,
			DepartmentName: departmentName,
			BranchName:     branchName,
			ProgramName:    programName,
			TelegramUser:   telegramUser,
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
			_, err = common.UpdateTelegram(taskReq, photoURLs...)
		} else {
			_, err = common.UpdateTelegram(taskReq)
		}
		if err != nil {
			log.Printf("Failed to update Telegram status: %q", err)
		}

		if assignmsgID > 0 {
			_, err = db.DB.Exec(`UPDATE telegram_chat SET assignto_id = 0 WHERE id = ?`, telegramID)
			if err != nil {
				log.Printf("Failed to update telegram_chat with message ID: %v", err)
			}
			_, err = common.DeleteTelegram(assignmsgID)
			if err != nil {
				log.Printf("Failed to Delete assign message!")
			}
		}

		// เตรียม photo URLs สำหรับ reply message
		var replyPhotoURLs []string
		for _, file := range uploadedFiles {
			if url, ok := file["url"].(string); ok {
				replyPhotoURLs = append(replyPhotoURLs, url)
			}
		}
		req.TelegramUser = telegramUser
		// ส่ง reply message ไปยัง Telegram
		log.Printf("Sending reply to Telegram - MessageID: %d, TelegramUser: %s, PhotoURLs count: %d", req.MessageID, req.TelegramUser, len(replyPhotoURLs))
		replyMessageID, err := common.ReplyToSpecificMessage(req, replyPhotoURLs...)
		if err != nil {
			log.Printf("Failed to send solution to Telegram: %v", err)
		}

		_, err = db.DB.Exec(`UPDATE telegram_chat SET solution_id = ? WHERE id = ?`, replyMessageID, telegramID)
		if err != nil {
			log.Printf("Failed to update telegram_chat with message ID: %v", err)
		}
	}

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
// @Param assignedto formData string false "User to assign the task"
// @Param assignedto_id formData int false "User ID to assign the task"
// @Param image formData file false "Updated resolution image files"
// @Param image_urls formData string false "Existing image URLs to keep (JSON array)"
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
	var Assignto string

	// ดึงข้อมูลเพิ่มเติมสำหรับ UpdateTelegram
	var phoneNumber int
	var departmentName, branchName, programName string
	var phoneID *int
	var systemID, departmentID int
	var text string
	var reportID int
	var ticketno, assignto, reportedby, resolvedat string
	var taskID, assigntoID int
	var createdAtStr string

	req.Solution = c.FormValue("solution")

	if assigntoStr := c.FormValue("assignto"); assigntoStr != "" {
		req.Assignto = assigntoStr
	}
	if assigntoIDStr := c.FormValue("assignedto_id"); assigntoIDStr != "" {
		req.AssignedtoID, _ = strconv.Atoi(assigntoIDStr)
	}
	// ดึงข้อมูล task ทั้งหมดที่จำเป็น
	err := db.DB.QueryRow(`
		SELECT IFNULL(phone_id, 0), IFNULL(system_id, 0), IFNULL(department_id, 0), 
		       IFNULL(text, ''), IFNULL(telegram_id, 0), IFNULL(ticket_no, ''), 
		       IFNULL(assignto, ''), IFNULL(reported_by, ''), IFNULL(assignto_id, 0)
		FROM tasks WHERE id = ?
	`, id).Scan(&phoneID, &systemID, &departmentID, &text, &telegramID, &ticketno, &assignto, &reportedby, &assigntoID)
	if err != nil {
		log.Printf("Failed to get task details: %v", err)
		return c.Status(404).JSON(fiber.Map{"error": "Task not found"})
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

	// ดึง report_id สำหรับ UpdateTelegram
	db.DB.QueryRow(`SELECT report_id FROM telegram_chat WHERE id = ?`, telegramID).Scan(&reportID)

	err = db.DB.QueryRow(`
		SELECT solution_id
		FROM tasks WHERE id = ?
	`, id).Scan(&resolutions)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Solution not found"})
	}

	err = db.DB.QueryRow(`
		SELECT text, telegram_id, IFNULL(file_paths, '[]')
		FROM resolutions WHERE id = ?
	`, resolutions).Scan(&existingResolution.Solution, &telegramID, &existingFilePathsJSON)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Resolution not found"})
	}

	// ดึงข้อมูล task

	err = db.DB.QueryRow(`
		SELECT r.tasks_id, t.ticket_no, IFNULL(t.assignto_id, 0), IFNULL(t.assignto, ''), IFNULL(t.reported_by, ''), t.created_at, t.resolved_at
		FROM resolutions r
		JOIN tasks t ON r.tasks_id = t.id
		WHERE r.id = ?
	`, resolutions).Scan(&taskID, &ticketno, &assigntoID, &assignto, &reportedby, &createdAtStr, &resolvedat)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Resolutions not found"})
	}

	// Parse created_at string to time
	CreatedAt := common.Fixtimefeature(createdAtStr)
	ResolvedAt := common.Fixtimefeature(resolvedat)

	// Parse ข้อมูลจาก request
	var keepImageURLs []string
	form, err := c.MultipartForm()
	if err != nil {
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request format"})
		}
		// รับ image_urls จาก JSON body
		keepImageURLs = req.ImageURLs
	} else {

		// รับ URL รูปเก่าที่ต้องการเก็บไว้
		imageURLsStr := c.FormValue("image_urls")
		if imageURLsStr != "" {
			if err := json.Unmarshal([]byte(imageURLsStr), &keepImageURLs); err != nil {
				log.Printf("Error parsing image_urls: %q", err)
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
			if existingFilePathsJSON != "" && existingFilePathsJSON != "[]" {
				var existingFiles []fiber.Map
				if err := json.Unmarshal([]byte(existingFilePathsJSON), &existingFiles); err == nil {
					for _, file := range existingFiles {
						if url, ok := file["url"].(string); ok {
							existingURLs = append(existingURLs, url)
						}
					}
				}
			}
			// ถ้า URLs ตรงกันทั้งหมด ไม่ต้องทำอะไร
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
					// URLs ตรงกันทั้งหมด ใช้ solution เดิมถ้าไม่ได้ส่งมาใหม่
					if req.Solution == "" {
						req.Solution = existingResolution.Solution
					}
					// อัปเดตเฉพาะ solution
					_, err = db.DB.Exec(`UPDATE resolutions SET text = ? WHERE id = ?`, req.Solution, resolutions)
					if err != nil {
						return c.Status(500).JSON(fiber.Map{"error": "Failed to update resolution"})
					}

					// อัปเดต Telegram
					var solutionMessageID int
					err = db.DB.QueryRow(`SELECT solution_id FROM telegram_chat WHERE id = ?`, telegramID).Scan(&solutionMessageID)
					if err == nil && solutionMessageID > 0 {
						req.TicketNo = ticketno
						if req.Assignto == "" {
							req.Assignto = assignto
						}
						req.CreatedAt = CreatedAt
						req.ResolvedAt = ResolvedAt
						var Urlenv string
						env := config.AppConfig.Environment
						if env == "dev" {
							Urlenv = "http://helpdesk-dev.nopadol.com/tasks/show/" + fmt.Sprintf("%d", taskID)
						} else {
							Urlenv = "http://helpdesk.nopadol.com/tasks/show/" + fmt.Sprintf("%d", taskID)
						}
						req.Url = Urlenv
						req.MessageID = reportID

						messageID, err := common.UpdatereplyToSpecificMessage(solutionMessageID, req, keepImageURLs...)
						if err != nil {
							log.Printf("Failed to update Telegram reply: %v", err)
						}
						db.DB.Exec(`UPDATE telegram_chat SET solution_id = ? WHERE id = ?`, messageID, telegramID)
					}

					return c.JSON(fiber.Map{"success": true, "message": "Resolution updated successfully"})
				}
			}
		}

		// ลบรูปเก่าทั้งหมดถ้าไม่ได้ส่ง image_urls มา หรือลบเฉพาะที่ไม่อยู่ในรายการ
		if existingFilePathsJSON != "" && existingFilePathsJSON != "[]" {
			var existingFiles []fiber.Map
			if err := json.Unmarshal([]byte(existingFilePathsJSON), &existingFiles); err == nil {
				for _, file := range existingFiles {
					if url, ok := file["url"].(string); ok {
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

		// อัปโหลดไฟล์ใหม่ถ้ามี
		if len(allFiles) > 0 {
			uploadedFiles, _ = common.HandleFileUploadsResolution(allFiles, ticketno)
		}

		// รวมรูปเก่าที่เก็บไว้กับรูปใหม่
		for _, keepURL := range keepImageURLs {
			uploadedFiles = append(uploadedFiles, fiber.Map{"url": keepURL})
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
	} else if len(keepImageURLs) > 0 {
		// ใช้เฉพาะรูปเก่าที่เก็บไว้
		var keepFiles []fiber.Map
		for _, url := range keepImageURLs {
			keepFiles = append(keepFiles, fiber.Map{"url": url})
		}
		filePathsBytes, _ := json.Marshal(keepFiles)
		filePathsJSON = string(filePathsBytes)
	} else {
		// ไม่มีไฟล์ใดๆ
		filePathsJSON = nil
	}

	// อัปเดต tasks ถ้ามีการส่ง assignto มา
	if req.Assignto != "" || req.AssignedtoID != 0 {
		_, err = db.DB.Exec(`UPDATE tasks SET assignto_id = ?, assignto = ? WHERE id = ?`, req.AssignedtoID, req.Assignto, id)
		if err != nil {
			log.Printf("Failed to update task assignto: %q", err)
		}
	}

	// อัปเดต resolution
	_, err = db.DB.Exec(`UPDATE resolutions SET text = ?, file_paths = ? WHERE id = ?`, req.Solution, filePathsJSON, resolutions)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update resolution"})
	}

	// ดึง resolved_at
	var Urlenv string
	env := config.AppConfig.Environment
	if env == "dev" {
		Urlenv = "http://helpdesk-dev.nopadol.com/tasks/show/" + id
	} else {
		Urlenv = "http://helpdesk.nopadol.com/tasks/show/" + id
	}
	// เตรียมข้อมูลสำหรับ Telegram
	req.TicketNo = ticketno
	var telegramUser string
	if req.AssignedtoID > 0 {
		err = db.DB.QueryRow(`SELECT IFNULL(name, ''), IFNULL(telegram_username, '') FROM responsibilities WHERE id = ?`, req.AssignedtoID).Scan(&Assignto, &telegramUser)
		if err != nil {
			log.Printf("Failed to get assignto name: %v", err)
			Assignto = req.Assignto // fallback to existing assignto
		}
	} else {
		Assignto = req.Assignto
	}
	req.CreatedAt = CreatedAt
	req.Url = Urlenv
	req.ResolvedAt = ResolvedAt
	req.TelegramUser = telegramUser

	// อัปเดตสถานะใน Telegram message ด้วยข้อมูลที่ครบ
	taskReq := models.TaskRequest{
		PhoneID:        phoneID,
		SystemID:       systemID,
		DepartmentID:   departmentID,
		Text:           text,
		MessageID:      reportID,
		Ticket:         ticketno,
		Assignto:       req.Assignto,
		ReportedBy:     reportedby,
		CreatedAt:      req.CreatedAt,
		ResolvedAt:     req.ResolvedAt,
		Status:         2,
		Url:            req.Url,
		PhoneNumber:    phoneNumber,
		DepartmentName: departmentName,
		BranchName:     branchName,
		ProgramName:    programName,
		TelegramUser:   telegramUser,
	}

	// ดึง file paths จาก task เดิม
	var existingFilePaths string
	db.DB.QueryRow(`SELECT IFNULL(file_paths, '[]') FROM tasks WHERE id = ?`, id).Scan(&existingFilePaths)

	var photoURLs []string
	if existingFilePaths != "" && existingFilePaths != "[]" {
		var existingFiles []fiber.Map
		if err := json.Unmarshal([]byte(existingFilePaths), &existingFiles); err == nil {
			for _, file := range existingFiles {
				if url, ok := file["url"].(string); ok {
					photoURLs = append(photoURLs, url)
				}
			}
		}
	}

	// อัปเดตสถานะใน Telegram
	if len(photoURLs) > 0 {
		_, err = common.UpdateTelegram(taskReq, photoURLs...)
	} else {
		_, err = common.UpdateTelegram(taskReq)
	}
	if err != nil {
		log.Printf("Failed to update Telegram status: %q", err)
	}

	// อัปเดต Telegram reply message ถ้ามี solution_id
	var solutionMessageID int
	err = db.DB.QueryRow(`SELECT solution_id FROM telegram_chat WHERE id = ?`, telegramID).Scan(&solutionMessageID)
	if err == nil && solutionMessageID > 0 {
		// เตรียม photo URLs จากไฟล์ทั้งหมด (เก่าและใหม่)
		var solutionPhotoURLs []string
		if filePathsJSON != nil {
			var allFiles []fiber.Map
			if err := json.Unmarshal([]byte(filePathsJSON.(string)), &allFiles); err == nil {
				for _, file := range allFiles {
					if url, ok := file["url"].(string); ok {
						solutionPhotoURLs = append(solutionPhotoURLs, url)
					}
				}
			}
		}

		// ตั้งค่าข้อมูลให้ครบถ้วนสำหรับ reply
		req.MessageID = reportID
		req.TicketNo = ticketno
		if req.Assignto == "" {
			req.Assignto = Assignto
		}
		req.CreatedAt = CreatedAt
		req.Url = Urlenv

		var messageID int
		// อัปเดต reply message
		log.Printf("Updating Telegram reply - SolutionMessageID: %d, MessageID: %d, TelegramUser: %s, PhotoURLs count: %d", solutionMessageID, req.MessageID, req.TelegramUser, len(solutionPhotoURLs))
		messageID, err = common.UpdatereplyToSpecificMessage(solutionMessageID, req, solutionPhotoURLs...)
		if err != nil {
			log.Printf("Failed to update Telegram reply: %v", err)
		} else {
			log.Printf("Successfully updated Telegram reply message with new ID: %d", messageID)
		}

		_, err = db.DB.Exec(`UPDATE telegram_chat SET solution_id = ? WHERE id = ?`, messageID, telegramID)
		if err != nil {
			log.Printf("Failed to update telegram_chat with message ID: %v", err)
		}
	}

	log.Printf("Updated resolution ID: %d", resolutions)

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
		log.Printf("Invalid id parameter: %s", idStr)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid id"})
	}

	log.Printf("Starting deletion process for task ID: %d", id)

	// Get message_id and file_paths before deleting
	var resolutions int
	var existingResolution models.ResolutionReq
	var telegramID int
	var existingFilePathsJSON string
	var messageID int

	// ดึง solution_id จาก tasks
	err = db.DB.QueryRow(`
		SELECT solution_id
		FROM tasks WHERE id = ?
	`, id).Scan(&resolutions)
	if err != nil {
		log.Printf("Task not found for ID %d: %v", id, err)
		return c.Status(404).JSON(fiber.Map{"error": "task not found"})
	}
	log.Printf("Found solution_id: %d for task ID: %d", resolutions, id)

	// ดึงข้อมูล resolution
	err = db.DB.QueryRow(`
		SELECT text, telegram_id, IFNULL(file_paths, '[]')
		FROM resolutions 
		WHERE id = ?
	`, resolutions).Scan(&existingResolution.Solution, &telegramID, &existingFilePathsJSON)

	if err != nil {
		log.Printf("Resolution not found for ID %d: %v", resolutions, err)
		return c.Status(404).JSON(fiber.Map{"error": "Resolution not found"})
	}
	log.Printf("Found resolution - telegramID: %d, file_paths: %s", telegramID, existingFilePathsJSON)

	// ดึง solution_id จาก telegram_chat (message ID สำหรับลบใน Telegram)
	err = db.DB.QueryRow(`
		SELECT solution_id
		FROM telegram_chat
		WHERE id = ?
	`, telegramID).Scan(&messageID)

	if err != nil {
		log.Printf("Telegram chat not found for ID %d: %v", telegramID, err)
		return c.Status(404).JSON(fiber.Map{"error": "Telegram chat not found"})
	}
	log.Printf("Found solution messageID: %d for telegramID: %d", messageID, telegramID)

	// Delete files from MinIO if they exist
	if existingFilePathsJSON != "" && existingFilePathsJSON != "[]" {
		var filePaths []fiber.Map
		if err := json.Unmarshal([]byte(existingFilePathsJSON), &filePaths); err == nil {
			log.Printf("Deleting %d files from MinIO", len(filePaths))
			for _, fp := range filePaths {
				if url, ok := fp["url"].(string); ok {
					if strings.Contains(url, "prefix=") {
						parts := strings.Split(url, "prefix=")
						if len(parts) > 1 {
							objectName := parts[1]
							log.Printf("Deleting file from MinIO: %s", objectName)
							common.DeleteImage(objectName)
						}
					}
				}
			}
		} else {
			log.Printf("Failed to unmarshal file paths JSON: %v", err)
		}
	}

	// Delete resolution from database
	_, err = db.DB.Exec(`DELETE FROM resolutions WHERE id=?`, resolutions)
	if err != nil {
		log.Printf("Failed to delete resolution ID %d: %v", resolutions, err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete resolutions"})
	}
	log.Printf("Successfully deleted resolution ID: %d", resolutions)

	// อัปเดต solution_id เป็น NULL ใน telegram_chat
	_, err = db.DB.Exec(`UPDATE telegram_chat SET solution_id = NULL WHERE id = ?`, telegramID)
	if err != nil {
		log.Printf("Failed to update telegram_chat solution_id to NULL for ID %d: %v", telegramID, err)
	} else {
		log.Printf("Successfully updated telegram_chat solution_id to NULL for ID: %d", telegramID)
	}

	// อัปเดต solution_id และ status ใน tasks
	_, err = db.DB.Exec(`UPDATE tasks SET solution_id = NULL, status = 0, resolved_at=NULL WHERE id = ?`, id)
	if err != nil {
		log.Printf("Failed to update tasks solution_id to NULL for ID %d: %v", id, err)
	} else {
		log.Printf("Successfully updated task ID %d: solution_id=NULL, status=0, resolved_at=NULL", id)
	}

	// ลบ solution message จาก Telegram ก่อน
	if messageID > 0 {
		log.Printf("Deleting solution message from Telegram, messageID: %d", messageID)
		_, err = common.DeleteTelegram(messageID)
		if err != nil {
			log.Printf("Failed to delete solution message from Telegram (messageID: %d): %v", messageID, err)
		} else {
			log.Printf("Successfully deleted solution message from Telegram (messageID: %d)", messageID)
		}
	}

	// อัปเดตสถานะใน Telegram message กลับเป็น "รอดำเนินการ"
	var reportID int
	err = db.DB.QueryRow(`SELECT report_id FROM telegram_chat WHERE id = ?`, telegramID).Scan(&reportID)
	log.Printf("Debug - telegramID: %d, reportID: %d, query error: %v", telegramID, reportID, err)

	if err != nil {
		log.Printf("Failed to get reportID for telegramID %d: %v", telegramID, err)
		return c.JSON(fiber.Map{"success": true, "warning": "Resolution deleted but failed to update Telegram status"})
	}

	if reportID <= 0 {
		log.Printf("Invalid reportID: %d for telegramID: %d", reportID, telegramID)
		return c.JSON(fiber.Map{"success": true, "warning": "Resolution deleted but invalid reportID"})
	}

	log.Printf("Processing Telegram update for reportID: %d", reportID)

	// ดึงข้อมูล task สำหรับอัปเดต Telegram
	var ticketno, assignto, reportedby string
	var phoneID *int
	var systemID, departmentID int
	var text string
	var assigntoID int
	var taskCreatedAtStr string

	err = db.DB.QueryRow(`
		SELECT ticket_no, IFNULL(assignto, ''), IFNULL(reported_by, ''), phone_id, system_id, department_id, text, created_at, IFNULL(assignto_id, 0)
		FROM tasks WHERE id = ?
	`, id).Scan(&ticketno, &assignto, &reportedby, &phoneID, &systemID, &departmentID, &text, &taskCreatedAtStr, &assigntoID)

	log.Printf("Debug - Task query error: %v", err)
	if err != nil {
		log.Printf("Failed to get task details for ID %d: %v", id, err)
		return c.JSON(fiber.Map{"success": true, "warning": "Resolution deleted but failed to get task details"})
	}

	// แปลง string เป็น time.Time
	CreatedAt := common.Fixtimefeature(taskCreatedAtStr)
	log.Printf("Task details: ticket=%s, assignto=%s, reportedby=%s, phoneID=%v, systemID=%d, departmentID=%d, assigntoID=%d, createdAt=%s",
		ticketno, assignto, reportedby, phoneID, systemID, departmentID, assigntoID, taskCreatedAtStr)

	// ดึงข้อมูลเพิ่มเติม
	var phoneNumber int
	var departmentName, branchName, programName string

	if phoneID != nil && *phoneID > 0 {
		err = db.DB.QueryRow(`
			SELECT p.number, d.name, b.name 
			FROM ip_phones p 
			JOIN departments d ON p.department_id = d.id 
			JOIN branches b ON d.branch_id = b.id 
			WHERE p.id = ?
		`, *phoneID).Scan(&phoneNumber, &departmentName, &branchName)
		if err != nil {
			log.Printf("Phone query error for phoneID %d: %v", *phoneID, err)
		} else {
			log.Printf("Phone details: number=%d, department=%s, branch=%s", phoneNumber, departmentName, branchName)
		}
	} else {
		err = db.DB.QueryRow(`
			SELECT d.name, b.name 
			FROM departments d 
			JOIN branches b ON d.branch_id = b.id 
			WHERE d.id = ?
		`, departmentID).Scan(&departmentName, &branchName)
		if err != nil {
			log.Printf("Department query error for departmentID %d: %v", departmentID, err)
		} else {
			log.Printf("Department details: department=%s, branch=%s", departmentName, branchName)
		}
	}

	if systemID > 0 {
		err = db.DB.QueryRow(`SELECT name FROM systems_program WHERE id = ?`, systemID).Scan(&programName)
		if err != nil {
			log.Printf("System query error for systemID %d: %v", systemID, err)
		} else {
			log.Printf("Program name: %s", programName)
		}
	}

	// ดึง telegram_user สำหรับ UpdateAssignedtoMsg
	var telegramUser string
	if assigntoID > 0 {
		err = db.DB.QueryRow(`SELECT IFNULL(telegram_user, '') FROM responsibilities WHERE id = ?`, assigntoID).Scan(&telegramUser)
		if err != nil {
			log.Printf("Telegram user query error for assigntoID %d: %v", assigntoID, err)
		} else {
			log.Printf("Telegram user: %s", telegramUser)
		}
	}

	// เตรียมข้อมูลสำหรับ UpdateTelegram
	var Urlenv string
	env := config.AppConfig.Environment
	if env == "dev" {
		Urlenv = "http://helpdesk-dev.nopadol.com/tasks/show/" + idStr
	} else {
		Urlenv = "http://helpdesk.nopadol.com/tasks/show/" + idStr
	}

	taskReq := models.TaskRequest{
		PhoneID:        phoneID,
		SystemID:       systemID,
		DepartmentID:   departmentID,
		Text:           text,
		MessageID:      reportID,
		Ticket:         ticketno,
		Assignto:       assignto,
		ReportedBy:     reportedby,
		CreatedAt:      CreatedAt,
		UpdatedAt:      "",
		Status:         0, // เปลี่ยนกลับเป็น "รอดำเนินการ"
		Url:            Urlenv,
		PhoneNumber:    phoneNumber,
		DepartmentName: departmentName,
		BranchName:     branchName,
		ProgramName:    programName,
		TelegramUser:   telegramUser,
	}

	log.Printf("TaskRequest prepared: MessageID=%d, Status=%d, Url=%s", taskReq.MessageID, taskReq.Status, taskReq.Url)

	// ดึง file paths จาก task เดิม (ไม่ใช่จาก resolution)
	var taskFilePathsJSON string
	err = db.DB.QueryRow(`SELECT IFNULL(file_paths, '[]') FROM tasks WHERE id = ?`, id).Scan(&taskFilePathsJSON)
	if err != nil {
		log.Printf("Task file paths query error for ID %d: %v", id, err)
		taskFilePathsJSON = "[]"
	}
	log.Printf("Task file paths: %s", taskFilePathsJSON)

	var photoURLs []string
	if taskFilePathsJSON != "" && taskFilePathsJSON != "[]" {
		var taskFiles []fiber.Map
		if err := json.Unmarshal([]byte(taskFilePathsJSON), &taskFiles); err == nil {
			for _, file := range taskFiles {
				if url, ok := file["url"].(string); ok {
					photoURLs = append(photoURLs, url)
				}
			}
			log.Printf("Extracted %d photo URLs from task files", len(photoURLs))
		} else {
			log.Printf("Failed to unmarshal task file paths JSON: %v", err)
		}
	}

	// อัปเดตสถานะใน Telegram
	log.Printf("Updating Telegram message (reportID: %d) with %d photos", reportID, len(photoURLs))

	var telegramUpdateErr error
	if len(photoURLs) > 0 {
		_, telegramUpdateErr = common.UpdateTelegram(taskReq, photoURLs...)
	} else {
		_, telegramUpdateErr = common.UpdateTelegram(taskReq)
	}

	if telegramUpdateErr != nil {
		log.Printf("Failed to update Telegram status for reportID %d: %v", reportID, telegramUpdateErr)
		return c.JSON(fiber.Map{
			"success":        true,
			"warning":        "Resolution deleted but failed to update Telegram status",
			"telegram_error": telegramUpdateErr.Error(),
		})
	} else {
		log.Printf("Successfully updated Telegram message to pending status for reportID: %d", reportID)
	}

	log.Printf("Completed deletion process for task ID: %d", id)
	return c.JSON(fiber.Map{"success": true})
}
