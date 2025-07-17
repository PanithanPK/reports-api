package models

// DashboardResponse สำหรับ response dashboard
// (ปรับให้ตรงกับ response ที่ dashboard.go ใช้)
type DashboardResponse struct {
	Success     bool                `json:"success"`
	Message     string              `json:"message"`
	Branches    []BranchDb          `json:"branches"`
	Departments []DepartmentDb      `json:"departments"`
	IPPhones    []IPPhoneDb         `json:"ip_phones"`
	Programs    []ProgramDb         `json:"programs"`
	Tasks       []TaskWithDetailsDb `json:"tasks"`
	ChartData   ChartData           `json:"chartdata"`
}

type TaskWithDetailsDb struct {
	ID             int     `json:"id"`
	PhoneID        int     `json:"phone_id"`
	Number         int     `json:"number"`
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

// ChartData สำหรับข้อมูลกราฟ
// (ใช้ใน calculateChartData)
type ChartData struct {
	YearStats []YearStat `json:"yearStats"`
}

// YearStat, MonthStat, BranchStat สำหรับโครงสร้างใหม่
// (ต้องตรงกับที่ใช้ใน dashboard.go)
type ProgramStat struct {
	ProgramId     *int   `json:"programId"`
	ProgramName   string `json:"programName"`
	TotalProblems int    `json:"total_problems"`
}

type IPPhoneStat struct {
	IPPhoneId     *int          `json:"ipphoneId"`
	IPPhoneName   string        `json:"ipphoneName"`
	TotalProblems int           `json:"total_problems"`
	Programs      []ProgramStat `json:"programs"`
}
type DepartmentStat struct {
	DepartmentId   *int          `json:"departmentId"`
	DepartmentName string        `json:"departmentName"`
	TotalProblems  int           `json:"total_problems"`
	IPPhones       []IPPhoneStat `json:"ipphones"`
}
type BranchStat struct {
	BranchId      *int             `json:"branchId"`
	BranchName    string           `json:"branchName"`
	TotalProblems int              `json:"total_problems"`
	Departments   []DepartmentStat `json:"departments"`
}
type MonthStat struct {
	Month    string       `json:"month"`
	Branches []BranchStat `json:"branches"`
}
type YearStat struct {
	Year   string      `json:"year"`
	Months []MonthStat `json:"months"`
}

type BarChartData struct {
	BranchName  string         `json:"branchName"`
	TotalCount  int            `json:"totalCount"`
	ProgramData map[string]int `json:"programData"`
}

type PieChartData struct {
	ProgramName string  `json:"programName"`
	Count       int     `json:"count"`
	Percentage  float64 `json:"percentage"`
}

type BranchDb struct {
	ID        int     `json:"id"`
	Name      *string `json:"name"`
	CreatedAt *string `json:"created_at"`
	UpdatedAt *string `json:"updated_at"`
}

type DepartmentDb struct {
	ID         int     `json:"id"`
	Name       *string `json:"name"`
	BranchID   *int    `json:"branch_id"`
	BranchName *string `json:"branch_name"`
	CreatedAt  *string `json:"created_at"`
	UpdatedAt  *string `json:"updated_at"`
}

type ProgramDb struct {
	ID        int     `json:"id"`
	Name      *string `json:"name"`
	CreatedAt *string `json:"created_at"`
	UpdatedAt *string `json:"updated_at"`
}

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
