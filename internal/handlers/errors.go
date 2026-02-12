package handlers

import (
	"database/sql"
	stderr "errors"
	"strings"

	"lqstudio-backend/pkg/errors"
)

// HandleServiceError converts service layer errors to appropriate AppErrors
// This is a helper function for handlers to use when calling services
func HandleServiceError(err error, resourceType string) error {
	if err == nil {
		return nil
	}

	// Check if it's already an AppError
	var appErr *errors.AppError
	if stderr.As(err, &appErr) {
		return appErr
	}

	// Handle database errors
	if stderr.Is(err, sql.ErrNoRows) {
		switch resourceType {
		case "package":
			return errors.NewPackageNotFoundError("")
		case "theme":
			return errors.NewThemeNotFoundError("")
		case "addon":
			return errors.NewAddonNotFoundError("")
		case "booking":
			return errors.NewBookingNotFoundError("")
		case "user":
			return errors.NewUserNotFoundError("")
		default:
			return errors.NewNotFoundError(resourceType)
		}
	}

	// Check error message patterns (for legacy errors)
	errMsg := err.Error()

	// Not found errors
	if strings.Contains(errMsg, "not found") {
		switch resourceType {
		case "package":
			return errors.NewPackageNotFoundError("")
		case "theme":
			return errors.NewThemeNotFoundError("")
		case "addon":
			return errors.NewAddonNotFoundError("")
		case "booking":
			return errors.NewBookingNotFoundError("")
		case "user":
			return errors.NewUserNotFoundError("")
		default:
			return errors.NewNotFoundError(resourceType)
		}
	}

	// Conflict/unavailable errors
	if strings.Contains(errMsg, "not available") ||
		strings.Contains(errMsg, "already booked") ||
		strings.Contains(errMsg, "conflict") {
		return errors.NewConflictError(errMsg)
	}

	// Validation errors
	if strings.Contains(errMsg, "invalid") ||
		strings.Contains(errMsg, "validation") {
		return errors.NewValidationError(errMsg)
	}

	// Already exists errors
	if strings.Contains(errMsg, "already exists") ||
		strings.Contains(errMsg, "duplicate") {
		return errors.NewAlreadyExistsError(resourceType)
	}

	// Default: return as internal error
	return errors.NewInternalError("An unexpected error occurred", err)
}
