package models

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"name"`
	CreatedBy int    `json:"created_by"`
}

type UserResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}
