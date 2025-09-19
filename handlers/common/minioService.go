package common

import (
	"context"
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

func HandleFileUploads(files []*multipart.FileHeader, ticketno string) ([]fiber.Map, []string) {
	var uploadedFiles []fiber.Map
	var errors []string

	// MinIO configuration from config
	endpoint := config.AppConfig.EndPoint
	accessKeyID := config.AppConfig.AccessKey
	secretAccessKey := config.AppConfig.SecretAccessKey
	useSSL := false
	bucketName := config.AppConfig.BucketName

	// Initialize MinIO client
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Printf("Failed to create MinIO client: %v", err)
		return uploadedFiles, []string{"Failed to initialize storage client"}
	}

	// loop through files and upload to MinIO
	for i, file := range files {
		// Load file image
		src, err := file.Open()
		// Check file type file image
		contentType := "image/jpeg"
		if filepath.Ext(file.Filename) == ".png" {
			contentType = "image/png"
		}
		if err != nil {
			log.Printf("Failed to open %s: %v", file.Filename, err)
			errors = append(errors, fmt.Sprintf("Failed to open %s: %v", file.Filename, err))
			continue
		}
		// Name Object
		dateStr := time.Now().Add(7 * time.Hour).Format("1504")
		filenameSafe := strings.ReplaceAll(file.Filename, " ", "-")
		filenameSafe = strings.ReplaceAll(filenameSafe, "(", "[")
		filenameSafe = strings.ReplaceAll(filenameSafe, ")", "]")
		objectName := fmt.Sprintf("%s-%02d-%s-%s", ticketno, i+1, dateStr, filenameSafe)

		// Upload to MinIO
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

		// Return URL for get file path
		fileURL := fmt.Sprintf("https://minio.sys9.co/api/v1/buckets/%s/objects/download?preview=true&prefix=%s", bucketName, objectName)
		uploadedFiles = append(uploadedFiles, fiber.Map{
			"url": fileURL,
		})
	}

	return uploadedFiles, errors
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

		dateStr := time.Now().Add(7 * time.Hour).Format("1504")
		filenameSafe := strings.ReplaceAll(file.Filename, " ", "-")
		filenameSafe = strings.ReplaceAll(filenameSafe, "(", "[")
		filenameSafe = strings.ReplaceAll(filenameSafe, ")", "]")
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

// HandleFileUploadsProgress handles file uploads for progress entries
func HandleFileUploadsProgress(files []*multipart.FileHeader, ticketNo string) ([]fiber.Map, []string) {
	var uploadedFiles []fiber.Map
	var errors []string

	// MinIO configuration from config
	endpoint := config.AppConfig.EndPoint
	accessKeyID := config.AppConfig.AccessKey
	secretAccessKey := config.AppConfig.SecretAccessKey
	useSSL := false
	bucketName := config.AppConfig.BucketName

	// Initialize MinIO client
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Printf("Failed to create MinIO client: %v", err)
		return uploadedFiles, []string{"Failed to initialize storage client"}
	}

	// loop through files and upload to MinIO
	for i, file := range files {
		// Load file image
		src, err := file.Open()
		// Check file type file image
		contentType := "image/jpeg"
		if filepath.Ext(file.Filename) == ".png" {
			contentType = "image/png"
		}
		if err != nil {
			log.Printf("Failed to open %s: %v", file.Filename, err)
			errors = append(errors, fmt.Sprintf("Failed to open %s: %v", file.Filename, err))
			continue
		}
		// Name Object for progress
		dateStr := time.Now().Add(7 * time.Hour).Format("1504")
		filenameSafe := strings.ReplaceAll(file.Filename, " ", "-")
		filenameSafe = strings.ReplaceAll(filenameSafe, "(", "[")
		filenameSafe = strings.ReplaceAll(filenameSafe, ")", "]")
		filename := "progress"
		objectName := fmt.Sprintf("%s-%s-%02d-%s-%s", ticketNo, filename, i+1, dateStr, filenameSafe)

		// Upload to MinIO
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

		// Return URL for get file path
		fileURL := fmt.Sprintf("https://minio.sys9.co/api/v1/buckets/%s/objects/download?preview=true&prefix=%s", bucketName, objectName)
		uploadedFiles = append(uploadedFiles, fiber.Map{
			"url": fileURL,
		})
	}

	return uploadedFiles, errors
}

func DeleteImage(objectName string) error {
	// MinIO configuration from config
	endpoint := config.AppConfig.EndPoint
	accessKeyID := config.AppConfig.AccessKey
	secretAccessKey := config.AppConfig.SecretAccessKey
	useSSL := false
	bucketName := config.AppConfig.BucketName

	// Initialize MinIO client
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return err
	}

	// Delete the object
	err = minioClient.RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		log.Printf("Failed to delete %s: %v", objectName, err)
		return err
	}

	log.Printf("Successfully deleted %s", objectName)
	return nil
}
