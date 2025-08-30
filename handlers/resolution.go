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

func CreateResolutionHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	var req models.ResolutionReq
	var uploadedFiles []fiber.Map

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	var ticketno string
	var telegramID int
	err := db.DB.QueryRow(`
				SELECT ticket_no, telegram_id
				FROM tasks
				WHERE id = ?
			`, id).Scan(&ticketno, &telegramID)

	if err != nil {
		log.Printf("Failed to retrieve ticket number: %v", err)
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

	if err != nil {
		log.Printf("Failed to retrieve ticket number: %v", err)
	}

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

		uploadedFiles, _ = handleFileUploadsResolution(allFiles, ticketno)

		// Convert string form values to int for multipart data
		if solutionByStr := c.FormValue("solution"); solutionByStr != "" {
			req.Solution = solutionByStr
		}

	}

	// ตรวจสอบว่ามี solution text หรือไม่
	if req.Solution == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Solution text is required"})
	}

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
	_, err = db.DB.Exec(`UPDATE tasks SET solution_id = ? WHERE id = ?`, resolutionID, id)
	if err != nil {
		log.Printf("Failed to update solution_id in tasks: %v", err)
	}

	// ส่ง solution ไปยัง Telegram ถ้ามี reportID
	if reportID > 0 {
		// เตรียม photo URLs สำหรับ Telegram
		var photoURLs []string
		for _, file := range uploadedFiles {
			if url, ok := file["url"].(string); ok {
				photoURLs = append(photoURLs, url)
			}
		}

		// ส่ง reply message ไปยัง Telegram
		replyMessageID, err := replyToSpecificMessage(reportID, ticketno, req.Solution, photoURLs)
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
	return c.JSON(fiber.Map{"success": true, "id": resolutionID})
}
