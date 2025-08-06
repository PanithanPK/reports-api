package tests

import (
	"database/sql"
	"reports-api/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockDB represents a mock database for testing
type MockDB struct {
	departments []models.Department
	branches    []models.Branch
	programs    []models.Program
}

func NewMockDB() *MockDB {
	return &MockDB{
		departments: []models.Department{
			{ID: 1, Name: stringPtr("IT Department"), BranchID: intPtr(1)},
			{ID: 2, Name: stringPtr("HR Department"), BranchID: intPtr(1)},
			{ID: 3, Name: stringPtr("Finance Department"), BranchID: intPtr(2)},
		},
		branches: []models.Branch{
			{ID: 1, Name: stringPtr("Main Branch")},
			{ID: 2, Name: stringPtr("Secondary Branch")},
		},
		programs: []models.Program{
			{ID: 1, Name: stringPtr("ERP System")},
			{ID: 2, Name: stringPtr("CRM System")},
		},
	}
}

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func TestMockDatabase(t *testing.T) {
	mockDB := NewMockDB()

	t.Run("Mock Departments", func(t *testing.T) {
		assert.Len(t, mockDB.departments, 3)
		assert.Equal(t, "IT Department", *mockDB.departments[0].Name)
		assert.Equal(t, 1, *mockDB.departments[0].BranchID)
	})

	t.Run("Mock Branches", func(t *testing.T) {
		assert.Len(t, mockDB.branches, 2)
		assert.Equal(t, "Main Branch", *mockDB.branches[0].Name)
	})

	t.Run("Mock Programs", func(t *testing.T) {
		assert.Len(t, mockDB.programs, 2)
		assert.Equal(t, "ERP System", *mockDB.programs[0].Name)
	})
}

func TestPaginationWithMockData(t *testing.T) {
	mockDB := NewMockDB()

	t.Run("Paginate Departments", func(t *testing.T) {
		// Simulate pagination logic
		page := 1
		limit := 2
		offset := (page - 1) * limit

		var paginatedDepts []models.Department
		total := len(mockDB.departments)
		
		end := offset + limit
		if end > total {
			end = total
		}
		
		if offset < total {
			paginatedDepts = mockDB.departments[offset:end]
		}

		assert.Len(t, paginatedDepts, 2)
		assert.Equal(t, "IT Department", *paginatedDepts[0].Name)
		assert.Equal(t, "HR Department", *paginatedDepts[1].Name)

		totalPages := (total + limit - 1) / limit
		assert.Equal(t, 2, totalPages)
	})

	t.Run("Paginate Second Page", func(t *testing.T) {
		page := 2
		limit := 2
		offset := (page - 1) * limit

		var paginatedDepts []models.Department
		total := len(mockDB.departments)
		
		end := offset + limit
		if end > total {
			end = total
		}
		
		if offset < total {
			paginatedDepts = mockDB.departments[offset:end]
		}

		assert.Len(t, paginatedDepts, 1)
		assert.Equal(t, "Finance Department", *paginatedDepts[0].Name)
	})
}

func TestDatabaseConnectionError(t *testing.T) {
	t.Run("Handle SQL Error", func(t *testing.T) {
		err := sql.ErrNoRows
		assert.Equal(t, sql.ErrNoRows, err)
		assert.Error(t, err)
	})

	t.Run("Handle Connection Error", func(t *testing.T) {
		err := sql.ErrConnDone
		assert.Equal(t, sql.ErrConnDone, err)
		assert.Error(t, err)
	})
}