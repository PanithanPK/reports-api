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
