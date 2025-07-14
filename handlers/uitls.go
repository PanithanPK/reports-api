package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reports-api/db"
)

// HealthCheckHandler provides a health check endpoint for Docker and monitoring
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Check DB connection
	dbStatus := "ok"
	if err := db.DB.Ping(); err != nil {
		dbStatus = "error"
	}

	handlersStatus := map[string]string{
		"ReportProblemHandler":    handlerStatus(checkHandlerHealth(ReportProblemHandler)),
		"GetProblemsHandler":      handlerStatus(checkHandlerHealth(GetProblemsHandler)),
		"AddUserHandler":          handlerStatus(checkHandlerHealth(AddUserHandler)),
		"GetUsersHandler":         handlerStatus(checkHandlerHealth(GetUsersHandler)),
		"AddProgramHandler":       handlerStatus(checkHandlerHealth(AddProgramHandler)),
		"GetProgramsHandler":      handlerStatus(checkHandlerHealth(GetProgramsHandler)),
		"AddBranchOfficeHandler":  handlerStatus(checkHandlerHealth(AddBranchOfficeHandler)),
		"GetBranchOfficesHandler": handlerStatus(checkHandlerHealth(GetBranchOfficesHandler)),
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"db":       dbStatus,
		"handlers": handlersStatus,
		"status":   "ok",
	})
}

func handlerStatus(ok bool) string {
	if ok {
		return "ok"
	}
	return "error"
}

func checkHandlerHealth(handlerFunc http.HandlerFunc) bool {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handlerFunc(w, req)
	resp := w.Result()
	return resp.StatusCode == http.StatusOK
}
