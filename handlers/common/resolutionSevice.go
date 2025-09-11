package common

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"mime/multipart"
	"path/filepath"
	"reports-api/config"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Alternative: Use string scanning then convert
func GetResolvedAtSafely(db *sql.DB, resolutionID int) (time.Time, error) {
	var resolvedAtStr string
	err := db.QueryRow(`SELECT DATE_FORMAT(resolved_at, '%Y-%m-%d %H:%i:%s') FROM resolutions WHERE id = ?`, resolutionID).Scan(&resolvedAtStr)
	if err != nil {
		return time.Time{}, err
	}

	if resolvedAtStr == "" {
		return time.Time{}, nil
	}

	resolvedAt, err := time.Parse("2006-01-02 15:04:05", resolvedAtStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse time %s: %v", resolvedAtStr, err)
	}

	return resolvedAt, nil
}

func HandleFileUploadsResolution(files []*multipart.FileHeader, ticketno string) ([]fiber.Map, []string) {
	var uploadedFiles []fiber.Map
	var errors []string

	// MinIO configuration from config
	endpoint := config.AppConfig.EndPoint
	accessKeyID := config.AppConfig.AccessKey
	secretAccessKey := config.AppConfig.SecretAccessKey
	useSSL := false
	bucketName := config.AppConfig.BucketName

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
			log.Printf("Failed to open file: %v", err)
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
