package models

// Department สำหรับแสดงข้อมูล departments
type Department struct {
	ID         int     `json:"id"`
	Name       *string `json:"name"`
	BranchID   *int    `json:"branch_id"`
	BranchName *string `json:"branch_name"`
	CreatedAt  *string `json:"created_at"`
	UpdatedAt  *string `json:"updated_at"`
	DeletedAt  *string `json:"deleted_at"`
}

// DepartmentRequest สำหรับรับข้อมูลเพิ่ม/แก้ไขแผนก
type DepartmentRequest struct {
	Name     *string `json:"name"`
	BranchID *int    `json:"branch_id"`
}

// DepartmentDetail สำหรับแสดงข้อมูลแผนกพร้อมรายละเอียดเพิ่มเติม
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
