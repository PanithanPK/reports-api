package models

// PaginationRequest represents pagination parameters
type PaginationRequest struct {
	Page  int `json:"page" query:"page"`
	Limit int `json:"limit" query:"limit"`
}

// PaginationResponse represents pagination metadata
type PaginationResponse struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Success    bool               `json:"success"`
	Data       interface{}        `json:"data"`
	Pagination PaginationResponse `json:"pagination"`
}