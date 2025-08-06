package tests

import (
	"reports-api/models"
	"reports-api/utils"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestGetPaginationParams(t *testing.T) {
	app := fiber.New()

	t.Run("Default values", func(t *testing.T) {
		app.Get("/test", func(c *fiber.Ctx) error {
			params := utils.GetPaginationParams(c)
			assert.Equal(t, 1, params.Page)
			assert.Equal(t, 10, params.Limit)
			return c.JSON(params)
		})
	})

	t.Run("Custom values", func(t *testing.T) {
		app.Get("/test2", func(c *fiber.Ctx) error {
			params := utils.GetPaginationParams(c)
			assert.Equal(t, 2, params.Page)
			assert.Equal(t, 20, params.Limit)
			return c.JSON(params)
		})
	})

	t.Run("Invalid values", func(t *testing.T) {
		app.Get("/test3", func(c *fiber.Ctx) error {
			params := utils.GetPaginationParams(c)
			assert.Equal(t, 1, params.Page) // Should default to 1
			assert.Equal(t, 10, params.Limit) // Should default to 10
			return c.JSON(params)
		})
	})
}

func TestCalculateOffset(t *testing.T) {
	tests := []struct {
		page     int
		limit    int
		expected int
	}{
		{1, 10, 0},
		{2, 10, 10},
		{3, 20, 40},
		{1, 5, 0},
	}

	for _, test := range tests {
		result := utils.CalculateOffset(test.page, test.limit)
		assert.Equal(t, test.expected, result)
	}
}

func TestCalculateTotalPages(t *testing.T) {
	tests := []struct {
		total    int
		limit    int
		expected int
	}{
		{0, 10, 0},
		{10, 10, 1},
		{15, 10, 2},
		{20, 10, 2},
		{21, 10, 3},
	}

	for _, test := range tests {
		result := utils.CalculateTotalPages(test.total, test.limit)
		assert.Equal(t, test.expected, result)
	}
}

func TestPaginationModels(t *testing.T) {
	t.Run("PaginationRequest", func(t *testing.T) {
		req := models.PaginationRequest{
			Page:  1,
			Limit: 10,
		}
		assert.Equal(t, 1, req.Page)
		assert.Equal(t, 10, req.Limit)
	})

	t.Run("PaginationResponse", func(t *testing.T) {
		resp := models.PaginationResponse{
			Page:       1,
			Limit:      10,
			Total:      100,
			TotalPages: 10,
		}
		assert.Equal(t, 1, resp.Page)
		assert.Equal(t, 10, resp.Limit)
		assert.Equal(t, 100, resp.Total)
		assert.Equal(t, 10, resp.TotalPages)
	})

	t.Run("PaginatedResponse", func(t *testing.T) {
		data := []string{"item1", "item2"}
		resp := models.PaginatedResponse{
			Success: true,
			Data:    data,
			Pagination: models.PaginationResponse{
				Page:       1,
				Limit:      10,
				Total:      2,
				TotalPages: 1,
			},
		}
		assert.True(t, resp.Success)
		assert.Equal(t, data, resp.Data)
		assert.Equal(t, 1, resp.Pagination.Page)
	})
}