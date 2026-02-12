package dto

// ErrorResponse represents a standardized error response
// Matches the frontend expectation from BACKEND_REQUIREMENTS.md
type ErrorResponse struct {
	Error   string                 `json:"error"`             // Human-readable error message
	Code    string                 `json:"code,omitempty"`    // Optional error code
	Details map[string]interface{} `json:"details,omitempty"` // Optional additional context
}

// NewErrorResponse creates a simple error response with just a message
func NewErrorResponse(message string) *ErrorResponse {
	return &ErrorResponse{
		Error: message,
	}
}

// NewErrorResponseWithCode creates an error response with a code
func NewErrorResponseWithCode(message, code string) *ErrorResponse {
	return &ErrorResponse{
		Error: message,
		Code:  code,
	}
}

// NewErrorResponseWithDetails creates a full error response with code and details
func NewErrorResponseWithDetails(message, code string, details map[string]interface{}) *ErrorResponse {
	return &ErrorResponse{
		Error:   message,
		Code:    code,
		Details: details,
	}
}

// WithDetails adds details to an existing error response
func (e *ErrorResponse) WithDetails(details map[string]interface{}) *ErrorResponse {
	e.Details = details
	return e
}

// WithCode adds a code to an existing error response
func (e *ErrorResponse) WithCode(code string) *ErrorResponse {
	e.Code = code
	return e
}
