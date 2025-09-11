package models

type ResolutionReq struct {
	Solution     string            `json:"solution"`
	FilePaths    map[string]string `json:"file_paths"`
	ImageURLs    []string          `json:"image_urls"`
	TelegramID   int               `json:"telegram_id"`
	TelegramUser string            `json:"telegram_user"`
	MessageID    int               `json:"message_id"`
	Url          string            `json:"url"`
	Assignto     string            `json:"assignto"`
	AssignedtoID int               `json:"assignedto_id" db:"assignedto_id"`
	TicketNo     string            `json:"ticket_no"`
	CreatedAt    string            `json:"created_at"`
	ResolvedAt   string            `json:"resolved_at"`
}

type UpdateResolutionReq struct {
	Solution   string            `json:"solution"`
	FilePaths  map[string]string `json:"file_paths"`
	ImageURLs  []string          `json:"image_urls"`
	TelegramID int               `json:"telegram_id"`
	ResolvedAt string            `json:"resolved_at"`
}
