package errors

import (
	"fmt"
	"net/http"
)

// AppError represents an application error with HTTP status code and error code
type AppError struct {
	Message    string                 // Human-readable error message
	Code       string                 // Error code for client
	StatusCode int                    // HTTP status code
	Details    map[string]interface{} // Optional additional details
	Err        error                  // Original error (for logging)
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the wrapped error
func (e *AppError) Unwrap() error {
	return e.Err
}

// WithDetails adds details to the error
func (e *AppError) WithDetails(details map[string]interface{}) *AppError {
	e.Details = details
	return e
}

// WithError wraps an underlying error
func (e *AppError) WithError(err error) *AppError {
	e.Err = err
	return e
}

// Error Code Constants
// Based on BACKEND_REQUIREMENTS.md error handling section
const (
	// Authentication Errors
	ErrCodeUnauthorized  = "UNAUTHORIZED"
	ErrCodeInvalidToken  = "INVALID_TOKEN"
	ErrCodeForbidden     = "FORBIDDEN"
	ErrCodeExpiredToken  = "EXPIRED_TOKEN"

	// Validation Errors
	ErrCodeValidation    = "VALIDATION_ERROR"
	ErrCodeRequiredField = "REQUIRED_FIELD"
	ErrCodeInvalidFormat = "INVALID_FORMAT"
	ErrCodeInvalidValue  = "INVALID_VALUE"
	ErrCodeInvalidInput  = "INVALID_INPUT"

	// Business Logic Errors
	ErrCodeSlotUnavailable = "SLOT_UNAVAILABLE"
	ErrCodePackageNotFound = "PACKAGE_NOT_FOUND"
	ErrCodeThemeNotFound   = "THEME_NOT_FOUND"
	ErrCodeAddonNotFound   = "ADDON_NOT_FOUND"
	ErrCodeBookingNotFound = "BOOKING_NOT_FOUND"
	ErrCodeUserNotFound    = "USER_NOT_FOUND"
	ErrCodeInvalidSlotCount = "INVALID_SLOT_COUNT"
	ErrCodeResourceNotFound = "RESOURCE_NOT_FOUND"
	ErrCodeConflict         = "CONFLICT"
	ErrCodeAlreadyExists    = "ALREADY_EXISTS"

	// File Upload Errors
	ErrCodeFileTooLarge     = "FILE_TOO_LARGE"
	ErrCodeInvalidFileType  = "INVALID_FILE_TYPE"
	ErrCodeUploadFailed     = "UPLOAD_FAILED"

	// Server Errors
	ErrCodeInternal       = "INTERNAL_ERROR"
	ErrCodeDatabaseError  = "DATABASE_ERROR"
	ErrCodeEmailSendFailed = "EMAIL_SEND_FAILED"
)

// ============================================================================
// Authentication Errors (401, 403)
// ============================================================================

// NewUnauthorizedError creates a 401 Unauthorized error
func NewUnauthorizedError(message string) *AppError {
	if message == "" {
		message = "Authentication required"
	}
	return &AppError{
		Message:    message,
		Code:       ErrCodeUnauthorized,
		StatusCode: http.StatusUnauthorized,
	}
}

// NewInvalidTokenError creates a 401 Invalid Token error
func NewInvalidTokenError(message string) *AppError {
	if message == "" {
		message = "Token is invalid or expired"
	}
	return &AppError{
		Message:    message,
		Code:       ErrCodeInvalidToken,
		StatusCode: http.StatusUnauthorized,
	}
}

// NewForbiddenError creates a 403 Forbidden error
func NewForbiddenError(message string) *AppError {
	if message == "" {
		message = "You do not have permission to perform this action"
	}
	return &AppError{
		Message:    message,
		Code:       ErrCodeForbidden,
		StatusCode: http.StatusForbidden,
	}
}

// ============================================================================
// Validation Errors (400)
// ============================================================================

// NewValidationError creates a 400 Validation error
func NewValidationError(message string) *AppError {
	if message == "" {
		message = "Invalid request data"
	}
	return &AppError{
		Message:    message,
		Code:       ErrCodeValidation,
		StatusCode: http.StatusBadRequest,
	}
}

// NewBadRequestError creates a 400 Bad Request error
func NewBadRequestError(message string) *AppError {
	if message == "" {
		message = "Invalid request"
	}
	return &AppError{
		Message:    message,
		Code:       ErrCodeInvalidInput,
		StatusCode: http.StatusBadRequest,
	}
}

// NewInvalidFormatError creates a 400 Invalid Format error
func NewInvalidFormatError(field, format string) *AppError {
	message := fmt.Sprintf("Invalid format for field '%s'", field)
	if format != "" {
		message = fmt.Sprintf("Invalid format for field '%s': expected %s", field, format)
	}
	return &AppError{
		Message:    message,
		Code:       ErrCodeInvalidFormat,
		StatusCode: http.StatusBadRequest,
		Details: map[string]interface{}{
			"field": field,
		},
	}
}

// ============================================================================
// Not Found Errors (404)
// ============================================================================

// NewNotFoundError creates a generic 404 Not Found error
func NewNotFoundError(resource string) *AppError {
	message := fmt.Sprintf("%s not found", resource)
	return &AppError{
		Message:    message,
		Code:       ErrCodeResourceNotFound,
		StatusCode: http.StatusNotFound,
	}
}

// NewPackageNotFoundError creates a 404 Package Not Found error
func NewPackageNotFoundError(id string) *AppError {
	return &AppError{
		Message:    fmt.Sprintf("Package with ID '%s' not found", id),
		Code:       ErrCodePackageNotFound,
		StatusCode: http.StatusNotFound,
		Details: map[string]interface{}{
			"packageId": id,
		},
	}
}

// NewThemeNotFoundError creates a 404 Theme Not Found error
func NewThemeNotFoundError(id string) *AppError {
	return &AppError{
		Message:    fmt.Sprintf("Theme with ID '%s' not found", id),
		Code:       ErrCodeThemeNotFound,
		StatusCode: http.StatusNotFound,
		Details: map[string]interface{}{
			"themeId": id,
		},
	}
}

// NewAddonNotFoundError creates a 404 Addon Not Found error
func NewAddonNotFoundError(id string) *AppError {
	return &AppError{
		Message:    fmt.Sprintf("Addon with ID '%s' not found", id),
		Code:       ErrCodeAddonNotFound,
		StatusCode: http.StatusNotFound,
		Details: map[string]interface{}{
			"addonId": id,
		},
	}
}

// NewBookingNotFoundError creates a 404 Booking Not Found error
func NewBookingNotFoundError(id string) *AppError {
	return &AppError{
		Message:    fmt.Sprintf("Booking with ID '%s' not found", id),
		Code:       ErrCodeBookingNotFound,
		StatusCode: http.StatusNotFound,
		Details: map[string]interface{}{
			"bookingId": id,
		},
	}
}

// NewUserNotFoundError creates a 404 User Not Found error
func NewUserNotFoundError(identifier string) *AppError {
	return &AppError{
		Message:    fmt.Sprintf("User '%s' not found", identifier),
		Code:       ErrCodeUserNotFound,
		StatusCode: http.StatusNotFound,
	}
}

// ============================================================================
// Conflict Errors (409)
// ============================================================================

// NewConflictError creates a 409 Conflict error
func NewConflictError(message string) *AppError {
	if message == "" {
		message = "Resource conflict"
	}
	return &AppError{
		Message:    message,
		Code:       ErrCodeConflict,
		StatusCode: http.StatusConflict,
	}
}

// NewSlotUnavailableError creates a 409 Slot Unavailable error
func NewSlotUnavailableError(date, time, themeID string) *AppError {
	return &AppError{
		Message:    "Selected time slot is no longer available",
		Code:       ErrCodeSlotUnavailable,
		StatusCode: http.StatusConflict,
		Details: map[string]interface{}{
			"slot": map[string]string{
				"date":    date,
				"time":    time,
				"themeId": themeID,
			},
		},
	}
}

// NewAlreadyExistsError creates a 409 Already Exists error
func NewAlreadyExistsError(resource string) *AppError {
	return &AppError{
		Message:    fmt.Sprintf("%s already exists", resource),
		Code:       ErrCodeAlreadyExists,
		StatusCode: http.StatusConflict,
	}
}

// ============================================================================
// Business Logic Errors (422)
// ============================================================================

// NewUnprocessableEntityError creates a 422 Unprocessable Entity error
func NewUnprocessableEntityError(message string) *AppError {
	if message == "" {
		message = "Request cannot be processed"
	}
	return &AppError{
		Message:    message,
		Code:       ErrCodeValidation,
		StatusCode: http.StatusUnprocessableEntity,
	}
}

// NewInvalidSlotCountError creates a 422 Invalid Slot Count error
func NewInvalidSlotCountError(expected, actual int) *AppError {
	return &AppError{
		Message:    fmt.Sprintf("Invalid number of slots: expected %d, got %d", expected, actual),
		Code:       ErrCodeInvalidSlotCount,
		StatusCode: http.StatusUnprocessableEntity,
		Details: map[string]interface{}{
			"expected": expected,
			"actual":   actual,
		},
	}
}

// ============================================================================
// File Upload Errors (400)
// ============================================================================

// NewFileTooLargeError creates a 400 File Too Large error
func NewFileTooLargeError(maxSize int64) *AppError {
	return &AppError{
		Message:    fmt.Sprintf("File size exceeds %dMB limit", maxSize/(1024*1024)),
		Code:       ErrCodeFileTooLarge,
		StatusCode: http.StatusBadRequest,
		Details: map[string]interface{}{
			"maxSizeBytes": maxSize,
		},
	}
}

// NewInvalidFileTypeError creates a 400 Invalid File Type error
func NewInvalidFileTypeError(allowedTypes []string) *AppError {
	return &AppError{
		Message:    "Only JPG, JPEG, and PNG files are allowed",
		Code:       ErrCodeInvalidFileType,
		StatusCode: http.StatusBadRequest,
		Details: map[string]interface{}{
			"allowedTypes": allowedTypes,
		},
	}
}

// NewUploadFailedError creates a 500 Upload Failed error
func NewUploadFailedError(err error) *AppError {
	return &AppError{
		Message:    "Failed to upload file. Please try again.",
		Code:       ErrCodeUploadFailed,
		StatusCode: http.StatusInternalServerError,
		Err:        err,
	}
}

// ============================================================================
// Server Errors (500)
// ============================================================================

// NewInternalError creates a 500 Internal Server Error
func NewInternalError(message string, err error) *AppError {
	if message == "" {
		message = "An unexpected error occurred"
	}
	return &AppError{
		Message:    message,
		Code:       ErrCodeInternal,
		StatusCode: http.StatusInternalServerError,
		Err:        err,
	}
}

// NewDatabaseError creates a 500 Database Error
func NewDatabaseError(operation string, err error) *AppError {
	return &AppError{
		Message:    fmt.Sprintf("Database error during %s", operation),
		Code:       ErrCodeDatabaseError,
		StatusCode: http.StatusInternalServerError,
		Err:        err,
	}
}

// NewEmailSendFailedError creates a 500 Email Send Failed error
// Note: This should be logged but may not always fail the request
func NewEmailSendFailedError(err error) *AppError {
	return &AppError{
		Message:    "Email notification could not be sent",
		Code:       ErrCodeEmailSendFailed,
		StatusCode: http.StatusInternalServerError,
		Err:        err,
	}
}
