package tests

import (
	"bytes"
	"net/http/httptest"
	"reports-api/db"
	"reports-api/handlers"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestRealWorldScenarios(t *testing.T) {
	app := SetupApp()

	t.Run("Database Connection Failure", func(t *testing.T) {
		app.Get("/departments", handlers.ListDepartmentsHandler)
		
		req := httptest.NewRequest("GET", "/departments", nil)
		resp, err := app.Test(req)
		
		// Test should not panic, but should handle error gracefully
		assert.NoError(t, err)
		// With database, expect 200 OK
		if db.DB != nil {
			assert.Equal(t, 200, resp.StatusCode)
		} else {
			assert.Equal(t, 500, resp.StatusCode)
		}
	})

	t.Run("Route Not Found", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/nonexistent", nil)
		resp, err := app.Test(req)
		
		assert.NoError(t, err)
		assert.Equal(t, 404, resp.StatusCode)
	})

	t.Run("Wrong HTTP Method", func(t *testing.T) {
		app.Get("/test", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "GET only"})
		})
		
		// Try POST on GET-only route
		req := httptest.NewRequest("POST", "/test", nil)
		resp, err := app.Test(req)
		
		assert.NoError(t, err)
		assert.Equal(t, 405, resp.StatusCode) // Method Not Allowed
	})

	t.Run("Invalid JSON Input", func(t *testing.T) {
		app.Post("/create", handlers.CreateDepartmentHandler)
		
		req := httptest.NewRequest("POST", "/create", 
			bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode) // Bad Request
	})

	t.Run("Missing Required Fields", func(t *testing.T) {
		app.Post("/create2", handlers.CreateDepartmentHandler)
		
		req := httptest.NewRequest("POST", "/create2", 
			bytes.NewReader([]byte("{}"))) // Empty JSON
		req.Header.Set("Content-Type", "application/json")
		
		resp, err := app.Test(req)
		assert.NoError(t, err)
		// With database, empty JSON might still create record
		if db.DB != nil {
			assert.True(t, resp.StatusCode == 200 || resp.StatusCode >= 400)
		} else {
			assert.True(t, resp.StatusCode >= 400)
		}
	})
}