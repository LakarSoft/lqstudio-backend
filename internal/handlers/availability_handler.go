package handlers

import (
	"net/http"
	"strings"

	"lqstudio-backend/internal/dto"
	"lqstudio-backend/internal/services"

	"github.com/labstack/echo/v4"
)

type AvailabilityHandler struct {
	bookingService *services.BookingService
}

func NewAvailabilityHandler(bookingService *services.BookingService) *AvailabilityHandler {
	return &AvailabilityHandler{
		bookingService: bookingService,
	}
}

// GetAvailability godoc
// @Summary Get available time slots
// @Description Check available time slots for a specific theme or all themes on a specific date. Use themeId="all" to check availability across all active themes.
// @Tags availability
// @Accept json
// @Produce json
// @Param themeId path string true "Theme ID or 'all' for all themes"
// @Param date query string true "Date in YYYY-MM-DD format"
// @Param packageId query string false "Package ID (optional)"
// @Success 200 {object} dto.ApiResponse{data=dto.AvailabilityResponse} "Returns availability with themeId in slots for specific theme, or without themeId for all themes"
// @Failure 400 {object} dto.ApiResponse "Invalid parameters"
// @Failure 404 {object} dto.ApiResponse "Theme not found"
// @Failure 500 {object} dto.ApiResponse
// @Router /api/themes/{themeId}/available-times [get]
func (h *AvailabilityHandler) GetAvailability(c echo.Context) error {
	// Get themeId from path parameter (as string)
	themeID := c.Param("themeId")
	if themeID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "theme ID is required")
	}

	// Get date from query parameter
	date := c.QueryParam("date")
	if date == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "date query parameter is required")
	}

	// Package ID is now optional
	packageID := c.QueryParam("packageId")

	// Create request
	req := dto.AvailabilityRequest{
		Date:      date,
		PackageID: packageID,
		ThemeID:   themeID,
	}

	// Get availability
	result, err := h.bookingService.GetAvailability(c.Request().Context(), &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return SendOK(c, result, "Availability retrieved successfully")
}
