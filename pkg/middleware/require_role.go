package middleware

import (
	"lqstudio-backend/pkg/errors"

	"github.com/labstack/echo/v4"
)

// RequireRole is a middleware that checks if the authenticated user has the required role
// This middleware should be used after RequireAuth middleware
func RequireRole(role string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get user role from context (set by RequireAuth middleware)
			userRole, ok := c.Get("user_role").(string)
			if !ok {
				return errors.NewUnauthorizedError("User role not found in context. Ensure RequireAuth middleware is applied first.")
			}

			// Check if user has the required role
			if userRole != role {
				return errors.NewForbiddenError("You do not have permission to perform this action")
			}

			return next(c)
		}
	}
}
