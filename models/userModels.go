package models

// User สำหรับแสดงข้อมูลผู้ใช้

type Data struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role,omitempty"`
}

// Credentials represents the structure of user login credentials
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse สำหรับตอบกลับ login
type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
	Data    *Data  `json:"data,omitempty"`
}

// RegisterUserRequest สำหรับรับข้อมูลสมัครสมาชิก
type RegisterUserRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	CreatedBy int    `json:"created_by"`
}

type UpdateUserRequest struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	UpdatedBy int    `json:"updated_by"`
}

type DeleteUserRequest struct {
	ID        int `json:"id"`
	DeletedBy int `json:"deleted_by"`
}
