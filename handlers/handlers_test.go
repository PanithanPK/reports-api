package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reports-api/db"
	"strings"
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

func TestSolveProblemHandler_InvalidBody(t *testing.T) {
	t.Log("เริ่มทดสอบกรณี body ไม่ถูกต้อง")
	req := httptest.NewRequest("PUT", "/problemEntry/solveProblem", nil) // ไม่มี body
	w := httptest.NewRecorder()

	SolveProblemHandler(w, req)
	resp := w.Result()
	t.Logf("Status code ที่ได้: %d", resp.StatusCode)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400 Bad Request, got %d", resp.StatusCode)
	}
}

func TestSolveProblemHandler_MissingFields(t *testing.T) {
	body := `{"problemId":0,"solution":"","solutionUser":""}`
	req := httptest.NewRequest("PUT", "/problemEntry/solveProblem", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	SolveProblemHandler(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 OK for missing fields, got %d", resp.StatusCode)
	}
	// ตรวจสอบ response body
	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)
	if res["success"] != false {
		t.Errorf("Expected success=false, got %v", res["success"])
	}
}
