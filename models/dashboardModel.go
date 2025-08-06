package models

// DashboardResponse model for dashboard response
type DashboardResponse struct {
	Success     bool                `json:"success"`
	Message     string              `json:"message"`
	ChartData   ChartData           `json:"chartdata"`
	Branches    []BranchDb          `json:"branches"`
	Departments []DepartmentDb      `json:"departments"`
	IPPhones    []IPPhoneDb         `json:"ip_phones"`
	Programs    []ProgramDb         `json:"programs"`
	Tasks       []TaskWithDetailsDb `json:"tasks"`
}

// TaskWithDetailsDb model for task with details
type TaskWithDetailsDb struct {
	ID             int     `json:"id"`
	PhoneID        *int    `json:"phone_id"`
	Number         *int    `json:"number"`
	PhoneName      *string `json:"phone_name"`
	SystemID       *int    `json:"system_id"`
	SystemName     *string `json:"system_name"`
	DepartmentID   *int    `json:"department_id"`
	DepartmentName *string `json:"department_name"`
	BranchID       *int    `json:"branch_id"`
	BranchName     *string `json:"branch_name"`
	Text           string  `json:"text"`
	Status         int     `json:"status"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
	Month          string  `json:"month"`
	Year           string  `json:"year"`
}

// ChartData model for chart data
type ChartData struct {
	YearStats []YearStat `json:"yearStats"`
}

// YearStat model for year statistics
type ProgramStat struct {
	ProgramId     *int   `json:"programId"`
	ProgramName   string `json:"programName"`
	TotalProblems int    `json:"total_problems"`
}

// IPPhoneStat model for IP phone statistics
type IPPhoneStat struct {
	IPPhoneId     *int          `json:"ipphoneId"`
	IPPhoneName   *string       `json:"ipphoneName"`
	TotalProblems int           `json:"total_problems"`
	Programs      []ProgramStat `json:"programs"`
}

// DepartmentStat model for department statistics
type DepartmentStat struct {
	DepartmentId   *int          `json:"departmentId"`
	DepartmentName string        `json:"departmentName"`
	TotalProblems  int           `json:"total_problems"`
	IPPhones       []IPPhoneStat `json:"ipphones"`
}

// BranchStat model for branch statistics
type BranchStat struct {
	BranchId      *int             `json:"branchId"`
	BranchName    string           `json:"branchName"`
	TotalProblems int              `json:"total_problems"`
	Departments   []DepartmentStat `json:"departments"`
}

// MonthStat model for month statistics
type MonthStat struct {
	Month    string       `json:"month"`
	Branches []BranchStat `json:"branches"`
}

// ChartData model for chart data
type YearStat struct {
	Year   string      `json:"year"`
	Months []MonthStat `json:"months"`
}

// BranchDb model for branch database representation
type BranchDb struct {
	ID        int     `json:"id"`
	Name      *string `json:"name"`
	CreatedAt *string `json:"created_at"`
	UpdatedAt *string `json:"updated_at"`
}

// DepartmentDb model for department database representation
type DepartmentDb struct {
	ID         int     `json:"id"`
	Name       *string `json:"name"`
	BranchID   *int    `json:"branch_id"`
	BranchName *string `json:"branch_name"`
	CreatedAt  *string `json:"created_at"`
	UpdatedAt  *string `json:"updated_at"`
	Scores     int     `json:"scores"`
}

// ProgramDb model for program database representation
type ProgramDb struct {
	ID        int     `json:"id"`
	Name      *string `json:"name"`
	CreatedAt *string `json:"created_at"`
	UpdatedAt *string `json:"updated_at"`
}

// IPPhoneDb model for IP phone database representation
type IPPhoneDb struct {
	ID             int     `json:"id"`
	Number         *int    `json:"number"`
	Name           *string `json:"name"`
	DepartmentID   *int    `json:"department_id"`
	DepartmentName *string `json:"department_name"`
	BranchID       *int    `json:"branch_id"`
	BranchName     *string `json:"branch_name"`
	CreatedAt      *string `json:"created_at"`
	UpdatedAt      *string `json:"updated_at"`
}
