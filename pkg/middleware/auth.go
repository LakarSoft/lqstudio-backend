package middleware

import (
	"strings"

	"lqstudio-backend/pkg/auth"
	"lqstudio-backend/pkg/errors"

	"github.com/labstack/echo/v4"
)

// RequireAuth is a middleware that validates JWT tokens from Authorization header
func RequireAuth(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return errors.NewUnauthorizedError("Missing authorization header")
			}

			// Check if it starts with "Bearer "
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return errors.NewUnauthorizedError("Invalid authorization header format. Expected: Bearer <token>")
			}

			tokenString := parts[1]
			if tokenString == "" {
				return errors.NewUnauthorizedError("Missing token")
			}

			// Validate token
			claims, err := auth.ValidateToken(tokenString, secret)
			if err != nil {
				return errors.NewInvalidTokenError("Token is invalid or expired")
			}

			// Store claims in context for use in handlers
			c.Set("user_id", claims.UserID)
			c.Set("user_email", claims.Email)
			c.Set("user_role", claims.Role)
			c.Set("claims", claims)

			return next(c)
		}
	}
}
