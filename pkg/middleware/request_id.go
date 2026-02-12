package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// RequestIDContextKey is the key used to store request ID in context
const RequestIDContextKey = "request_id"

// RequestIDHeader is the header name for request ID
const RequestIDHeader = "X-Request-ID"

// RequestID middleware generates a unique request ID for each request
// and stores it in the context for use in responses and logging
func RequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Check if request ID already exists in header (from client or load balancer)
			requestID := c.Request().Header.Get(RequestIDHeader)

			// Generate new UUID if not present
			if requestID == "" {
				requestID = uuid.New().String()
			}

			// Set in response header
			c.Response().Header().Set(RequestIDHeader, requestID)

			// Store in context for handlers to use
			c.Set(RequestIDContextKey, requestID)

			return next(c)
		}
	}
}

// GetRequestID retrieves the request ID from the context
func GetRequestID(c echo.Context) string {
	if requestID, ok := c.Get(RequestIDContextKey).(string); ok {
		return requestID
	}
	return ""
}
