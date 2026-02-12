package handlers

import (
	"lqstudio-backend/internal/database"
	"net/http"

	"github.com/labstack/echo/v4"
)

type HealthHandler struct {
	db *database.Connection
}

func NewHealthHandler(db *database.Connection) *HealthHandler {
	return &HealthHandler{db: db}
}

// HealthCheck godoc
// @Summary Health check
// @Description Check if the API and database are healthy
// @Tags health
// @Produce json
// @Success 200 {object} dto.ApiResponse{data=map[string]string}
// @Failure 503 {object} dto.ApiResponse
// @Router /health [get]
func (h *HealthHandler) HealthCheck(c echo.Context) error {
	// Check database health
	if err := h.db.Health(c.Request().Context()); err != nil {
		return SendSuccess(c, http.StatusServiceUnavailable, map[string]string{
			"status":   "unhealthy",
			"database": "down",
		}, "Service is unhealthy")
	}

	return SendOK(c, map[string]string{
		"status":   "healthy",
		"database": "up",
	}, "Service is healthy")
}
