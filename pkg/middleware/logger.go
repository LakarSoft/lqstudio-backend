package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func Logger(logger *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			// Process request
			err := next(c)

			// Skip logging for metrics endpoint to reduce log clutter
			if c.Request().RequestURI != "/metrics" {
				// Log request
				logger.Info("HTTP Request",
					zap.String("method", c.Request().Method),
					zap.String("uri", c.Request().RequestURI),
					zap.Int("status", c.Response().Status),
					zap.Duration("latency", time.Since(start)),
					zap.String("remote_ip", c.RealIP()),
					zap.String("user_agent", c.Request().UserAgent()),
				)
			}

			return err
		}
	}
}
