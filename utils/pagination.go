package utils

import (
	"reports-api/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// GetPaginationParams extracts and validates pagination parameters from query string
func GetPaginationParams(c *fiber.Ctx) models.PaginationRequest {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 500 {
		limit = 10
	}

	return models.PaginationRequest{
		Page:  page,
		Limit: limit,
	}
}

// CalculateOffset calculates the SQL OFFSET value
func CalculateOffset(page, limit int) int {
	return (page - 1) * limit
}

// CalculateTotalPages calculates total pages based on total records and limit
func CalculateTotalPages(total, limit int) int {
	if total == 0 {
		return 0
	}
	return (total + limit - 1) / limit
}
