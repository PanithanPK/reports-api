package common

import (
	"bytes"
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

	// Image processing configuration
	imageConfig := DefaultImageConfig()

	// loop through files and upload to MinIO
	for i, file := range files {
		// Check if file is an image
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			log.Printf("Skipping non-image file: %s", file.Filename)
			errors = append(errors, fmt.Sprintf("Unsupported file type: %s", file.Filename))
			continue
		}

		// Process image (resize and compress)
		processedImage, contentType, err := ProcessImage(file, imageConfig)
		if err != nil {
			log.Printf("Failed to process image %s: %v", file.Filename, err)
			errors = append(errors, fmt.Sprintf("Failed to process image %s: %v", file.Filename, err))
			continue
		}

		// Name Object
		dateStr := time.Now().Add(7 * time.Hour).Format("1504")
		filenameSafe := strings.ReplaceAll(file.Filename, " ", "-")
		filenameSafe = strings.ReplaceAll(filenameSafe, "(", "[")
		filenameSafe = strings.ReplaceAll(filenameSafe, ")", "]")

		// Update filename extension based on processed format
		if contentType == "image/jpeg" && (ext == ".png") {
			// Replace .png with .jpg if converted
			filenameSafe = strings.TrimSuffix(filenameSafe, ext) + ".jpg"
		}

		objectName := fmt.Sprintf("%s-%02d-%s-%s", ticketno, i+1, dateStr, filenameSafe)

		// Upload processed image to MinIO
		_, err = minioClient.PutObject(
			context.Background(),
			bucketName,
			objectName,
			bytes.NewReader(processedImage.Bytes()),
			int64(processedImage.Len()),
			minio.PutObjectOptions{ContentType: contentType},
		)

		if err != nil {
			log.Printf("Failed to upload processed image %s: %v", file.Filename, err)
			errors = append(errors, fmt.Sprintf("Failed to upload processed image %s: %v", file.Filename, err))
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

	// Image processing configuration
	imageConfig := DefaultImageConfig()

	for i, file := range files {
		// Check if file is an image
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			log.Printf("Skipping non-image file: %s", file.Filename)
			errors = append(errors, fmt.Sprintf("Unsupported file type: %s", file.Filename))
			continue
		}

		// Process image (resize and compress)
		processedImage, contentType, err := ProcessImage(file, imageConfig)
		if err != nil {
			log.Printf("Failed to process image %s: %v", file.Filename, err)
			errors = append(errors, fmt.Sprintf("Failed to process image %s: %v", file.Filename, err))
			continue
		}

		dateStr := time.Now().Add(7 * time.Hour).Format("1504")
		filenameSafe := strings.ReplaceAll(file.Filename, " ", "-")
		filenameSafe = strings.ReplaceAll(filenameSafe, "(", "[")
		filenameSafe = strings.ReplaceAll(filenameSafe, ")", "]")

		// Update filename extension based on processed format
		if contentType == "image/jpeg" && (ext == ".png") {
			// Replace .png with .jpg if converted
			filenameSafe = strings.TrimSuffix(filenameSafe, ext) + ".jpg"
		}

		filename := "solution"
		objectName := fmt.Sprintf("%s-%s-%02d-%s-%s", ticketno, filename, i+1, dateStr, filenameSafe)

		_, err = minioClient.PutObject(
			context.Background(),
			bucketName,
			objectName,
			bytes.NewReader(processedImage.Bytes()),
			int64(processedImage.Len()),
			minio.PutObjectOptions{ContentType: contentType},
		)

		if err != nil {
			log.Printf("Failed to upload processed image %s: %v", file.Filename, err)
			errors = append(errors, fmt.Sprintf("Failed to upload processed image %s: %v", file.Filename, err))
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

	// Image processing configuration
	imageConfig := DefaultImageConfig()

	// loop through files and upload to MinIO
	for i, file := range files {
		// Check if file is an image
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			log.Printf("Skipping non-image file: %s", file.Filename)
			errors = append(errors, fmt.Sprintf("Unsupported file type: %s", file.Filename))
			continue
		}

		// Process image (resize and compress)
		processedImage, contentType, err := ProcessImage(file, imageConfig)
		if err != nil {
			log.Printf("Failed to process image %s: %v", file.Filename, err)
			errors = append(errors, fmt.Sprintf("Failed to process image %s: %v", file.Filename, err))
			continue
		}

		// Name Object for progress
		dateStr := time.Now().Add(7 * time.Hour).Format("1504")
		filenameSafe := strings.ReplaceAll(file.Filename, " ", "-")
		filenameSafe = strings.ReplaceAll(filenameSafe, "(", "[")
		filenameSafe = strings.ReplaceAll(filenameSafe, ")", "]")

		// Update filename extension based on processed format
		if contentType == "image/jpeg" && (ext == ".png") {
			// Replace .png with .jpg if converted
			filenameSafe = strings.TrimSuffix(filenameSafe, ext) + ".jpg"
		}

		filename := "progress"
		objectName := fmt.Sprintf("%s-%s-%02d-%s-%s", ticketNo, filename, i+1, dateStr, filenameSafe)

		// Upload processed image to MinIO
		_, err = minioClient.PutObject(
			context.Background(),
			bucketName,
			objectName,
			bytes.NewReader(processedImage.Bytes()),
			int64(processedImage.Len()),
			minio.PutObjectOptions{ContentType: contentType},
		)

		if err != nil {
			log.Printf("Failed to upload processed image %s: %v", file.Filename, err)
			errors = append(errors, fmt.Sprintf("Failed to upload processed image %s: %v", file.Filename, err))
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
