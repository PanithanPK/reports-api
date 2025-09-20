package common

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
)

// ImageConfig holds configuration for image processing
type ImageConfig struct {
	MaxWidth    uint
	MaxHeight   uint
	Quality     int   // JPEG quality (1-100)
	MaxFileSize int64 // Maximum file size in bytes (e.g., 10MB for Telegram)
}

// DefaultImageConfig returns default configuration optimized for Telegram
func DefaultImageConfig() ImageConfig {
	return ImageConfig{
		MaxWidth:    1920,             // Max width for images
		MaxHeight:   1080,             // Max height for images
		Quality:     85,               // JPEG quality
		MaxFileSize: 10 * 1024 * 1024, // 10MB - Telegram's limit
	}
}

// ProcessImage resizes and compresses an image file
func ProcessImage(file *multipart.FileHeader, config ImageConfig) (*bytes.Buffer, string, error) {
	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return nil, "", fmt.Errorf("failed to open file: %v", err)
	}
	defer src.Close()

	// Get file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))

	// Decode image based on file type
	var img image.Image
	var originalFormat string

	switch ext {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(src)
		originalFormat = "jpeg"
	case ".png":
		img, err = png.Decode(src)
		originalFormat = "png"
	default:
		return nil, "", fmt.Errorf("unsupported image format: %s", ext)
	}

	if err != nil {
		return nil, "", fmt.Errorf("failed to decode image: %v", err)
	}

	// Get original dimensions
	bounds := img.Bounds()
	originalWidth := uint(bounds.Dx())
	originalHeight := uint(bounds.Dy())

	log.Printf("Original image dimensions: %dx%d", originalWidth, originalHeight)

	// Calculate new dimensions while maintaining aspect ratio
	newWidth, newHeight := calculateNewDimensions(originalWidth, originalHeight, config.MaxWidth, config.MaxHeight)

	// Resize image if needed
	var resizedImg image.Image
	if newWidth != originalWidth || newHeight != originalHeight {
		log.Printf("Resizing image to: %dx%d", newWidth, newHeight)
		resizedImg = resize.Resize(newWidth, newHeight, img, resize.Lanczos3)
	} else {
		resizedImg = img
	}

	// Encode the processed image
	var buf bytes.Buffer
	var contentType string

	// For PNG images, try to compress as JPEG first if it results in smaller file
	if originalFormat == "png" {
		// Try JPEG compression first
		var jpegBuf bytes.Buffer
		err = jpeg.Encode(&jpegBuf, resizedImg, &jpeg.Options{Quality: config.Quality})
		if err == nil && int64(jpegBuf.Len()) < config.MaxFileSize {
			buf = jpegBuf
			contentType = "image/jpeg"
			log.Printf("Converted PNG to JPEG for better compression. Size: %d bytes", buf.Len())
		} else {
			// Fall back to PNG
			err = png.Encode(&buf, resizedImg)
			if err != nil {
				return nil, "", fmt.Errorf("failed to encode PNG: %v", err)
			}
			contentType = "image/png"
			log.Printf("Kept as PNG. Size: %d bytes", buf.Len())
		}
	} else {
		// JPEG format
		err = jpeg.Encode(&buf, resizedImg, &jpeg.Options{Quality: config.Quality})
		if err != nil {
			return nil, "", fmt.Errorf("failed to encode JPEG: %v", err)
		}
		contentType = "image/jpeg"
		log.Printf("Processed as JPEG. Size: %d bytes", buf.Len())
	}

	// Check if file is still too large
	if int64(buf.Len()) > config.MaxFileSize {
		log.Printf("Warning: Processed image size (%d bytes) exceeds limit (%d bytes)", buf.Len(), config.MaxFileSize)

		// Try with lower quality for JPEG
		if contentType == "image/jpeg" && config.Quality > 50 {
			return processWithLowerQuality(resizedImg, config.MaxFileSize)
		}
	}

	log.Printf("Image processing completed. Final size: %d bytes, Content-Type: %s", buf.Len(), contentType)
	return &buf, contentType, nil
}

// processWithLowerQuality tries to compress JPEG with progressively lower quality
func processWithLowerQuality(img image.Image, maxSize int64) (*bytes.Buffer, string, error) {
	qualities := []int{70, 60, 50, 40, 30}

	for _, quality := range qualities {
		var buf bytes.Buffer
		err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality})
		if err != nil {
			continue
		}

		if int64(buf.Len()) <= maxSize {
			log.Printf("Achieved target size with quality %d. Size: %d bytes", quality, buf.Len())
			return &buf, "image/jpeg", nil
		}
	}

	// If still too large, return the lowest quality version
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 30})
	if err != nil {
		return nil, "", fmt.Errorf("failed to encode with lowest quality: %v", err)
	}

	log.Printf("Using lowest quality (30). Final size: %d bytes", buf.Len())
	return &buf, "image/jpeg", nil
}

