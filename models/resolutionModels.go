package models

type ResolutionReq struct {
	Solution   string            `json:"solution"`
	FilePaths  map[string]string `json:"file_paths"`
	TelegramID int               `json:"telegram_id"`
	ResolvedAt string            `json:"resolved_at"`
}
