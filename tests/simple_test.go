package tests

import (
	"reports-api/models"
	"reports-api/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaginationUtils(t *testing.T) {
	t.Run("CalculateOffset", func(t *testing.T) {
		assert.Equal(t, 0, utils.CalculateOffset(1, 10))
		assert.Equal(t, 10, utils.CalculateOffset(2, 10))
		assert.Equal(t, 40, utils.CalculateOffset(3, 20))
	})

	t.Run("CalculateTotalPages", func(t *testing.T) {
		assert.Equal(t, 0, utils.CalculateTotalPages(0, 10))
		assert.Equal(t, 1, utils.CalculateTotalPages(10, 10))
		assert.Equal(t, 2, utils.CalculateTotalPages(15, 10))
		assert.Equal(t, 3, utils.CalculateTotalPages(21, 10))
	})
}

func TestPaginationStructs(t *testing.T) {
	t.Run("PaginationRequest", func(t *testing.T) {
		req := models.PaginationRequest{Page: 1, Limit: 10}
		assert.Equal(t, 1, req.Page)
		assert.Equal(t, 10, req.Limit)
	})

	t.Run("PaginationResponse", func(t *testing.T) {
		resp := models.PaginationResponse{
			Page: 1, Limit: 10, Total: 100, TotalPages: 10,
		}
		assert.Equal(t, 1, resp.Page)
		assert.Equal(t, 100, resp.Total)
	})

	t.Run("PaginatedResponse", func(t *testing.T) {
		data := []string{"item1", "item2"}
		resp := models.PaginatedResponse{
			Success: true,
			Data:    data,
			Pagination: models.PaginationResponse{Page: 1, Limit: 10},
		}
		assert.True(t, resp.Success)
		assert.Equal(t, data, resp.Data)
	})
}

func TestBasicPaginationLogic(t *testing.T) {
	t.Run("Basic pagination calculation", func(t *testing.T) {
		// Test basic pagination math
		page := 2
		limit := 10
		offset := (page - 1) * limit
		assert.Equal(t, 10, offset)
		
		total := 25
		totalPages := (total + limit - 1) / limit
		assert.Equal(t, 3, totalPages)
	})
}

func TestDepartmentModel(t *testing.T) {
	name := "IT Department"
	branchID := 1
	dept := models.Department{
		ID:       1,
		Name:     &name,
		BranchID: &branchID,
	}
	assert.Equal(t, 1, dept.ID)
	assert.Equal(t, "IT Department", *dept.Name)
	assert.Equal(t, 1, *dept.BranchID)
}

func TestTaskModel(t *testing.T) {
	req := models.TaskStatusUpdateRequest{
		ID:        1,
		Status:    1,
		UpdatedBy: 123,
	}
	assert.Equal(t, 1, req.ID)
	assert.Equal(t, 1, req.Status)
	assert.Equal(t, 123, req.UpdatedBy)
}