// calculateNewDimensions calculates new dimensions while maintaining aspect ratio
func calculateNewDimensions(originalWidth, originalHeight, maxWidth, maxHeight uint) (uint, uint) {
	if originalWidth <= maxWidth && originalHeight <= maxHeight {
		return originalWidth, originalHeight
	}

	// Calculate scaling factors
	widthRatio := float64(maxWidth) / float64(originalWidth)
	heightRatio := float64(maxHeight) / float64(originalHeight)

	// Use the smaller ratio to maintain aspect ratio
	ratio := widthRatio
	if heightRatio < widthRatio {
		ratio = heightRatio
	}

	newWidth := uint(float64(originalWidth) * ratio)
	newHeight := uint(float64(originalHeight) * ratio)

	return newWidth, newHeight
}

// ProcessImageFromReader processes an image from an io.Reader (for downloaded images)
func ProcessImageFromReader(reader io.Reader, filename string, config ImageConfig) (*bytes.Buffer, string, error) {
	// Read all data into memory
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read image data: %v", err)
	}

	// Create a reader from the data
	dataReader := bytes.NewReader(data)

	// Get file extension
	ext := strings.ToLower(filepath.Ext(filename))

	// Decode image based on file type
	var img image.Image
	var originalFormat string

	switch ext {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(dataReader)
		originalFormat = "jpeg"
	case ".png":
		img, err = png.Decode(dataReader)
		originalFormat = "png"
	default:
		return nil, "", fmt.Errorf("unsupported image format: %s", ext)
	}

	if err != nil {
		return nil, "", fmt.Errorf("failed to decode image: %v", err)
	}

	// Get original dimensions
	bounds := img.Bounds()
	originalWidth := uint(bounds.Dx())
	originalHeight := uint(bounds.Dy())

	log.Printf("Downloaded image dimensions: %dx%d", originalWidth, originalHeight)

	// Calculate new dimensions while maintaining aspect ratio
	newWidth, newHeight := calculateNewDimensions(originalWidth, originalHeight, config.MaxWidth, config.MaxHeight)

	// Resize image if needed
	var resizedImg image.Image
	if newWidth != originalWidth || newHeight != originalHeight {
		log.Printf("Resizing downloaded image to: %dx%d", newWidth, newHeight)
		resizedImg = resize.Resize(newWidth, newHeight, img, resize.Lanczos3)
	} else {
		resizedImg = img
	}

	// Encode the processed image
	var buf bytes.Buffer
	var contentType string

	// For PNG images, try to compress as JPEG first if it results in smaller file
	if originalFormat == "png" {
		// Try JPEG compression first
		var jpegBuf bytes.Buffer
		err = jpeg.Encode(&jpegBuf, resizedImg, &jpeg.Options{Quality: config.Quality})
		if err == nil && int64(jpegBuf.Len()) < config.MaxFileSize {
			buf = jpegBuf
			contentType = "image/jpeg"
			log.Printf("Converted downloaded PNG to JPEG. Size: %d bytes", buf.Len())
		} else {
			// Fall back to PNG
			err = png.Encode(&buf, resizedImg)
			if err != nil {
				return nil, "", fmt.Errorf("failed to encode PNG: %v", err)
			}
			contentType = "image/png"
			log.Printf("Kept downloaded image as PNG. Size: %d bytes", buf.Len())
		}
	} else {
		// JPEG format
		err = jpeg.Encode(&buf, resizedImg, &jpeg.Options{Quality: config.Quality})
		if err != nil {
			return nil, "", fmt.Errorf("failed to encode JPEG: %v", err)
		}
		contentType = "image/jpeg"
		log.Printf("Processed downloaded image as JPEG. Size: %d bytes", buf.Len())
	}

	// Check if file is still too large
	if int64(buf.Len()) > config.MaxFileSize {
		log.Printf("Warning: Downloaded image size (%d bytes) exceeds limit (%d bytes)", buf.Len(), config.MaxFileSize)

		// Try with lower quality for JPEG
		if contentType == "image/jpeg" && config.Quality > 50 {
			return processWithLowerQuality(resizedImg, config.MaxFileSize)
		}
	}

	log.Printf("Downloaded image processing completed. Final size: %d bytes, Content-Type: %s", buf.Len(), contentType)
	return &buf, contentType, nil
}
