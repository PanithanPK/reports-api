package models

// TaskWithDetails สำหรับแสดง tasks พร้อมรายละเอียดจาก ip_phone และ systems_program
type TaskWithDetails struct {
	ID             int    `json:"id"`
	PhoneID        int    `json:"phone_id"`
	PhoneName      string `json:"phone_name"`
	SystemID       int    `json:"system_id"`
	SystemName     string `json:"system_name"`
	DepartmentID   int    `json:"department_id"`
	DepartmentName string `json:"department_name"`
	BranchID       int    `json:"branch_id"`
	BranchName     string `json:"branch_name"`
	Text           string `json:"text"`
	Status         string `json:"status"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}
