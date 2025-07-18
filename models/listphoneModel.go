package models

// IPPhone model for displaying IP phone information
type IPPhone struct {
	ID             int     `json:"id"`
	Number         *int    `json:"number"`
	Name           *string `json:"name"`
	DepartmentID   int     `json:"department_id"`
	DepartmentName *string `json:"department_name"`
	BranchID       *int    `json:"branch_id"`
	BranchName     *string `json:"branch_name"`
	CreatedAt      *string `json:"created_at"`
	UpdatedAt      *string `json:"updated_at"`
	DeletedAt      *string `json:"deleted_at"`
	CreatedBy      *int    `json:"created_by"`
	UpdatedBy      *int    `json:"uodated_by"`
	DeletedBy      *int    `json:"deleted_by"`
}

// IPPhoneRequest model for creating or updating an IP phone
type IPPhoneRequest struct {
	Number       *int    `json:"number"`
	Name         *string `json:"name"`
	DepartmentID int     `json:"department_id"`
	CreatedBy    *int    `json:"created_by"`
	UpdatedBy    *int    `json:"uodated_by"`
}
