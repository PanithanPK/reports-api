package middleware

import (
	"encoding/json"
	"log"
	"reports-api/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// HeaderMiddleware adds common headers to responses
func HeaderMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/json")
		c.Set("Applications", "API")
		c.Set("Version", "1.0")
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Set("Content-Security-Policy", "default-src 'self'")
		return c.Next()
	}
}

// LoggingMiddleware logs incoming requests and responses with timing
func LoggingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Log incoming request
		log.Printf("[REQUEST] %s %s | IP: %s | User-Agent: %s",
			c.Method(),
			c.OriginalURL(),
			c.IP(),
			c.Get("User-Agent"),
		)

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Get response status and body
		status := c.Response().StatusCode()
		responseBody := string(c.Response().Body())

		// Truncate response body if too long for logging
		if len(responseBody) > 10 {
			responseBody = responseBody[:10] + "..."
		}

		// Log response
		log.Printf("[RESPONSE] %s %s | Status: %d | Duration: %v | Response: %s",
			c.Method(),
			c.OriginalURL(),
			status,
			duration,
			responseBody,
		)

		// Log with color coding based on status
		statusEmoji := "✅"
		if status >= 400 && status < 500 {
			statusEmoji = "⚠️"
		} else if status >= 500 {
			statusEmoji = "❌"
		}

		log.Printf("%s [SUMMARY] %s %s | %d | %v",
			statusEmoji,
			c.Method(),
			c.OriginalURL(),
			status,
			duration,
		)

		return err
	}
}

// StandardResponse represents the standard API response format

// CompressionMiddleware enables gzip compression for responses
func CompressionMiddleware() fiber.Handler {
	return compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // Balance between speed and compression ratio
	})
}

// ResponseStandardizationMiddleware converts responses to standard format
func ResponseStandardizationMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Generate request ID
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}
		c.Set("X-Request-ID", requestID)

		// Process the request
		err := c.Next()

		// Only standardize JSON responses
		contentType := c.Get("Content-Type")
		if contentType != "application/json" {
			return err
		}

		// Get current response
		status := c.Response().StatusCode()
		body := c.Response().Body()

		// Skip if already in standard format or empty
		if len(body) == 0 {
			return err
		}

		// Try to parse existing response
		var existingData any
		if json.Unmarshal(body, &existingData) != nil {
			return err
		}

		// Check if already standardized
		if standardResp, ok := existingData.(map[string]any); ok {
			if _, hasSuccess := standardResp["success"]; hasSuccess {
				return err // Already standardized
			}
		}

		// Create standard response
		standardResponse := models.StandardResponse{
			Success:   status >= 200 && status < 300,
			Timestamp: time.Now().Format(time.RFC3339),
			RequestID: requestID,
		}

		if status >= 400 {
			// Error response
			standardResponse.Message = getStatusMessage(status)
			standardResponse.Error = existingData
		} else {
			// Success response
			standardResponse.Message = getStatusMessage(status)
			standardResponse.Data = existingData
		}

		// Set the standardized response
		return c.JSON(standardResponse)
	}
}

// generateRequestID creates a simple request ID
func generateRequestID() string {
	return time.Now().Format("20060102150405") + "-" + time.Now().Format("000")
}

// getStatusMessage returns appropriate message for status code
func getStatusMessage(status int) string {
	switch {
	case status >= 200 && status < 300:
		return "Success"
	case status == 400:
		return "Bad Request"
	case status == 401:
		return "Unauthorized"
	case status == 403:
		return "Forbidden"
	case status == 404:
		return "Not Found"
	case status == 429:
		return "Too Many Requests"
	case status >= 500:
		return "Internal Server Error"
	default:
		return "Unknown Status"
	}
}

// RateLimiter creates a rate limiting middleware
func RateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(429).JSON(fiber.Map{
				"error": "Too many requests",
			})
		},
	})
}
