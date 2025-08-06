package tests

import (
	"bytes"
	"net/http/httptest"
	"reports-api/handlers"
	"reports-api/models"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupTestApp() *fiber.App {
	app := fiber.New()
	return app
}

func TestDepartmentHandlers(t *testing.T) {
	app := setupTestApp()

	t.Run("CreateDepartmentHandler - Invalid JSON", func(t *testing.T) {
		app.Post("/department", handlers.CreateDepartmentHandler)

		req := httptest.NewRequest("POST", "/department", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
	})

	t.Run("UpdateDepartmentHandler - Invalid ID", func(t *testing.T) {
		app.Put("/department/:id", handlers.UpdateDepartmentHandler)

		req := httptest.NewRequest("PUT", "/department/invalid", bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
	})

	t.Run("DeleteDepartmentHandler - Invalid ID", func(t *testing.T) {
		app.Delete("/department/:id", handlers.DeleteDepartmentHandler)

		req := httptest.NewRequest("DELETE", "/department/invalid", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
	})
}

func TestBranchHandlers(t *testing.T) {
	app := setupTestApp()

	t.Run("CreateBranchHandler - Invalid JSON", func(t *testing.T) {
		app.Post("/branch", handlers.CreateBranchHandler)

		req := httptest.NewRequest("POST", "/branch", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
	})

	t.Run("UpdateBranchHandler - Invalid ID", func(t *testing.T) {
		app.Put("/branch/:id", handlers.UpdateBranchHandler)

		req := httptest.NewRequest("PUT", "/branch/invalid", bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
	})
}

func TestProgramHandlers(t *testing.T) {
	app := setupTestApp()

	t.Run("CreateProgramHandler - Invalid JSON", func(t *testing.T) {
		app.Post("/program", handlers.CreateProgramHandler)

		req := httptest.NewRequest("POST", "/program", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
	})

	t.Run("UpdateProgramHandler - Invalid ID", func(t *testing.T) {
		app.Put("/program/:id", handlers.UpdateProgramHandler)

		req := httptest.NewRequest("PUT", "/program/invalid", bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
	})
}

func TestIPPhoneHandlers(t *testing.T) {
	app := setupTestApp()

	t.Run("CreateIPPhoneHandler - Invalid JSON", func(t *testing.T) {
		app.Post("/ipphone", handlers.CreateIPPhoneHandler)

		req := httptest.NewRequest("POST", "/ipphone", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
	})

	t.Run("UpdateIPPhoneHandler - Invalid ID", func(t *testing.T) {
		app.Put("/ipphone/:id", handlers.UpdateIPPhoneHandler)

		req := httptest.NewRequest("PUT", "/ipphone/invalid", bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
	})
}

func TestTaskHandlers(t *testing.T) {
	app := setupTestApp()

	t.Run("CreateTaskHandler - Invalid JSON", func(t *testing.T) {
		app.Post("/task", handlers.CreateTaskHandler)

		req := httptest.NewRequest("POST", "/task", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
	})

	t.Run("UpdateTaskHandler - Invalid ID", func(t *testing.T) {
		app.Put("/task/:id", handlers.UpdateTaskHandler)

		req := httptest.NewRequest("PUT", "/task/invalid", bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
	})

	t.Run("GetTaskDetailHandler - Invalid ID", func(t *testing.T) {
		app.Get("/task/:id", handlers.GetTaskDetailHandler)

		req := httptest.NewRequest("GET", "/task/invalid", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
	})
}

func TestModelsValidation(t *testing.T) {
	t.Run("DepartmentRequest", func(t *testing.T) {
		name := "Test Department"
		branchID := 1
		req := models.DepartmentRequest{
			Name:     &name,
			BranchID: &branchID,
		}
		assert.Equal(t, "Test Department", *req.Name)
		assert.Equal(t, 1, *req.BranchID)
	})

	t.Run("TaskStatusUpdateRequest", func(t *testing.T) {
		req := models.TaskStatusUpdateRequest{
			ID:        1,
			Status:    1,
			UpdatedBy: 123,
		}
		assert.Equal(t, 1, req.ID)
		assert.Equal(t, 1, req.Status)
		assert.Equal(t, 123, req.UpdatedBy)
	})
}
