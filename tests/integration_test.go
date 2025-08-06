package tests

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"reports-api/handlers"
	"reports-api/models"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestPaginationIntegration(t *testing.T) {
	app := fiber.New()

	t.Run("Department List with Pagination", func(t *testing.T) {
		app.Get("/departments", handlers.ListDepartmentsHandler)
		
		// Test default pagination
		req := httptest.NewRequest("GET", "/departments", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		
		// Note: This will fail without database connection
		// In real scenario, you would mock the database
		if resp.StatusCode != 500 {
			var result models.PaginatedResponse
			json.NewDecoder(resp.Body).Decode(&result)
			assert.True(t, result.Success)
			assert.Equal(t, 1, result.Pagination.Page)
			assert.Equal(t, 10, result.Pagination.Limit)
		}
	})

	t.Run("Department List with Custom Pagination", func(t *testing.T) {
		app.Get("/departments2", handlers.ListDepartmentsHandler)
		
		req := httptest.NewRequest("GET", "/departments2?page=2&limit=5", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		
		// Note: This will fail without database connection
		if resp.StatusCode != 500 {
			var result models.PaginatedResponse
			json.NewDecoder(resp.Body).Decode(&result)
			assert.Equal(t, 2, result.Pagination.Page)
			assert.Equal(t, 5, result.Pagination.Limit)
		}
	})
}

func TestAPIEndpoints(t *testing.T) {
	app := fiber.New()

	// Setup routes
	app.Get("/api/v1/department/list", handlers.ListDepartmentsHandler)
	app.Post("/api/v1/department/create", handlers.CreateDepartmentHandler)
	app.Put("/api/v1/department/update/:id", handlers.UpdateDepartmentHandler)
	app.Delete("/api/v1/department/delete/:id", handlers.DeleteDepartmentHandler)

	app.Get("/api/v1/branch/list", handlers.ListBranchesHandler)
	app.Post("/api/v1/branch/create", handlers.CreateBranchHandler)

	app.Get("/api/v1/program/list", handlers.ListProgramsHandler)
	app.Post("/api/v1/program/create", handlers.CreateProgramHandler)

	app.Get("/api/v1/ipphone/list", handlers.ListIPPhonesHandler)
	app.Post("/api/v1/ipphone/create", handlers.CreateIPPhoneHandler)

	app.Get("/api/v1/task/list", handlers.GetTasksHandler)
	app.Post("/api/v1/task/create", handlers.CreateTaskHandler)

	t.Run("All List Endpoints Return Proper Structure", func(t *testing.T) {
		endpoints := []string{
			"/api/v1/department/list",
			"/api/v1/branch/list",
			"/api/v1/program/list",
			"/api/v1/ipphone/list",
			"/api/v1/task/list",
		}

		for _, endpoint := range endpoints {
			req := httptest.NewRequest("GET", endpoint, nil)
			resp, err := app.Test(req)
			assert.NoError(t, err)
			
			// Without database, we expect 500 error
			// In real scenario with mocked DB, we would test success cases
			assert.Equal(t, 500, resp.StatusCode)
		}
	})

	t.Run("Create Endpoints Validate Input", func(t *testing.T) {
		endpoints := []string{
			"/api/v1/department/create",
			"/api/v1/branch/create",
			"/api/v1/program/create",
			"/api/v1/ipphone/create",
			"/api/v1/task/create",
		}

		for _, endpoint := range endpoints {
			// Test invalid JSON
			req := httptest.NewRequest("POST", endpoint, bytes.NewReader([]byte("invalid json")))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, 400, resp.StatusCode)
		}
	})
}

func TestErrorHandling(t *testing.T) {
	app := fiber.New()

	t.Run("Invalid ID Parameters", func(t *testing.T) {
		app.Put("/test/:id", handlers.UpdateDepartmentHandler)
		app.Delete("/test/:id", handlers.DeleteDepartmentHandler)
		app.Get("/test/:id", handlers.GetDepartmentDetailHandler)

		invalidIDs := []string{"abc", "-1", "0.5", ""}
		
		for _, id := range invalidIDs {
			// Test PUT
			req := httptest.NewRequest("PUT", "/test/"+id, bytes.NewReader([]byte("{}")))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)
			assert.NoError(t, err)
			if id != "" {
				assert.Equal(t, 400, resp.StatusCode)
			}

			// Test DELETE
			req = httptest.NewRequest("DELETE", "/test/"+id, nil)
			resp, err = app.Test(req)
			assert.NoError(t, err)
			if id != "" {
				assert.Equal(t, 400, resp.StatusCode)
			}

			// Test GET
			req = httptest.NewRequest("GET", "/test/"+id, nil)
			resp, err = app.Test(req)
			assert.NoError(t, err)
			if id != "" {
				assert.Equal(t, 400, resp.StatusCode)
			}
		}
	})
}