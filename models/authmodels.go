package models

// Credentials represents the structure of user login credentials
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// User สำหรับแสดงข้อมูลผู้ใช้
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

// RegisterRequest สำหรับรับข้อมูลสมัครสมาชิก
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// UpdateUserRequest สำหรับรับข้อมูลแก้ไขผู้ใช้
type UpdateUserRequest struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// LogoutRequest สำหรับรับข้อมูล logout
type LogoutRequest struct {
	ID int `json:"id"`
}

// LoginResponse สำหรับตอบกลับ login
type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
	User    *User  `json:"user,omitempty"`
}
