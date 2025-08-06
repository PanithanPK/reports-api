package tests

import (
	"net/http/httptest"
	"reports-api/handlers"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestRouteExistence(t *testing.T) {
	app := fiber.New()

	// Setup routes like in real application
	app.Get("/api/v1/department/list", handlers.ListDepartmentsHandler)
	app.Post("/api/v1/department/create", handlers.CreateDepartmentHandler)
	app.Put("/api/v1/department/update/:id", handlers.UpdateDepartmentHandler)
	app.Delete("/api/v1/department/delete/:id", handlers.DeleteDepartmentHandler)

	t.Run("Department List Route Exists", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/department/list", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		// Should not be 404 (route exists)
		assert.NotEqual(t, 404, resp.StatusCode)
	})

	t.Run("Wrong Route Returns 404", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/wrong/route", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 404, resp.StatusCode)
	})

	t.Run("Wrong HTTP Method Returns 405", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/department/list", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 405, resp.StatusCode)
	})
}

func TestDatabaseConnectionFailure(t *testing.T) {
	app := fiber.New()
	app.Get("/department/list", handlers.ListDepartmentsHandler)

	t.Run("Database Error Returns 500", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/department/list", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		// Without database connection, should return 500
		assert.Equal(t, 500, resp.StatusCode)
	})
}

func TestHandlerPanic(t *testing.T) {
	app := fiber.New()

	// Test handler that might panic
	app.Get("/panic", func(c *fiber.Ctx) error {
		panic("something went wrong")
	})

	t.Run("Handler Panic Should Be Caught", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/panic", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		// Fiber should catch panic and return 500
		assert.Equal(t, 500, resp.StatusCode)
	})
}