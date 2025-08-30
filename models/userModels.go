package models

// Data model for user information
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

// LoginResponse represents the structure of the response after a user logs in
type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
	Data    *Data  `json:"data,omitempty"`
}

// RegisterUserRequest represents the structure of the request to register a new user
type RegisterUserRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	CreatedBy int    `json:"created_by"`
}

// UserDetailResponse represents the structure of user details response
type UpdateUserRequest struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	UpdatedBy int    `json:"updated_by"`
}

// UserDetailResponse represents the structure of user details response
type DeleteUserRequest struct {
	ID        int `json:"id"`
	DeletedBy int `json:"deleted_by"`
}

type ResponseRequest struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	TelegramUsername string `json:"telegram_username"`
}
