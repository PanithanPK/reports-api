package middleware

import (
	"encoding/json"
	"log"
	"reports-api/models"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/google/uuid"
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

// CompressionMiddleware enables gzip compression for responses
func CompressionMiddleware() fiber.Handler {
	return compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // Balance between speed and compression ratio
	})
}

// ResponseStandardizationMiddleware converts responses to standard format with enhanced error handling
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

		// Handle middleware errors
		if err != nil {
			return handleMiddlewareError(c, err, requestID)
		}

		// Only standardize JSON responses
		contentType := c.Get("Content-Type")
		if contentType != "application/json" {
			return nil
		}

		// Get current response
		status := c.Response().StatusCode()
		body := c.Response().Body()

		// Skip if empty response
		if len(body) == 0 {
			return nil
		}

		// Try to parse existing response
		var existingData any
		if err := json.Unmarshal(body, &existingData); err != nil {
			// Handle malformed JSON
			return c.Status(500).JSON(models.StandardResponse{
				Success:   false,
				Message:   "Internal Server Error",
				Error:     "Invalid JSON response",
				Timestamp: time.Now().Format(time.RFC3339),
				RequestID: requestID,
			})
		}

		// Check if already standardized
		if standardResp, ok := existingData.(map[string]any); ok {
			if _, hasSuccess := standardResp["success"]; hasSuccess {
				return nil // Already standardized
			}
		}

		// Create standard response
		standardResponse := models.StandardResponse{
			Success:   status >= 200 && status < 300,
			Timestamp: time.Now().Format(time.RFC3339),
			RequestID: requestID,
		}

		if status >= 400 {
			// Error response with detailed error handling
			standardResponse.Message = getStatusMessage(status)
			standardResponse.Error = processErrorData(existingData)
		} else {
			// Success response
			standardResponse.Message = getStatusMessage(status)
			standardResponse.Data = existingData
		}

		// Set the standardized response
		return c.JSON(standardResponse)
	}
}

// handleMiddlewareError handles errors that occur during middleware processing
func handleMiddlewareError(c *fiber.Ctx, err error, requestID string) error {
	log.Printf("Middleware error: %v", err)

	// Determine appropriate status code based on error type
	status := 500
	message := "Internal Server Error"

	if fiberErr, ok := err.(*fiber.Error); ok {
		status = fiberErr.Code
		message = fiberErr.Message
	}

	return c.Status(status).JSON(models.StandardResponse{
		Success:   false,
		Message:   message,
		Error:     err.Error(),
		Timestamp: time.Now().Format(time.RFC3339),
		RequestID: requestID,
	})
}

// processErrorData processes error data to provide more structured error information
func processErrorData(data any) any {
	if errorMap, ok := data.(map[string]any); ok {
		// If it's already a structured error, return as is
		if _, hasError := errorMap["error"]; hasError {
			return data
		}
	}

	// For simple error messages, wrap in a structured format
	if errorStr, ok := data.(string); ok {
		return map[string]any{
			"error": errorStr,
			"type":  "validation_error",
		}
	}

	return data
}

// generateRequestID creates a unique request ID using UUID
func generateRequestID() string {
	return uuid.New().String()
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

// RateLimiterConfig holds configuration for different endpoints
type RateLimiterConfig struct {
	Max        int
	Expiration time.Duration
}

// getEndpointLimits returns rate limit configuration for different endpoints
func getEndpointLimits(path string) RateLimiterConfig {
	switch {
	case strings.Contains(path, "/api/v1/problem/create"):
		return RateLimiterConfig{Max: 10, Expiration: 1 * time.Minute}
	case strings.Contains(path, "/api/v1/problem/update"):
		return RateLimiterConfig{Max: 20, Expiration: 1 * time.Minute}
	case strings.Contains(path, "/api/v1/problem/delete"):
		return RateLimiterConfig{Max: 5, Expiration: 1 * time.Minute}
	case strings.Contains(path, "/api/v1/problem/list"):
		return RateLimiterConfig{Max: 200, Expiration: 1 * time.Minute}
	default:
		return RateLimiterConfig{Max: 100, Expiration: 1 * time.Minute}
	}
}

// RateLimiter creates an endpoint-specific rate limiting middleware
func RateLimiter() fiber.Handler {
	return func(c *fiber.Ctx) error {
		config := getEndpointLimits(c.Path())

		limiterMiddleware := limiter.New(limiter.Config{
			Max:        config.Max,
			Expiration: config.Expiration,
			KeyGenerator: func(c *fiber.Ctx) string {
				return c.IP() + ":" + c.Path()
			},
			LimitReached: func(c *fiber.Ctx) error {
				return c.Status(429).JSON(fiber.Map{
					"error":  "Too many requests for this endpoint",
					"limit":  config.Max,
					"window": config.Expiration.String(),
				})
			},
		})

		return limiterMiddleware(c)
	}
}
