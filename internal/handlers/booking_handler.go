package handlers

import (
	"net/http"
	"strings"

	"lqstudio-backend/internal/config"
	"lqstudio-backend/internal/dto"
	"lqstudio-backend/internal/services"
	"lqstudio-backend/pkg/errors"
	"lqstudio-backend/pkg/upload"

	"github.com/labstack/echo/v4"
)

type BookingHandler struct {
	bookingService *services.BookingService
	uploadConfig   *config.UploadConfig
}

func NewBookingHandler(
	bookingService *services.BookingService,
	uploadConfig *config.UploadConfig,
) *BookingHandler {
	return &BookingHandler{
		bookingService: bookingService,
		uploadConfig:   uploadConfig,
	}
}

// CreateBooking godoc
// @Summary Create new booking
// @Description Create a new booking (customer-facing, no auth required).
// @Description
// @Description **Slot rules depend on the package type:**
// @Description
// @Description **1-slot / 2-slot packages** — send exactly 1 or 2 slots, each with a `themeId`:
// @Description ```json
// @Description {
// @Description   "packageId": "pkg-single",
// @Description   "slots": [
// @Description     { "date": "2026-02-24", "time": "2:20 PM", "themeId": "theme-A" }
// @Description   ],
// @Description   "addons": [],
// @Description   "customer": { "name": "Anas", "email": "anas@example.com", "phone": "0123456789" }
// @Description }
// @Description ```
// @Description
// @Description **Studio-level (3-slot / 60-min) packages** — send exactly 3 time slots **without** `themeId`.
// @Description The backend automatically books every active theme for those 3 slots.
// @Description Sending `themeId` is harmless but it will be ignored.
// @Description ```json
// @Description {
// @Description   "packageId": "pkg-studio",
// @Description   "slots": [
// @Description     { "date": "2026-02-24", "time": "2:20 PM" },
// @Description     { "date": "2026-02-24", "time": "2:40 PM" },
// @Description     { "date": "2026-02-24", "time": "3:00 PM" }
// @Description   ],
// @Description   "addons": [],
// @Description   "customer": { "name": "Anas", "email": "anas@example.com", "phone": "0123456789" }
// @Description }
// @Description ```
// @Tags bookings
// @Accept json
// @Produce json
// @Param request body dto.BookingRequest true "Booking data"
// @Success 201 {object} dto.ApiResponse{data=dto.BookingResponse}
// @Failure 400 {object} dto.ApiResponse "Invalid request, wrong slot count, or missing themeId for non-studio package"
// @Failure 404 {object} dto.ApiResponse "Package, theme, or addon not found"
// @Failure 409 {object} dto.ApiResponse "One or more requested slots are already booked"
// @Failure 500 {object} dto.ApiResponse
// @Router /api/bookings [post]
func (h *BookingHandler) CreateBooking(c echo.Context) error {
	var req dto.BookingRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	result, err := h.bookingService.CreateBooking(c.Request().Context(), &req)
	if err != nil {
		// Check if it's our custom AppError - return it directly so the error handler can process it
		if appErr, ok := err.(*errors.AppError); ok {
			return appErr
		}

		// Check if it's a conflict error (slot unavailable)
		if strings.Contains(err.Error(), "not available") || strings.Contains(err.Error(), "conflict") || strings.Contains(err.Error(), "already booked") {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}
		// Check if it's a validation error
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "not active") || strings.Contains(err.Error(), "invalid") {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return SendCreated(c, result, "Booking created successfully")
}

// GetBookingByID godoc
// @Summary Get booking by ID
// @Description Get a booking by ID (public - for customers to check their booking)
// @Tags bookings
// @Accept json
// @Produce json
// @Param id path string true "Booking ID"
// @Success 200 {object} dto.ApiResponse{data=dto.BookingResponse}
// @Failure 400 {object} dto.ApiResponse
// @Failure 404 {object} dto.ApiResponse
// @Failure 500 {object} dto.ApiResponse
// @Router /api/bookings/{id} [get]
func (h *BookingHandler) GetBookingByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid booking ID")
	}

	booking, err := h.bookingService.GetBookingByID(c.Request().Context(), id)
	if err != nil {
		// Check if it's our custom AppError - return it directly
		if appErr, ok := err.(*errors.AppError); ok {
			return appErr
		}

		// Legacy error handling
		if strings.Contains(err.Error(), "not found") {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return SendOK(c, booking, "Booking retrieved successfully")
}

// UploadPaymentScreenshot godoc
// @Summary Upload payment screenshot
// @Description Upload payment screenshot for a booking
// @Tags bookings
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Booking ID"
// @Param file formData file true "Payment screenshot image"
// @Success 200 {object} dto.ApiResponse{data=dto.BookingResponse}
// @Failure 400 {object} dto.ApiResponse "Invalid file or booking ID"
// @Failure 404 {object} dto.ApiResponse "Booking not found"
// @Failure 413 {object} dto.ApiResponse "File too large"
// @Failure 415 {object} dto.ApiResponse "Invalid file type"
// @Failure 500 {object} dto.ApiResponse
// @Router /api/bookings/{id}/payment-screenshot [post]
func (h *BookingHandler) UploadPaymentScreenshot(c echo.Context) error {
	// Get booking ID from URL parameter
	bookingID := c.Param("id")
	if bookingID == "" {
		appErr := errors.NewBadRequestError("Booking ID is required")
		return c.JSON(appErr.StatusCode, dto.NewErrorResponseWithCode(appErr.Message, appErr.Code))
	}

	// Get file from form data
	fileHeader, err := c.FormFile("file")
	if err != nil {
		appErr := errors.NewBadRequestError("File is required")
		return c.JSON(appErr.StatusCode, dto.NewErrorResponseWithCode(appErr.Message, appErr.Code))
	}

	// Validate file (size and MIME type)
	if err := upload.ValidateFile(fileHeader, h.uploadConfig.MaxFileSize, h.uploadConfig.AllowedTypes); err != nil {
		// Determine if it's a file size or type error
		if strings.Contains(err.Error(), "exceeds") {
			appErr := errors.NewFileTooLargeError(h.uploadConfig.MaxFileSize)
			return c.JSON(appErr.StatusCode, dto.NewErrorResponseWithCode(appErr.Message, appErr.Code))
		}
		appErr := errors.NewInvalidFileTypeError(h.uploadConfig.AllowedTypes)
		return c.JSON(appErr.StatusCode, dto.NewErrorResponseWithCode(appErr.Message, appErr.Code))
	}

	// Open file
	file, err := fileHeader.Open()
	if err != nil {
		appErr := errors.NewUploadFailedError(err)
		return c.JSON(appErr.StatusCode, dto.NewErrorResponseWithCode(appErr.Message, appErr.Code))
	}
	defer file.Close()

	// Generate unique filename
	uniqueFilename := upload.GenerateUniqueFilename(fileHeader.Filename, "payment")

	// Save file and get URL path
	urlPath, err := upload.SaveFile(file, uniqueFilename, h.uploadConfig.StoragePath, "payment-screenshots")
	if err != nil {
		appErr := errors.NewUploadFailedError(err)
		return c.JSON(appErr.StatusCode, dto.NewErrorResponseWithCode(appErr.Message, appErr.Code))
	}

	// Update booking with payment screenshot URL
	booking, err := h.bookingService.UpdatePaymentScreenshot(c.Request().Context(), bookingID, urlPath)
	if err != nil {
		// Check if booking not found
		if appErr, ok := err.(*errors.AppError); ok {
			return c.JSON(appErr.StatusCode, dto.NewErrorResponseWithCode(appErr.Message, appErr.Code))
		}
		// Generic error
		appErr := errors.NewInternalError("Failed to update booking", err)
		return c.JSON(appErr.StatusCode, dto.NewErrorResponseWithCode(appErr.Message, appErr.Code))
	}

	return SendOK(c, booking, "Payment screenshot uploaded successfully")
}
