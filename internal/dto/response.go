package dto

import "time"

// ApiResponse is the standardized response wrapper for all API endpoints
type ApiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ApiError   `json:"error,omitempty"`
	Meta    *Metadata   `json:"meta"`
}

// ApiError represents error information in the response
type ApiError struct {
	Message string                 `json:"message"`
	Code    string                 `json:"code,omitempty"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// Metadata contains request tracking and timestamp information
type Metadata struct {
	Timestamp string `json:"timestamp"`
	RequestID string `json:"request_id,omitempty"`
}

// NewApiSuccessResponse creates a standardized success response
func NewApiSuccessResponse(data interface{}, message string) *ApiResponse {
	if message == "" {
		message = "Operation completed successfully"
	}
	return &ApiResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta: &Metadata{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	}
}

// NewApiErrorResponse creates a standardized error response
func NewApiErrorResponse(message, code string, details map[string]interface{}) *ApiResponse {
	return &ApiResponse{
		Success: false,
		Error: &ApiError{
			Message: message,
			Code:    code,
			Details: details,
		},
		Meta: &Metadata{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	}
}

// NewApiErrorResponseFromError creates error response from ErrorResponse
func NewApiErrorResponseFromError(err *ErrorResponse) *ApiResponse {
	return &ApiResponse{
		Success: false,
		Error: &ApiError{
			Message: err.Error,
			Code:    err.Code,
			Details: err.Details,
		},
		Meta: &Metadata{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	}
}

// WithRequestID adds request ID to the response metadata
func (r *ApiResponse) WithRequestID(requestID string) *ApiResponse {
	if r.Meta != nil {
		r.Meta.RequestID = requestID
	}
	return r
}

// WithMessage sets or updates the message
func (r *ApiResponse) WithMessage(message string) *ApiResponse {
	r.Message = message
	return r
}
