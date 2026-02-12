package middleware

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func CORS(allowedOrigins []string) echo.MiddlewareFunc {
	// Debug: Print what origins the middleware is receiving
	fmt.Printf("[CORS Middleware] Allowed Origins (%d): %v\n", len(allowedOrigins), allowedOrigins)

	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: allowedOrigins,
		AllowMethods: []string{
			echo.GET,
			echo.POST,
			echo.PUT,
			echo.PATCH,
			echo.DELETE,
			echo.OPTIONS,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
		},
		AllowCredentials: true,
		MaxAge:           300,
	})
}
