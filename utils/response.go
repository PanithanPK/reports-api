package utils

import (
	"reports-api/constants"
	"reports-api/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

// SuccessResponse returns a standardized success response
func SuccessResponse(c *fiber.Ctx, data interface{}) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": constants.MessageSuccess,
		"data":    data,
	})
}

// ErrorResponse returns a standardized error response
func ErrorResponse(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"success": false,
		"error":   message,
	})
}

// CreatedResponse returns a standardized created response
func CreatedResponse(c *fiber.Ctx, id interface{}) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": constants.MessageCreated,
		"id":      id,
	})
}

// UpdatedResponse returns a standardized updated response
func UpdatedResponse(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": constants.MessageUpdated,
	})
}

// DeletedResponse returns a standardized deleted response
func DeletedResponse(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": constants.MessageDeleted,
	})
}

// PaginatedSuccessResponse returns a paginated success response
func PaginatedSuccessResponse(c *fiber.Ctx, data interface{}, pagination models.PaginationResponse) error {
	return c.JSON(models.PaginatedResponse{
		Success:    true,
		Message:    constants.MessageSuccess,
		Data:       data,
		Pagination: pagination,
		Timestamp:  time.Now().Format(time.RFC3339),
	})
}
