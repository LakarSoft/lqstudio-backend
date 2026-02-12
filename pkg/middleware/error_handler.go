package middleware

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"lqstudio-backend/internal/dto"
	appErrors "lqstudio-backend/pkg/errors"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

// ErrorHandler is the custom error handler for Echo framework
// It converts errors to standardized API error responses
func ErrorHandler(err error, c echo.Context) {
	// Prevent double sending response
	if c.Response().Committed {
		return
	}

	var (
		code         = http.StatusInternalServerError
		errorMessage string
		errorCode    string
		errorDetails map[string]interface{}
	)

	// Check if it's our custom AppError
	var appErr *appErrors.AppError
	if errors.As(err, &appErr) {
		code = appErr.StatusCode
		errorMessage = appErr.Message
		errorCode = appErr.Code
		errorDetails = appErr.Details
		// Log the underlying error if present
		if appErr.Err != nil {
			c.Logger().Errorf("AppError: %s, underlying: %v", appErr.Message, appErr.Err)
		}
	} else if echoErr, ok := err.(*echo.HTTPError); ok {
		// Handle Echo HTTP errors
		code = echoErr.Code
		errorMessage = extractMessage(echoErr.Message)
		errorCode = mapHTTPCodeToErrorCode(code)

		// Log internal server errors
		if code == http.StatusInternalServerError {
			c.Logger().Errorf("Echo HTTPError: %v", echoErr.Internal)
		}
	} else {
		// Handle standard Go errors and database errors
		code, errorMessage, errorCode, errorDetails = handleStandardError(err, c)
	}

	// Don't expose internal error details in production
	if code == http.StatusInternalServerError {
		// Keep the error code but sanitize the message
		if errorCode == "" {
			errorCode = appErrors.ErrCodeInternal
		}
		errorMessage = "An unexpected error occurred"
		errorDetails = nil

		// Log the actual error
		c.Logger().Errorf("Internal error: %v", err)
	}

	// Create standardized API error response
	response := dto.NewApiErrorResponse(errorMessage, errorCode, errorDetails)

	// Add request ID if available
	if requestID := GetRequestID(c); requestID != "" {
		response.WithRequestID(requestID)
	}

	// Send error response
	if err := c.JSON(code, response); err != nil {
		c.Logger().Errorf("Failed to send error response: %v", err)
	}
}

// handleStandardError handles standard Go errors and database errors
func handleStandardError(err error, c echo.Context) (int, string, string, map[string]interface{}) {
	// Check for database errors
	if errors.Is(err, sql.ErrNoRows) {
		return http.StatusNotFound, "Resource not found", appErrors.ErrCodeResourceNotFound, nil
	}

	// Check for PostgreSQL errors
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return handlePostgresError(pqErr)
	}

	// Check error message patterns (for legacy errors)
	errMsg := err.Error()

	// Not found errors
	if strings.Contains(errMsg, "not found") {
		code := appErrors.ErrCodeResourceNotFound
		if strings.Contains(errMsg, "package") {
			code = appErrors.ErrCodePackageNotFound
		} else if strings.Contains(errMsg, "theme") {
			code = appErrors.ErrCodeThemeNotFound
		} else if strings.Contains(errMsg, "addon") {
			code = appErrors.ErrCodeAddonNotFound
		} else if strings.Contains(errMsg, "booking") {
			code = appErrors.ErrCodeBookingNotFound
		}
		return http.StatusNotFound, errMsg, code, nil
	}

	// Conflict/unavailable errors
	if strings.Contains(errMsg, "not available") ||
		strings.Contains(errMsg, "already booked") ||
		strings.Contains(errMsg, "conflict") {
		return http.StatusConflict, errMsg, appErrors.ErrCodeSlotUnavailable, nil
	}

	// Validation errors
	if strings.Contains(errMsg, "invalid") ||
		strings.Contains(errMsg, "required") ||
		strings.Contains(errMsg, "validation") {
		return http.StatusBadRequest, errMsg, appErrors.ErrCodeValidation, nil
	}

	// Default internal server error
	c.Logger().Errorf("Unhandled error: %v", err)
	return http.StatusInternalServerError, "An unexpected error occurred", appErrors.ErrCodeInternal, nil
}

// handlePostgresError converts PostgreSQL errors to appropriate HTTP errors
func handlePostgresError(pqErr *pq.Error) (int, string, string, map[string]interface{}) {
	switch pqErr.Code {
	case "23505": // unique_violation
		return http.StatusConflict, "Resource already exists", appErrors.ErrCodeAlreadyExists, nil
	case "23503": // foreign_key_violation
		return http.StatusBadRequest, "Invalid reference to related resource", appErrors.ErrCodeValidation, nil
	case "23502": // not_null_violation
		return http.StatusBadRequest, "Required field is missing", appErrors.ErrCodeRequiredField, nil
	case "22P02": // invalid_text_representation
		return http.StatusBadRequest, "Invalid data format", appErrors.ErrCodeInvalidFormat, nil
	default:
		return http.StatusInternalServerError, "Database error occurred", appErrors.ErrCodeDatabaseError, nil
	}
}

// mapHTTPCodeToErrorCode maps HTTP status codes to error codes
func mapHTTPCodeToErrorCode(code int) string {
	switch code {
	case http.StatusBadRequest:
		return appErrors.ErrCodeInvalidInput
	case http.StatusUnauthorized:
		return appErrors.ErrCodeUnauthorized
	case http.StatusForbidden:
		return appErrors.ErrCodeForbidden
	case http.StatusNotFound:
		return appErrors.ErrCodeResourceNotFound
	case http.StatusConflict:
		return appErrors.ErrCodeConflict
	case http.StatusUnprocessableEntity:
		return appErrors.ErrCodeValidation
	case http.StatusInternalServerError:
		return appErrors.ErrCodeInternal
	default:
		return appErrors.ErrCodeInternal
	}
}

// extractMessage extracts string message from Echo error message (can be string or error)
func extractMessage(message interface{}) string {
	switch v := message.(type) {
	case string:
		return v
	case error:
		return v.Error()
	default:
		return "An error occurred"
	}
}
