package handlers

import (
	"net/http"
	"net/http/httptest"
	"reports-api/db"
	"testing"
)

func TestDBConnection(t *testing.T) {
	err := db.InitDB()
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}
	if db.DB == nil {
		t.Fatal("DB instance is nil after InitDB")
	}
}

func TestHealthCheckHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	HealthCheckHandler(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", resp.StatusCode)
	}
}
