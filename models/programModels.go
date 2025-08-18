package models

// Program model for displaying program information
type Program struct {
	ID        int     `json:"id"`
	Name      *string `json:"name"`
	TypeID    *int    `json:"type_id"`
	TypeName  string  `json:"type_name"`
	CreatedAt *string `json:"created_at"`
	UpdatedAt *string `json:"updated_at"`
	DeletedAt *string `json:"deleted_at"`
	CreatedBy *int    `json:"created_by"`
	UpdatedBy *int    `json:"updated_by"`
	DeletedBy *int    `json:"deleted_by"`
}

// ProgramRequest model for receiving program data
type ProgramRequest struct {
	Name      *string `json:"name"`
	TypeID    *int    `json:"type_id"`
	CreatedBy *int    `json:"created_by"`
	UpdatedBy *int    `json:"updated_by"`
}

type Type struct {
	ID        int     `json:"id"`
	Name      *string `json:"name"`
	CreatedAt *string `json:"created_at"`
}
