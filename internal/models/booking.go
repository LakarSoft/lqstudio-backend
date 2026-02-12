package models

import (
	"time"

	"github.com/shopspring/decimal"
)

// Booking represents a studio booking
type Booking struct {
	ID                     string           `json:"id"`
	PackageID              string           `json:"packageId"`
	BookingDate            time.Time        `json:"bookingDate"`
	StartTime              string           `json:"startTime"`
	Status                 BookingStatus    `json:"status"`
	CustomerName           string           `json:"customerName"`
	CustomerEmail          string           `json:"customerEmail"`
	CustomerPhone          string           `json:"customerPhone"`
	CustomerNotes          string           `json:"customerNotes"`
	PackageAmount          decimal.Decimal  `json:"packageAmount"`
	AddOnsAmount           decimal.Decimal  `json:"addOnsAmount"`
	TotalAmount            decimal.Decimal  `json:"totalAmount"`
	PaymentScreenshotURL   *string          `json:"paymentScreenshotUrl,omitempty"`
	PaymentConfirmedAt     *time.Time       `json:"paymentConfirmedAt,omitempty"`
	AdminNotes             string           `json:"adminNotes"`
	CreatedAt              time.Time        `json:"createdAt"`
	UpdatedAt              time.Time        `json:"updatedAt"`

	// Related entities (loaded separately)
	Slots                  []BookingSlot    `json:"slots,omitempty"`
	Addons                 []BookingAddon   `json:"addons,omitempty"`
}

// BookingStatus represents booking status
type BookingStatus string

const (
	BookingStatusPending  BookingStatus = "PENDING"
	BookingStatusApproved BookingStatus = "APPROVED"
	BookingStatusRejected BookingStatus = "REJECTED"
)

// BookingSlot represents a single 20-minute theme slot
type BookingSlot struct {
	ID        int32     `json:"id"`
	BookingID string    `json:"bookingId"`
	ThemeID   string    `json:"themeId"`
	Date      time.Time `json:"date"`
	Time      string    `json:"time"` // HH:MM format
}

// BookingAddon represents an add-on line item
type BookingAddon struct {
	ID        int32 `json:"id"`
	BookingID string `json:"bookingId"`
	AddonID   string `json:"addonId"`
	Quantity  int   `json:"quantity"`
}

// AvailableSlot represents an available time slot
type AvailableSlot struct {
	StartTime string `json:"startTime"`
	Available bool   `json:"available"`
}
