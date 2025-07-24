package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"reports-api/db"
	"reports-api/models"

	"github.com/gorilla/mux"
)

func ListScoresHandler(w http.ResponseWriter, r *http.Request) {
	query := `SELECT department_id, year, month, score FROM scores`
	rows, err := db.DB.Query(query)
	if err != nil {
		http.Error(w, "Failed to query scores", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var scores []models.Score
	for rows.Next() {
		var score models.Score
		err := rows.Scan(&score.DepartmentID, &score.Year, &score.Month, &score.Score)
		if err != nil {
			http.Error(w, "Failed to scan score", http.StatusInternalServerError)
			return
		}
		scores = append(scores, score)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, "Row error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"scores": scores})
}

func GetScoreDetailHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	query := `SELECT department_id, year, month, score FROM scores WHERE department_id = ?`
	row := db.DB.QueryRow(query, id)

	var score models.Score
	err := row.Scan(&score.DepartmentID, &score.Year, &score.Month, &score.Score)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Score not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to query score", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(score)
}

func UpdateScoreHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var score models.ScoreUpdateRequest
	err := json.NewDecoder(r.Body).Decode(&score)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var query string
	var args []interface{}

	if score.Year > 0 && score.Month > 0 {

		query = `UPDATE scores SET score = ? WHERE department_id = ? AND year = ? AND month = ?`
		args = []interface{}{score.Score, id, score.Year, score.Month}
	} else {

		query = `UPDATE scores SET score = ? WHERE department_id = ?`
		args = []interface{}{score.Score, id}
	}

	_, err = db.DB.Exec(query, args...)
	if err != nil {
		http.Error(w, "Failed to update score", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteScoreHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var score models.Score
	err := json.NewDecoder(r.Body).Decode(&score)

	var query string
	var args []interface{}

	if err == nil && score.Year > 0 && score.Month > 0 {

		query = `DELETE FROM scores WHERE department_id = ? AND year = ? AND month = ?`
		args = []interface{}{id, score.Year, score.Month}
	} else {

		query = `DELETE FROM scores WHERE department_id = ?`
		args = []interface{}{id}
	}

	_, err = db.DB.Exec(query, args...)
	if err != nil {
		http.Error(w, "Failed to delete score", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
