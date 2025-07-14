package models

import (
	"database/sql"
	"encoding/json"
)

type ProgramRequest struct {
	Name string `json:"name"`
}

type ProgramResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type GetProgramsResponse struct {
	Success  bool      `json:"success"`
	Message  string    `json:"message"`
	Programs []Program `json:"programs"`
}

// Program represents a program in the system
type Program struct {
	ID   sql.NullInt64  `json:"-"`
	Name sql.NullString `json:"-"`
}

// MarshalJSON implements custom JSON marshaling for Program
func (p *Program) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID   *int64  `json:"id"`
		Name *string `json:"name"`
	}{
		ID:   nullInt64ToPtr(p.ID),
		Name: nullStringToPtr(p.Name),
	})
}

type DeleteAllProgramsResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
