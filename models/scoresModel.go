package models

type Score struct {
	DepartmentID int `json:"department_id"`
	Year         int `json:"year"`
	Month        int `json:"month"`
	Score        int `json:"score"`
}

type ScoreUpdateRequest struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Score int `json:"score"`
}
