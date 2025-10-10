package utils

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// ParseIDParam parses and validates an ID parameter from the URL
func ParseIDParam(c *fiber.Ctx, paramName string) (int, error) {
	idStr := c.Params(paramName)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}
	if id <= 0 {
		return 0, fiber.ErrBadRequest
	}
	return id, nil
}

// ValidateIDParam validates an ID parameter and returns error response if invalid
func ValidateIDParam(c *fiber.Ctx, paramName string) (int, error) {
	id, err := ParseIDParam(c, paramName)
	if err != nil {
		return 0, c.Status(400).JSON(fiber.Map{"error": "Invalid id"})
	}
	return id, nil
}
