package models

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"name"`
	CreatedBy int    `json:"created_by"`
}

type UserResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data"`
	Timestamp string      `json:"timestamp,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
}
