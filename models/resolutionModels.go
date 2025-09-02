package models

type ResolutionReq struct {
	Solution   string            `json:"solution"`
	FilePaths  map[string]string `json:"file_paths"`
	TelegramID int               `json:"telegram_id"`
	MessageID  int               `json:"message_id"`
	Url        string            `json:"url"`
	Assignto   string            `json:"assignto"`
	TicketNo   string            `json:"ticket_no"`
	CreatedAt  string            `json:"created_at"`
	ResolvedAt string            `json:"resolved_at"`
}

type UpdateResolutionReq struct {
	Solution   string            `json:"solution"`
	FilePaths  map[string]string `json:"file_paths"`
	TelegramID int               `json:"telegram_id"`
	ResolvedAt string            `json:"resolved_at"`
}
