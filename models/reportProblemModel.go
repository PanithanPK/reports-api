package models

// TaskRequest model for receiving task data
type TaskRequest struct {
	PhoneID          *int    `json:"phone_id" db:"phone_id"`
	PhoneElse        *string `json:"phone_else" db:"phone_else"`
	Ticket           string  `json:"ticket_no" db:"ticket_no"`
	SystemID         int     `json:"system_id"`
	IssueElse        string  `json:"issue_else" db:"issue_else"`
	IssueTypeID      int     `json:"issue_type" db:"issue_type"`
	DepartmentID     int     `json:"department_id"`
	Text             string  `json:"text"`
	ReportedBy       string  `json:"reported_by"`
	AssignedtoID     int     `json:"assignedto_id" db:"assignedto_id"`
	Assignto         string  `json:"assign_to"`
	Status           int     `json:"status"`
	CreatedBy        int     `json:"created_by"`
	UpdatedBy        int     `json:"updated_by"`
	ResolvedAt       string  `json:"resolved_at"`
	Telegram         bool    `json:"telegram"`
	TelegramUser     string  `json:"telegram_user"`
	MessageID        int     `json:"message_id"`
	UpdatedAt        string  `json:"updated_at"`
	PreviousAssignto string  `json:"previous_assignto"`
	AssigntoID       int     `json:"-"`
	PhoneNumber      int     `json:"-"`
	DepartmentName   string  `json:"-"`
	BranchName       string  `json:"-"`
	ProgramName      string  `json:"-"`
	Url              string  `json:"-"`
	CreatedAt        string  `json:"-"`
}

type TaskRequestUpdate struct {
	PhoneID      *int    `json:"phone_id" db:"phone_id"`
	PhoneElse    *string `json:"phone_else" db:"phone_else"`
	SystemID     int     `json:"system_id"`
	IssueElse    string  `json:"issue_else" db:"issue_else"`
	IssueTypeID  int     `json:"issue_type" db:"issue_type"`
	AssignedtoID int     `json:"assignedto_id" db:"assignedto_id"`
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
	ID             int               `json:"id" db:"id"`
	Ticket         string            `json:"ticket_no" db:"ticket_no"`
	PhoneID        *int              `json:"phone_id"`
	PhoneElse      *string           `json:"phone_else"`
	Number         *int              `json:"number"`
	PhoneName      *string           `json:"phone_name"`
	SystemID       int               `json:"system_id"`
	IssueElse      string            `json:"issue_else"`
	IssueTypeID    int               `json:"issue_type"`
	SystemName     string            `json:"system_name"`
	SystemType     string            `json:"system_type"`
	DepartmentID   int               `json:"department_id"`
	DepartmentName string            `json:"department_name"`
	BranchID       int               `json:"branch_id"`
	BranchName     string            `json:"branch_name"`
	Text           string            `json:"text"`
	AssignedtoID   int               `json:"assignedto_id" db:"assignedto_id"`
	Assignto       *string           `json:"assign_to"`
	ReportedBy     *string           `json:"reported_by"`
	Status         int               `json:"status"`
	FilePaths      map[string]string `json:"file_paths"`
	ResolvedAt     string            `json:"resolved_at"`
	CreatedAt      string            `json:"created_at"`
	UpdatedAt      string            `json:"updated_at"`
	Overdue        interface{}       `json:"overdue"`
}

type AssignRequest struct {
	AssignedtoID   int    `json:"assignedto_id"`
	Assignto       string `json:"assign_to"`
	UpdatedBy      int    `json:"updated_by"`
	UpdateTelegram bool   `json:"update_telegram"`
}
