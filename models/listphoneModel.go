package models

// IPPhone สำหรับแสดงข้อมูล ip_phones
type IPPhone struct {
	ID             int     `json:"id"`
	Number         *int    `json:"number"`
	Name           *string `json:"name"`
	DepartmentID   int     `json:"department_id"`
	DepartmentName *string `json:"department_name"`
	BranchID       *int    `json:"branch_id"`
	BranchName     *string `json:"branch_name"`
	CreatedAt      *string `json:"created_at"`
	UpdatedAt      *string `json:"updated_at"`
	DeletedAt      *string `json:"deleted_at"`
	CreatedBy      *int    `json:"created_by"`
	UpdatedBy      *int    `json:"uodated_by"`
	DeletedBy      *int    `json:"deleted_by"`
}

// IPPhoneRequest สำหรับรับข้อมูลเพิ่ม/แก้ไข ip_phone
type IPPhoneRequest struct {
	Number       *int    `json:"number"`
	Name         *string `json:"name"`
	DepartmentID int     `json:"department_id"`
	CreatedBy    *int    `json:"created_by"`
	UpdatedBy    *int    `json:"uodated_by"`
}
