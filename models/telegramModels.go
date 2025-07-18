package models

// TelegramRequest model for sending Telegram messages
type TelegramRequest struct {
	Reportmessage  string `json:"reportmessage"`
	BranchName     string `json:"branchName"`
	DepartmentName string `json:"departmentName"`
	Number         string `json:"number"`
	IPPhoneName    string `json:"ipphoneName"`
	URL            string `json:"url"`
	// URLTs          string `json:"urlts"`
}

// TelegramResponse model for receiving Telegram responses
type TelegramResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
