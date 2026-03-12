package dto

import (
	"time"

	"lqstudio-backend/internal/models"

	"github.com/shopspring/decimal"
)

// BookingRequest represents a booking creation request
// Matches the frontend booking flow structure with nested slots, addons, and customer info
type BookingRequest struct {
	PackageID string         `json:"packageId" validate:"required"`
	Slots     []SlotRequest  `json:"slots" validate:"required,min=1,max=3,dive"`
	Addons    []AddonRequest `json:"addons,omitempty,dive"`
	Customer  CustomerInfo   `json:"customer" validate:"required"`
}

// UpdateBookingRequest represents a booking update request (admin only).
// Payload is identical to BookingRequest minus packageId — the package cannot change.
// Slot rules are the same: 1/2-slot packages require themeId per slot;
// 3-slot (studio-level) packages omit themeId and the backend auto-assigns all themes.
type UpdateBookingRequest struct {
	Slots    []SlotRequest  `json:"slots" validate:"required,min=1,max=3,dive"`
	Addons   []AddonRequest `json:"addons,omitempty"`
	Customer CustomerInfo   `json:"customer" validate:"required"`
}

// SlotRequest represents a single booking slot
// Matches the frontend BookingSlot structure with camelCase JSON tags
type SlotRequest struct {
	Date    string `json:"date" validate:"required"` // ISO date string (YYYY-MM-DD)
	Time    string `json:"time" validate:"required"` // Time in HH:mm format (e.g., "10:00")
	ThemeID string `json:"themeId,omitempty"`        // Required for 1/2-slot packages; omitted for studio-level (3-slot) packages — backend auto-assigns all themes
}

// AddonRequest represents an add-on selection in a booking
// Matches the frontend SelectedAddon structure with camelCase JSON tags
type AddonRequest struct {
	AddonID  string `json:"addonId" validate:"required"`
	Quantity int    `json:"quantity" validate:"required,min=1"`
}

// CustomerInfo represents customer information
// Matches the frontend CustomerInfo structure with camelCase JSON tags
type CustomerInfo struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Phone string `json:"phone" validate:"required"`
	Notes string `json:"notes,omitempty"`
}

// BookingResponse represents a booking in API responses
// Matches the frontend Booking entity structure with camelCase JSON tags
type BookingResponse struct {
	ID                   string              `json:"id"`
	PackageID            string              `json:"packageId"`
	Slots                []SlotResponse      `json:"slots"`
	Addons               []AddonItemResponse `json:"addons,omitempty"`
	Customer             CustomerInfo        `json:"customer"`
	Status               string              `json:"status"` // PENDING | APPROVED | REJECTED | COMPLETED
	TotalPrice           float64             `json:"totalPrice"`
	PaymentScreenshotURL string              `json:"paymentScreenshotUrl,omitempty"`
	AdminNotes           string              `json:"adminNotes,omitempty"` // Admin-only notes for this booking
	CreatedAt            string              `json:"createdAt"`            // ISO 8601 datetime
	UpdatedAt            string              `json:"updatedAt"`            // ISO 8601 datetime
}

// SlotResponse represents a booked slot in responses
// Matches the frontend BookingSlot structure
type SlotResponse struct {
	Date    string `json:"date"`    // ISO date string (YYYY-MM-DD)
	Time    string `json:"time"`    // Time in HH:mm format
	ThemeID string `json:"themeId"` // Foreign key to Theme
}

// AddonItemResponse represents a booked add-on in responses
// Matches the frontend SelectedAddon structure
type AddonItemResponse struct {
	AddonID  string `json:"addonId"`
	Quantity int    `json:"quantity"`
}

// AvailabilityRequest for checking available time slots
// Matches the frontend expectations with camelCase JSON tags
type AvailabilityRequest struct {
	Date      string `json:"date" validate:"required"`    // YYYY-MM-DD
	PackageID string `json:"packageId,omitempty"`         // Optional: for package-specific availability
	ThemeID   string `json:"themeId" validate:"required"` // Theme to check availability for, or "all" for all themes
}

