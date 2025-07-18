package models

// TaskRequest model for receiving task data
type TaskRequest struct {
	PhoneID   int    `json:"phone_id"`
	SystemID  int    `json:"system_id"`
	Text      string `json:"text"`
	Status    int    `json:"status"`
	CreatedBy int    `json:"created_by"`
	UpdatedBy int    `json:"updated_by"`
}

// TaskWithDetailsDb model for task with details in the database
type TaskWithDetails struct {
	ID             int    `json:"id"`
	PhoneID        int    `json:"phone_id"`
	Number         int    `json:"number"`
	PhoneName      string `json:"phone_name"`
	SystemID       int    `json:"system_id"`
	SystemName     string `json:"system_name"`
	DepartmentID   int    `json:"department_id"`
	DepartmentName string `json:"department_name"`
	BranchID       int    `json:"branch_id"`
	BranchName     string `json:"branch_name"`
	Text           string `json:"text"`
	Status         int    `json:"status"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}
