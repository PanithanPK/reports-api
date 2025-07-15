package models

// Branch สำหรับแสดงข้อมูล branches
type Branch struct {
	ID        int     `json:"id"`
	Name      *string `json:"name"`
	CreatedAt *string `json:"created_at"`
	UpdatedAt *string `json:"updated_at"`
	DeletedAt *string `json:"deleted_at"`
	CreatedBy *int    `json:"created_by"`
	UpdatedBy *int    `json:"updated_by"`
	DeletedBy *int    `json:"deleted_by"`
}

// BranchRequest สำหรับรับข้อมูลเพิ่ม/แก้ไขสาขา
type BranchRequest struct {
	Name      *string `json:"name"`
	CreatedBy *int    `json:"created_by"`
	UpdatedBy *int    `json:"updated_by"`
}
