package constants

// Task Status Constants
const (
	TaskStatusPending    = 0
	TaskStatusInProgress = 1
	TaskStatusResolved   = 2
)

// HTTP Status Messages
const (
	MessageSuccess         = "Success"
	MessageCreated         = "Created successfully"
	MessageUpdated         = "Updated successfully"
	MessageDeleted         = "Deleted successfully"
	MessageBadRequest      = "Bad Request"
	MessageUnauthorized    = "Unauthorized"
	MessageNotFound        = "Not Found"
	MessageInternalError   = "Internal Server Error"
	MessageDatabaseError   = "Database error"
	MessageInvalidRequest  = "Invalid request body"
	MessageInvalidID       = "Invalid id"
)

// Error Messages
const (
	ErrorFailedToQuery   = "Failed to query"
	ErrorFailedToCreate  = "Failed to create"
	ErrorFailedToUpdate  = "Failed to update"
	ErrorFailedToDelete  = "Failed to delete"
	ErrorNotFound        = "Not found"
	ErrorInvalidRequest  = "Invalid request"
	ErrorDatabaseError   = "Database error"
)

// Pagination Defaults
const (
	DefaultPage  = 1
	DefaultLimit = 10
	MaxLimit     = 100
)
