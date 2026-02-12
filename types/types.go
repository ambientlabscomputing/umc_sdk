package types

// ComponentState represents the runtime state of a UMC
type ComponentState string

const (
	ComponentStateStarting ComponentState = "starting"
	ComponentStateRunning  ComponentState = "running"
	ComponentStateStopping ComponentState = "stopping"
	ComponentStateStopped  ComponentState = "stopped"
	ComponentStateFailed   ComponentState = "failed"
)

// HealthStatus represents health check result
type HealthStatus struct {
	Healthy bool                `json:"healthy"`
	Status  ComponentState      `json:"status"`
	Message string              `json:"message,omitempty"`
	Details map[string]string   `json:"details,omitempty"`
}

// ErrorResponse is a standard error response structure
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// NewError creates a new error response
func NewError(code, message string) *ErrorResponse {
	return &ErrorResponse{
		Code:    code,
		Message: message,
	}
}

// NewErrorWithDetails creates a new error response with details
func NewErrorWithDetails(code, message, details string) *ErrorResponse {
	return &ErrorResponse{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// Error implements the error interface
func (e *ErrorResponse) Error() string {
	if e.Details != "" {
		return e.Code + ": " + e.Message + " (" + e.Details + ")"
	}
	return e.Code + ": " + e.Message
}

// Common error codes for UMC operations
const (
	ErrCodeUnknown            = "UNKNOWN"
	ErrCodeInvalidRequest     = "INVALID_REQUEST"
	ErrCodeNotFound           = "NOT_FOUND"
	ErrCodeAlreadyExists      = "ALREADY_EXISTS"
	ErrCodePermissionDenied   = "PERMISSION_DENIED"
	ErrCodeUnavailable        = "UNAVAILABLE"
	ErrCodeTimeout            = "TIMEOUT"
	ErrCodeInternal           = "INTERNAL_ERROR"
	ErrCodeNotImplemented     = "NOT_IMPLEMENTED"
)
