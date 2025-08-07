package models

// TaskRequest model for receiving task data
type TaskRequest struct {
	PhoneID        *int   `json:"phone_id" db:"phone_id"`
	Ticket         string `json:"ticket_no" db:"ticket_no"`
	SystemID       int    `json:"system_id"`
	DepartmentID   int    `json:"department_id"`
	Text           string `json:"text"`
	Status         int    `json:"status"`
	CreatedBy      int    `json:"created_by"`
	UpdatedBy      int    `json:"updated_by"`
	Telegram       bool   `json:"telegram"`
	PhoneNumber    int    `json:"-"`
	DepartmentName string `json:"-"`
	BranchName     string `json:"-"`
	ProgramName    string `json:"-"`
	Url            string `json:"-"`
	CreatedAt      string `json:"-"`
}

type TaskStatusUpdateRequest struct {
	ID        int `json:"id" db:"id"`
	Status    int `json:"status"`
	UpdatedBy int `json:"updated_by"`
}

// TaskWithDetailsDb model for task with details in the database
type TaskWithDetails struct {
	ID             int         `json:"id" db:"id"`
	Ticket         string      `json:"ticket_no" db:"ticket_no"`
	PhoneID        *int        `json:"phone_id"`
	Number         *int        `json:"number"`
	PhoneName      *string     `json:"phone_name"`
	SystemID       int         `json:"system_id"`
	SystemName     string      `json:"system_name"`
	DepartmentID   int         `json:"department_id"`
	DepartmentName string      `json:"department_name"`
	BranchID       int         `json:"branch_id"`
	BranchName     string      `json:"branch_name"`
	Text           string      `json:"text"`
	Status         int         `json:"status"`
	CreatedAt      string      `json:"created_at"`
	UpdatedAt      string      `json:"updated_at"`
	Overdue        interface{} `json:"overdue"`
}
