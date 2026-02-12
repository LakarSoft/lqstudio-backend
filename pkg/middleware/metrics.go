package middleware

import (
	"lqstudio-backend/pkg/metrics"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

// Metrics is a middleware that automatically instruments HTTP requests with Prometheus metrics
func Metrics() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Track in-flight requests
			metrics.HTTPRequestsInFlight.Inc()
			defer metrics.HTTPRequestsInFlight.Dec()

			// Start timer
			start := time.Now()

			// Process request
			err := next(c)

			// Calculate duration
			duration := time.Since(start).Seconds()

			// Get response status
			status := c.Response().Status
			if err != nil {
				// If there's an error and status hasn't been set, use 500
				if status == 0 {
					status = 500
				}
			}

			// Normalize route path to prevent cardinality explosion
			// Echo provides the matched route path, which is already parameterized
			path := c.Path()
			if path == "" {
				path = c.Request().URL.Path
			}

			// Record metrics
			method := c.Request().Method
			statusStr := strconv.Itoa(status)

			metrics.HTTPRequestsTotal.WithLabelValues(method, path, statusStr).Inc()
			metrics.HTTPRequestDuration.WithLabelValues(method, path).Observe(duration)

			return err
		}
	}
}
