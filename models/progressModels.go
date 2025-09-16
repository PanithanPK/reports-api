package models

type ProgressEntry struct {
	ID        int               `json:"id"`
	Text      string            `json:"text"`
	FilePaths map[string]string `json:"file_paths,omitempty"`
	CreatedAt string            `json:"created_at"`
}
