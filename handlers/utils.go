package handlers

import (
	"log"
	"os"
	"time"
)

// Logger levels
const (
	INFO  = "INFO"
	WARN  = "WARN"
	ERROR = "ERROR"
	DEBUG = "DEBUG"
)

// CustomLogger provides structured logging with levels
type CustomLogger struct {
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger
	Debug *log.Logger
}

// Global logger instances
var (
	ProblemLogger   *CustomLogger
	SolutionLogger  *CustomLogger
	ProgramLogger   *CustomLogger
	BranchLogger    *CustomLogger
	UserLogger      *CustomLogger
	DashboardLogger *CustomLogger
)

// Initialize loggers
func init() {
	ProblemLogger = NewLogger("PROBLEM")
	SolutionLogger = NewLogger("SOLUTION")
	ProgramLogger = NewLogger("PROGRAM")
	BranchLogger = NewLogger("BRANCH")
	UserLogger = NewLogger("USER")
	DashboardLogger = NewLogger("DASHBOARD")
}

// NewLogger creates a new custom logger with the given prefix
func NewLogger(prefix string) *CustomLogger {
	return &CustomLogger{
		Info:  log.New(os.Stdout, "["+prefix+"][INFO] ", log.Ldate|log.Ltime),
		Warn:  log.New(os.Stdout, "["+prefix+"][WARN] ", log.Ldate|log.Ltime),
		Error: log.New(os.Stderr, "["+prefix+"][ERROR] ", log.Ldate|log.Ltime),
		Debug: log.New(os.Stdout, "["+prefix+"][DEBUG] ", log.Ldate|log.Ltime),
	}
}

// LogRequest logs incoming HTTP requests
func LogRequest(logger *CustomLogger, method, path, remoteAddr string) {
	logger.Info.Printf("📥 %s %s from %s", method, path, remoteAddr)
}

// LogResponse logs HTTP responses with timing
func LogResponse(logger *CustomLogger, method, path string, statusCode int, duration time.Duration) {
	emoji := "✅"
	if statusCode >= 400 {
		emoji = "❌"
	} else if statusCode >= 300 {
		emoji = "⚠️"
	}

	logger.Info.Printf("%s %s %s completed in %v (Status: %d)", emoji, method, path, duration, statusCode)
}

// LogDatabaseOperation logs database operations
func LogDatabaseOperation(logger *CustomLogger, operation, table string, rowsAffected int64, err error) {
	if err != nil {
		logger.Error.Printf("❌ Database %s on %s failed: %v", operation, table, err)
	} else {
		logger.Info.Printf("✅ Database %s on %s successful (%d rows affected)", operation, table, rowsAffected)
	}
}

// LogValidationError logs validation errors
func LogValidationError(logger *CustomLogger, field, value string, remoteAddr string) {
	logger.Warn.Printf("⚠️ Validation error - Field: %s, Value: %s, From: %s", field, value, remoteAddr)
}

// LogSecurityEvent logs security-related events
func LogSecurityEvent(logger *CustomLogger, event, details, remoteAddr string) {
	logger.Warn.Printf("🔒 Security event - %s: %s from %s", event, details, remoteAddr)
}

// LogPerformance logs performance metrics
func LogPerformance(logger *CustomLogger, operation string, duration time.Duration, details string) {
	if duration > 1*time.Second {
		logger.Warn.Printf("🐌 Slow operation - %s took %v: %s", operation, duration, details)
	} else {
		logger.Debug.Printf("⚡ %s completed in %v: %s", operation, duration, details)
	}
}

// LogDataProcessing logs data processing operations
func LogDataProcessing(logger *CustomLogger, operation string, count int, details string) {
	logger.Info.Printf("📊 %s processed %d records: %s", operation, count, details)
}

// LogError logs errors with context
func LogError(logger *CustomLogger, operation string, err error, context string) {
	logger.Error.Printf("❌ %s failed: %v | Context: %s", operation, err, context)
}

// LogSuccess logs successful operations
func LogSuccess(logger *CustomLogger, operation string, details string) {
	logger.Info.Printf("✅ %s successful: %s", operation, details)
}

// LogWarning logs warnings
func LogWarning(logger *CustomLogger, operation string, details string) {
	logger.Warn.Printf("⚠️ %s warning: %s", operation, details)
}

// LogInfo logs informational messages
func LogInfo(logger *CustomLogger, operation string, details string) {
	logger.Info.Printf("ℹ️ %s: %s", operation, details)
}

// LogDebug logs debug information
func LogDebug(logger *CustomLogger, operation string, details string) {
	logger.Debug.Printf("🔍 %s: %s", operation, details)
}
