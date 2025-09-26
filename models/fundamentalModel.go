package models

// Models Middleware
type StandardResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	Data      any    `json:"data,omitempty"`
	Error     any    `json:"error,omitempty"`
	Timestamp string `json:"timestamp"`
	RequestID string `json:"request_id,omitempty"`
}

// Models Config
type Config struct {
	EndPoint        string
	AccessKey       string
	SecretAccessKey string
	BucketName      string
	Environment     string
	ChatID          string
	BotToken        string
}

// Models ImageProcessor
type ImageConfig struct {
	MaxWidth    uint
	MaxHeight   uint
	Quality     int   // JPEG quality (1-100)
	MaxFileSize int64 // Maximum file size in bytes (e.g., 10MB for Telegram)
}
