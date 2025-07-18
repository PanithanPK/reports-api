package models

// Branch model for displaying branch information
type Branch struct {
	ID        int     `json:"id"`
	Name      *string `json:"name"`
	CreatedAt *string `json:"created_at"`
	UpdatedAt *string `json:"updated_at"`
	DeletedAt *string `json:"deleted_at"`
	CreatedBy *int    `json:"created_by"`
	UpdatedBy *int    `json:"updated_by"`
	DeletedBy *int    `json:"deleted_by"`
}

// BranchRequest model for receiving branch data
type BranchRequest struct {
	Name      *string `json:"name"`
	CreatedBy *int    `json:"created_by"`
	UpdatedBy *int    `json:"updated_by"`
}

// BranchDetail model for displaying branch details
type BranchDetail struct {
	ID               int     `json:"id"`
	Name             *string `json:"name"`
	CreatedAt        *string `json:"created_at"`
	UpdatedAt        *string `json:"updated_at"`
	DepartmentsCount *int    `json:"departments_count"`
	IPPhonesCount    *int    `json:"ip_phones_count"`
}
