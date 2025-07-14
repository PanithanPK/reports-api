package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"reports-api/db"
	"reports-api/models"
	"time"

	"github.com/gorilla/mux"
)

// SolveProblemHandler handles PUT requests to solve problems
func SolveProblemHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req models.SolveProblemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.ProblemID == 0 || req.Solution == "" || req.SolutionUser == "" {
		response := models.SolveProblemResponse{
			Success: false,
			Message: "Missing required fields: problemId, solution, solutionUser",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Set solution date to current time if not provided
	solutionDate := sql.NullTime{Time: time.Now(), Valid: true}
	if req.SolutionDate.Valid {
		solutionDate = req.SolutionDate
	}

	// Update database
	query := `UPDATE report_problem SET solution = ?, solution_date = ?, solution_user = ?, status = 'Solved', updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	result, err := db.DB.Exec(query, req.Solution, solutionDate, req.SolutionUser, req.ProblemID)
	if err != nil {
		log.Printf("Error updating problem: %v", err)
		response := models.SolveProblemResponse{
			Success: false,
			Message: "Failed to solve problem",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		response := models.SolveProblemResponse{
			Success: false,
			Message: "Problem not found",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := models.SolveProblemResponse{
		Success: true,
		Message: "Problem solved successfully",
	}
	json.NewEncoder(w).Encode(response)
}

func ResetSolutionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	problemID := vars["id"]

	// Reset solution fields in database
	query := `UPDATE report_problem SET solution = NULL, solution_date = NULL, solution_user = NULL, status = 'Pending', updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	result, err := db.DB.Exec(query, problemID)
	if err != nil {
		log.Printf("Error resetting solution: %v", err)
		response := models.ResetSolutionResponse{
			Success: false,
			Message: "Failed to reset solution",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		response := models.ResetSolutionResponse{
			Success: false,
			Message: "Problem not found",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := models.ResetSolutionResponse{
		Success: true,
		Message: "Solution reset successfully",
	}
	json.NewEncoder(w).Encode(response)
}

// UpdateSolutionHandler handles PUT requests to update solution details
func UpdateSolutionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	problemID := vars["id"]

	var req models.UpdateSolutionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Solution == "" || req.SolutionUser == "" {
		response := models.UpdateSolutionResponse{
			Success: false,
			Message: "Missing required fields: solution, solutionUser",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Update solution details in database
	query := `UPDATE report_problem SET solution = ?, solution_user = ?, solution_date = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	result, err := db.DB.Exec(query, req.Solution, req.SolutionUser, problemID)
	if err != nil {
		log.Printf("Error updating solution: %v", err)
		response := models.UpdateSolutionResponse{
			Success: false,
			Message: "Failed to update solution",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		response := models.UpdateSolutionResponse{
			Success: false,
			Message: "Problem not found",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := models.UpdateSolutionResponse{
		Success: true,
		Message: "Solution updated successfully",
	}
	json.NewEncoder(w).Encode(response)
}

func UpdateProblemANDSolveProblemHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	problemID := vars["id"]

	var req models.UpdateProblemANDSolveProblemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update problem details in database
	query := `UPDATE report_problem SET problem = ?, solution = ?, solution_user = ?, solution_date = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	result, err := db.DB.Exec(query, req.Problem, req.Solution, req.SolutionUser, problemID)
	if err != nil {
		log.Printf("Error updating problem: %v", err)
		response := models.UpdateProblemANDSolveProblemResponse{
			Success: false,
			Message: "Failed to update problem",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		response := models.UpdateProblemANDSolveProblemResponse{
			Success: false,
			Message: "Problem not found",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := models.UpdateProblemANDSolveProblemResponse{
		Success: true,
		Message: "Problem updated successfully",
	}
	json.NewEncoder(w).Encode(response)
}

func DeleteProblemandSolveProblemHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	problemID := vars["id"]

	// Delete problem and solution from database
	query := `DELETE FROM report_problem WHERE id = ?`
	result, err := db.DB.Exec(query, problemID)
	if err != nil {
		log.Printf("Error deleting problem: %v", err)
		response := models.DeleteProblemandSolveProblemResponse{
			Success: false,
			Message: "Failed to delete problem",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		response := models.DeleteProblemandSolveProblemResponse{
			Success: false,
			Message: "Problem not found",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := models.DeleteProblemandSolveProblemResponse{
		Success: true,
		Message: "Problem deleted successfully",
	}
	json.NewEncoder(w).Encode(response)
}