// AvailabilityResponse with available time slots
// Matches the frontend expectations with camelCase JSON tags
type AvailabilityResponse struct {
	Date  string              `json:"date"` // YYYY-MM-DD
	Slots []AvailableSlotInfo `json:"slots"`
}

// AvailableSlotInfo represents a single available time slot
// Matches the frontend structure with time, available status, and optional themeId
type AvailableSlotInfo struct {
	Time      string `json:"time"`              // HH:mm format
	Available bool   `json:"available"`         // true if slot is available
	ThemeID   string `json:"themeId,omitempty"` // Theme ID for this slot (omitted when checking all themes)
}

// UpdateBookingStatusRequest for admin status updates
// Matches the frontend expectations with camelCase JSON tags
type UpdateBookingStatusRequest struct {
	Status     string `json:"status" validate:"required,oneof=PENDING APPROVED REJECTED COMPLETED"`
	AdminNotes string `json:"adminNotes,omitempty"`
}

// UpdateBookingStatusResponse extends BookingResponse with email notification status
// Used to inform admin if customer notification email was sent successfully
type UpdateBookingStatusResponse struct {
	Booking               *BookingResponse `json:"booking"`
	EmailNotificationSent bool             `json:"emailNotificationSent"`
	EmailError            string           `json:"emailError,omitempty"`
}

// UpdateAdminNotesRequest for admin-only notes updates on a booking
// Allows the admin to set or clear the internal notes without touching the booking status
type UpdateAdminNotesRequest struct {
	AdminNotes string `json:"adminNotes"`
}

// BookingFilters for admin booking list with filtering, sorting, and pagination
type BookingFilters struct {
	Status    string `query:"status"`    // Filter by status (PENDING, APPROVED, REJECTED, COMPLETED)
	Email     string `query:"email"`     // Filter by customer email (partial match)
	PackageID string `query:"packageId"` // Filter by package ID
	ThemeID   string `query:"themeId"`   // Filter by theme (joins booking_slots)
	SlotDate  string `query:"slotDate"`  // Filter by slot date (YYYY-MM-DD)
	DateFrom  string `query:"dateFrom"`  // Filter by slot date from (YYYY-MM-DD)
	DateTo    string `query:"dateTo"`    // Filter by slot date to (YYYY-MM-DD)
	Search    string `query:"search"`    // Search customer name, email, or phone
	SortBy    string `query:"sortBy"`    // Sort field: createdAt, updatedAt, totalPrice, status (default: createdAt)
	Order     string `query:"order"`     // Sort order: asc, desc (default: desc)
	Page      int    `query:"page"`      // Page number (default: 1)
	Limit     int    `query:"limit"`     // Items per page (default: 20)
}

// Defaults sets default values for unset filter fields
func (f *BookingFilters) Defaults() {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Limit <= 0 {
		f.Limit = 20
	}
	if f.Limit > 100 {
		f.Limit = 100
	}
	if f.SortBy == "" {
		f.SortBy = "createdAt"
	}
	if f.Order == "" {
		f.Order = "desc"
	}
}

// Offset calculates the SQL OFFSET from page and limit
func (f *BookingFilters) Offset() int {
	return (f.Page - 1) * f.Limit
}

// PaginatedBookingsResponse wraps booking list with pagination metadata
type PaginatedBookingsResponse struct {
	Data       []*BookingResponse `json:"data"`
	Pagination PaginationInfo     `json:"pagination"`
}

// PaginationInfo provides pagination metadata
type PaginationInfo struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
}

