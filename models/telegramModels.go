package models

// TelegramRequest model for sending Telegram messages
type TelegramRequest struct {
	Reportmessage  string `json:"reportmessage"`
	BranchName     string `json:"branchName"`
	DepartmentName string `json:"departmentName"`
	Program        string `json:"program"`
	CreatedA       string `json:"createdA"`
	URL            string `json:"url"`

	// URLTs          string `json:"urlts"`
}

// TelegramResponse model for receiving Telegram responses
type TelegramResponse struct {
	Ok     bool `json:"ok"`
	Result struct {
		MessageID int `json:"message_id"`
	} `json:"result"`
}

type TelegramPayload struct {
	ChatID           string `json:"chat_id"`
	Text             string `json:"text"`
	ReplyToMessageID int    `json:"reply_to_message_id"`
	ParseMode        string `json:"parse_mode,omitempty"`
}
