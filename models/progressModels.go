package models

type ProgressEntry struct {
	ID        int               `json:"id"`
	Ticketno  string            `json:"ticket_no"`
	Text      string            `json:"text"`
	FilePaths map[string]string `json:"file_paths,omitempty"`
	AssignTo  string            `json:"assignto"`
	ImageURLs []string          `json:"-"`
	UpdateAt  string            `json:"updated_at"`
	CreatedAt string            `json:"created_at"`
}

type UpdateProgress struct {
	ID        int               `json:"id"`
	Text      string            `json:"text"`
	FilePaths map[string]string `json:"file_paths,omitempty"`
	ImageURLs []string          `json:"image_urls"`
	UpdateAt  string            `json:"updated_at"`
	CreatedAt string            `json:"created_at"`
}
