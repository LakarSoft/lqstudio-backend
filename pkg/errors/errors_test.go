package errors

import (
	"errors"
	"net/http"
	"testing"
)

func TestNewPackageNotFoundError(t *testing.T) {
	packageID := "pkg-123"
	err := NewPackageNotFoundError(packageID)

	if err.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, err.StatusCode)
	}

	if err.Code != ErrCodePackageNotFound {
		t.Errorf("Expected code %s, got %s", ErrCodePackageNotFound, err.Code)
	}

	if err.Details["packageId"] != packageID {
		t.Errorf("Expected packageId in details to be %s, got %v", packageID, err.Details["packageId"])
	}
}

func TestNewSlotUnavailableError(t *testing.T) {
	date := "2025-03-15"
	time := "10:00"
	themeID := "theme-minimalist"

	err := NewSlotUnavailableError(date, time, themeID)

	if err.StatusCode != http.StatusConflict {
		t.Errorf("Expected status code %d, got %d", http.StatusConflict, err.StatusCode)
	}

	if err.Code != ErrCodeSlotUnavailable {
		t.Errorf("Expected code %s, got %s", ErrCodeSlotUnavailable, err.Code)
	}

	slot, ok := err.Details["slot"].(map[string]string)
	if !ok {
		t.Fatal("Expected slot details to be a map[string]string")
	}

	if slot["date"] != date {
		t.Errorf("Expected date to be %s, got %s", date, slot["date"])
	}

	if slot["time"] != time {
		t.Errorf("Expected time to be %s, got %s", time, slot["time"])
	}

	if slot["themeId"] != themeID {
		t.Errorf("Expected themeId to be %s, got %s", themeID, slot["themeId"])
	}
}

func TestNewInvalidSlotCountError(t *testing.T) {
	expected := 3
	actual := 1

	err := NewInvalidSlotCountError(expected, actual)

	if err.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("Expected status code %d, got %d", http.StatusUnprocessableEntity, err.StatusCode)
	}

	if err.Code != ErrCodeInvalidSlotCount {
		t.Errorf("Expected code %s, got %s", ErrCodeInvalidSlotCount, err.Code)
	}

	if err.Details["expected"] != expected {
		t.Errorf("Expected 'expected' to be %d, got %v", expected, err.Details["expected"])
	}

	if err.Details["actual"] != actual {
		t.Errorf("Expected 'actual' to be %d, got %v", actual, err.Details["actual"])
	}
}

func TestWithDetails(t *testing.T) {
	err := NewValidationError("Test error")
	details := map[string]interface{}{
		"field": "email",
		"issue": "invalid format",
	}

	err = err.WithDetails(details)

	if err.Details["field"] != "email" {
		t.Errorf("Expected field to be 'email', got %v", err.Details["field"])
	}

	if err.Details["issue"] != "invalid format" {
		t.Errorf("Expected issue to be 'invalid format', got %v", err.Details["issue"])
	}
}

func TestWithError(t *testing.T) {
	originalErr := errors.New("database connection failed")
	appErr := NewDatabaseError("create booking", originalErr)

	if appErr.Err == nil {
		t.Fatal("Expected wrapped error to be set")
	}

	if !errors.Is(appErr, originalErr) {
		t.Error("Expected error chain to include original error")
	}
}

func TestErrorInterface(t *testing.T) {
	err := NewBookingNotFoundError("bk-123")

	// Test that AppError implements error interface
	var _ error = err

	expectedMessage := "Booking with ID 'bk-123' not found"
	if err.Error() != expectedMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedMessage, err.Error())
	}
}

func TestUnwrap(t *testing.T) {
	originalErr := errors.New("original error")
	appErr := NewInternalError("Something went wrong", originalErr)

	unwrapped := appErr.Unwrap()
	if unwrapped != originalErr {
		t.Error("Expected Unwrap to return the original error")
	}
}
