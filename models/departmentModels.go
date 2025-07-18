package models

// Department model for displaying department information
type Department struct {
	ID         int     `json:"id"`
	Name       *string `json:"name"`
	BranchID   *int    `json:"branch_id"`
	BranchName *string `json:"branch_name"`
	CreatedAt  *string `json:"created_at"`
	UpdatedAt  *string `json:"updated_at"`
	DeletedAt  *string `json:"deleted_at"`
}

// Department Request model for creating or updating a department
type DepartmentRequest struct {
	Name     *string `json:"name"`
	BranchID *int    `json:"branch_id"`
}

// Department Detail model for detailed view of a department
type DepartmentDetail struct {
	ID            int     `json:"id"`
	Name          *string `json:"name"`
	BranchID      *int    `json:"branch_id"`
	BranchName    *string `json:"branch_name"`
	CreatedAt     *string `json:"created_at"`
	UpdatedAt     *string `json:"updated_at"`
	IPPhonesCount *int    `json:"ip_phones_count"`
	TasksCount    *int    `json:"tasks_count"`
}