// ToBookingResponse converts a domain model Booking to a BookingResponse DTO
// Handles decimal.Decimal to float64 conversion and time formatting
func ToBookingResponse(booking *models.Booking) *BookingResponse {
	if booking == nil {
		return nil
	}

	// Convert decimal.Decimal to float64
	totalPrice, _ := booking.TotalAmount.Float64()

	// Convert slots
	slots := make([]SlotResponse, len(booking.Slots))
	for i, slot := range booking.Slots {
		slots[i] = SlotResponse{
			Date:    slot.Date.Format("2006-01-02"), // YYYY-MM-DD format
			Time:    slot.Time,
			ThemeID: slot.ThemeID,
		}
	}

	// Convert addons
	addons := make([]AddonItemResponse, len(booking.Addons))
	for i, addon := range booking.Addons {
		addons[i] = AddonItemResponse{
			AddonID:  addon.AddonID,
			Quantity: addon.Quantity,
		}
	}

	// Build customer info
	customer := CustomerInfo{
		Name:  booking.CustomerName,
		Email: booking.CustomerEmail,
		Phone: booking.CustomerPhone,
		Notes: booking.CustomerNotes,
	}

	// Get payment screenshot URL
	paymentURL := ""
	if booking.PaymentScreenshotURL != nil {
		paymentURL = *booking.PaymentScreenshotURL
	}

	return &BookingResponse{
		ID:                   booking.ID,
		PackageID:            booking.PackageID,
		Slots:                slots,
		Addons:               addons,
		Customer:             customer,
		Status:               string(booking.Status),
		TotalPrice:           totalPrice,
		PaymentScreenshotURL: paymentURL,
		AdminNotes:           booking.AdminNotes,
		CreatedAt:            booking.CreatedAt.Format(time.RFC3339),
		UpdatedAt:            booking.UpdatedAt.Format(time.RFC3339),
	}
}

// ToBookingsResponse converts a slice of domain model Bookings to a slice of BookingResponse DTOs
func ToBookingsResponse(bookings []*models.Booking) []*BookingResponse {
	if bookings == nil {
		return []*BookingResponse{}
	}

	responses := make([]*BookingResponse, len(bookings))
	for i, booking := range bookings {
		responses[i] = ToBookingResponse(booking)
	}
	return responses
}

// ToBookingModel converts BookingRequest to domain model Booking
// This creates the basic booking structure; slots and addons need to be set separately
func (r *BookingRequest) ToBookingModel(packageAmount, addonsAmount decimal.Decimal) *models.Booking {
	totalAmount := packageAmount.Add(addonsAmount)

	// Convert slots
	slots := make([]models.BookingSlot, len(r.Slots))
	for i, slot := range r.Slots {
		// Parse date string to time.Time
		date, _ := time.Parse("2006-01-02", slot.Date)
		slots[i] = models.BookingSlot{
			ThemeID: slot.ThemeID,
			Date:    date,
			Time:    slot.Time,
		}
	}

	// Convert addons
	addons := make([]models.BookingAddon, len(r.Addons))
	for i, addon := range r.Addons {
		addons[i] = models.BookingAddon{
			AddonID:  addon.AddonID,
			Quantity: addon.Quantity,
		}
	}

	return &models.Booking{
		CustomerName:  r.Customer.Name,
		CustomerEmail: r.Customer.Email,
		CustomerPhone: r.Customer.Phone,
		CustomerNotes: r.Customer.Notes,
		PackageID:     r.PackageID,
		Status:        models.BookingStatusPending,
		PackageAmount: packageAmount,
		AddOnsAmount:  addonsAmount,
		TotalAmount:   totalAmount,
		Slots:         slots,
		Addons:        addons,
	}
}

// ToAvailableSlotInfos converts models.AvailableSlot to AvailableSlotInfo DTOs
func ToAvailableSlotInfos(slots []models.AvailableSlot, themeID string) []AvailableSlotInfo {
	if slots == nil {
		return []AvailableSlotInfo{}
	}

	infos := make([]AvailableSlotInfo, len(slots))
	for i, slot := range slots {
		infos[i] = AvailableSlotInfo{
			Time:      slot.StartTime,
			Available: slot.Available,
			ThemeID:   themeID,
		}
	}
	return infos
}
