package models

import (
	"time"

	"github.com/shopspring/decimal"
)

// Package represents a booking package
type Package struct {
	ID              string          `json:"id"`
	Name            string          `json:"name"`
	Description     string          `json:"description"`
	Offers          []string        `json:"offers"`
	DurationMinutes int32           `json:"durationMinutes"`
	Price           decimal.Decimal `json:"price"`
	Discount        decimal.Decimal `json:"discount"` // Percentage (0-100)
	ImageURL        string          `json:"imageUrl"`
	IsActive        bool            `json:"isActive"`
	CreatedAt       time.Time       `json:"createdAt"`
	UpdatedAt       time.Time       `json:"updatedAt"`
}

// FinalPrice calculates the final price after applying discount
// Formula: price * (1 - discount/100)
func (p *Package) FinalPrice() decimal.Decimal {
	if p.Discount.IsZero() {
		return p.Price
	}

	// Convert discount percentage to decimal (e.g., 10% -> 0.10)
	discountMultiplier := p.Discount.Div(decimal.NewFromInt(100))

	// Calculate discount amount
	discountAmount := p.Price.Mul(discountMultiplier)

	// Return final price
	return p.Price.Sub(discountAmount)
}

// RequiredSlots calculates the number of 20-minute slots needed
func (p *Package) RequiredSlots() int {
	return int(p.DurationMinutes / 20)
}

// ValidatePackage checks if package data is valid
func (p *Package) ValidatePackage() error {
	if p.DurationMinutes <= 0 || p.DurationMinutes%20 != 0 {
		return ErrInvalidDuration
	}
	if p.Price.IsNegative() {
		return ErrInvalidPrice
	}
	if p.Discount.IsNegative() || p.Discount.GreaterThan(decimal.NewFromInt(100)) {
		return ErrInvalidDiscount
	}
	return nil
}

// Errors
var (
	ErrInvalidDuration = NewValidationError("duration must be a positive multiple of 20 minutes")
	ErrInvalidPrice    = NewValidationError("price cannot be negative")
	ErrInvalidDiscount = NewValidationError("discount must be between 0 and 100")
)

// ValidationError represents a validation error
type ValidationError struct {
	Message string
}

func NewValidationError(msg string) *ValidationError {
	return &ValidationError{Message: msg}
}

func (e *ValidationError) Error() string {
	return e.Message
}
