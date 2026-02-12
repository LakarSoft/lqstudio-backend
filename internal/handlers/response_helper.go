package handlers

import (
	"lqstudio-backend/internal/dto"
	"lqstudio-backend/pkg/middleware"
	"net/http"

	"github.com/labstack/echo/v4"
)

// SendSuccess sends a standardized success response
func SendSuccess(c echo.Context, statusCode int, data interface{}, message string) error {
	response := dto.NewApiSuccessResponse(data, message)

	// Add request ID if available
	if requestID := middleware.GetRequestID(c); requestID != "" {
		response.WithRequestID(requestID)
	}

	return c.JSON(statusCode, response)
}

// SendCreated sends a 201 Created success response
func SendCreated(c echo.Context, data interface{}, message string) error {
	if message == "" {
		message = "Resource created successfully"
	}
	return SendSuccess(c, http.StatusCreated, data, message)
}

// SendOK sends a 200 OK success response
func SendOK(c echo.Context, data interface{}, message string) error {
	if message == "" {
		message = "Operation completed successfully"
	}
	return SendSuccess(c, http.StatusOK, data, message)
}

// SendDeleted sends a success response for delete operations
func SendDeleted(c echo.Context, message string) error {
	if message == "" {
		message = "Resource deleted successfully"
	}
	return SendSuccess(c, http.StatusOK, nil, message)
}
