package models

// TaskRequest model for receiving task data
type TaskRequest struct {
	PhoneID        *int   `json:"phone_id" db:"phone_id"`
	Ticket         string `json:"ticket_no" db:"ticket_no"`
	SystemID       int    `json:"system_id"`
	IssueElse      string `json:"issue_else" db:"issue_else"`
	IssueTypeID    int    `json:"issue_type" db:"issue_type"`
	DepartmentID   int    `json:"department_id"`
	Text           string `json:"text"`
	ReportedBy     string `json:"reported_by"`
	Assignto       string `json:"assign_to"`
	Status         int    `json:"status"`
	CreatedBy      int    `json:"created_by"`
	UpdatedBy      int    `json:"updated_by"`
	Telegram       bool   `json:"telegram"`
	MessageID      int    `json:"message_id"`
	UpdatedAt      string `json:"updated_at"`
	PhoneNumber    int    `json:"-"`
	DepartmentName string `json:"-"`
	BranchName     string `json:"-"`
	ProgramName    string `json:"-"`
	Url            string `json:"-"`
	CreatedAt      string `json:"-"`
}

type TaskRequestUpdate struct {
	PhoneID      *int    `json:"phone_id" db:"phone_id"`
	SystemID     int     `json:"system_id"`
	IssueElse    string  `json:"issue_else" db:"issue_else"`
	IssueTypeID  int     `json:"issue_type" db:"issue_type"`
	Assignto     *string `json:"assign_to"`
	ReportedBy   *string `json:"reported_by"`
	DepartmentID int     `json:"department_id"`
	Status       int     `json:"status"`
	Text         string  `json:"text"`
	Solution     string  `json:"solution"`
	UpdatedBy    int     `json:"updated_by"`
}

type TaskStatusUpdateRequest struct {
	ID        int `json:"id" db:"id"`
	Status    int `json:"status"`
	UpdatedBy int `json:"updated_by"`
}

// TaskWithDetailsDb model for task with details in the database
type TaskWithDetails struct {
	ID                int               `json:"id" db:"id"`
	Ticket            string            `json:"ticket_no" db:"ticket_no"`
	PhoneID           *int              `json:"phone_id"`
	Number            *int              `json:"number"`
	PhoneName         *string           `json:"phone_name"`
	SystemID          int               `json:"system_id"`
	IssueElse         string            `json:"issue_else"`
	IssueTypeID       int               `json:"issue_type"`
	SystemName        string            `json:"system_name"`
	SystemType        string            `json:"system_type"`
	Solution          string            `json:"solution"`
	DepartmentID      int               `json:"department_id"`
	DepartmentName    string            `json:"department_name"`
	BranchID          int               `json:"branch_id"`
	BranchName        string            `json:"branch_name"`
	Text              string            `json:"text"`
	Assignto          *string           `json:"assign_to"`
	ReportedBy        *string           `json:"reported_by"`
	Status            int               `json:"status"`
	FilePaths         map[string]string `json:"file_paths"`
	FilePathsSolution map[string]string `json:"file_paths_solution"`
	ResolvedAt        string            `json:"resolved_at"`
	CreatedAt         string            `json:"created_at"`
	UpdatedAt         string            `json:"updated_at"`
	Overdue           interface{}       `json:"overdue"`
}
