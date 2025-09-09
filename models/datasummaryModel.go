package models

import (
	"time"
)

type DataTask struct {
	ID             int        `json:"id"`
	TicketNo       *string    `json:"ticket_no"`
	PhoneName      *string    `json:"phone_name"`
	IssueTypeName  *string    `json:"issue_type_name"`
	SystemName     *string    `json:"system_name"`
	BranchName     *string    `json:"branch_name"`
	DepartmentName *string    `json:"department_name"`
	Text           *string    `json:"text"`
	ReportedBy     *string    `json:"reported_by"`
	AssigntoName   *string    `json:"assignto_name"`
	SolutionText   *string    `json:"solution_text"`
	StatusText     *string    `json:"status_text"`
	FilePaths      *string    `json:"file_paths"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	ResolvedAt     *time.Time `json:"resolved_at"`
	DeletedAt      *time.Time `json:"deleted_at"`
	CreatedBy      *int       `json:"created_by"`
	UpdatedBy      *int       `json:"updated_by"`
}

type DataIssueType struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type DataSystemProgram struct {
	ID        int        `json:"id"`
	Name      *string    `json:"name"`
	Priority  *int       `json:"priority"`
	Type      *int       `json:"type"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	CreatedBy *int       `json:"created_by"`
	UpdatedBy *int       `json:"updated_by"`
	DeletedBy *int       `json:"deleted_by"`
}

type DataBranch struct {
	ID        int        `json:"id"`
	Name      *string    `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	CreatedBy *int       `json:"created_by"`
	UpdatedBy *int       `json:"updated_by"`
	DeletedBy *int       `json:"deleted_by"`
}

type DataDepartment struct {
	ID        int        `json:"id"`
	Name      *string    `json:"name"`
	BranchID  *int       `json:"branch_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type DataIPPhone struct {
	ID           int        `json:"id"`
	Number       *string    `json:"number"`
	Name         *string    `json:"name"`
	DepartmentID *int       `json:"department_id"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at"`
}

type DataResolution struct {
	ID         int        `json:"id"`
	TasksID    *int       `json:"tasks_id"`
	Text       *string    `json:"text"`
	TelegramID *string    `json:"telegram_id"`
	ResolvedAt *time.Time `json:"resolved_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
}

type DataResponsibility struct {
	ID               int       `json:"id"`
	TelegramUsername *string   `json:"telegram_username"`
	Name             *string   `json:"name"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
